package mysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"sort"
	"strconv"

	"github.com/devopsmakers/go-sakila-remora/remora"
	// Blank import is fine here
	_ "github.com/go-sql-driver/mysql"
	jww "github.com/spf13/jwalterweatherman"
)

// MySQL type for interface
type MySQL struct{}

// Check - logic to decide whether this service is healthy
func (m MySQL) Check(c *remora.Config) remora.Result {

	// variables for results
	var status int
	var body = bytes.NewBufferString("OK")

	// Connect to DB
	user := c.Service.User
	pass := c.Service.Pass
	host := c.Service.Host
	port := c.Service.Port
	lag, _ := strconv.Atoi(c.AcceptableLag)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, pass, host, port))
	if err != nil {
		status = 2
		body = bytes.NewBufferString(fmt.Sprintf("Error: %s", err.Error()))
		jww.ERROR.Println(body.String())
		return remora.Result{StatusCode: status, Body: *body}
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		status = 2
		body = bytes.NewBufferString(fmt.Sprintf("Error: %s", err.Error()))
		jww.ERROR.Println(body.String())
		return remora.Result{StatusCode: status, Body: *body}
	}

	//Connection is all good, first, let's check if replication is running.

	// Get slave status
	rows, err := db.Query("SHOW SLAVE STATUS")
	if err != nil {
		status = 2
		body = bytes.NewBufferString(fmt.Sprintf("Error: %s", err.Error()))
		jww.ERROR.Println(body.String())
		return remora.Result{StatusCode: status, Body: *body}
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		status = 2
		body = bytes.NewBufferString(fmt.Sprintf("Error: %s", err.Error()))
		jww.ERROR.Println(body.String())
		return remora.Result{StatusCode: status, Body: *body}
	}

	// Create slice for values
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Create a map
	slaveStatus := make(map[string]string)
	for rows.Next() {
		// get raw bytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			status = 2
			body = bytes.NewBufferString(fmt.Sprintf("Error: %s", err.Error()))
			jww.ERROR.Println(body.String())
			return remora.Result{StatusCode: status, Body: *body}
		}
		var value string

		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			slaveStatus[string(columns[i])] = string(value)
		}
	}

	// make sure we are a slave
	if len(slaveStatus) < 1 {
		status = 2
		body = bytes.NewBufferString("This server is not a slave")
		jww.ERROR.Println(body.String())
		return remora.Result{StatusCode: status, Body: *body}
	}

	// Test lag, IO and SQL running
	secsBehind, _ := strconv.Atoi(slaveStatus["Seconds_Behind_Master"])

	if int(secsBehind) > int(lag) ||
		slaveStatus["Slave_IO_Running"] != "Yes" ||
		slaveStatus["Slave_SQL_Running"] != "Yes" {
		status = 1
	}

	// Give some useful output
	body = bytes.NewBufferString("")

	// Nicely order our output
	var sortKeys []string
	for k, v := range slaveStatus {
		if len(v) > 0 {
			sortKeys = append(sortKeys, k)
		}
	}

	sort.Strings(sortKeys)
	for _, k := range sortKeys {
		body.WriteString(fmt.Sprintf("%s: %s\n", k, slaveStatus[k]))
	}

	return remora.Result{StatusCode: status, Body: *body}
}
