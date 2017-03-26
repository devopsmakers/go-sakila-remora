package mysql

import (
	"bytes"

	"github.com/devopsmakers/go-sakila-remora/remora"
)

// MySQL type for interface
type MySQL struct{}

// Check - logic to decide whether this service is healthy
func (m MySQL) Check(c *remora.Config) remora.Result {
	status := 1
	body := bytes.NewBufferString("Just testing stuff out here")

	return remora.Result{StatusCode: status, Body: *body}
}
