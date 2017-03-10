package remora

import "encoding/json"

// SlaveStatus is used to marshal required values into
type SlaveStatus struct {
	secondsBehindMaster int
}

// GetSlaveStatus returns values obtained from SHOW SLAVE STATUS
func GetSlaveStatus(config *Config) (int, []byte, error) {
	vals := make(map[string]string)
	vals["seconds_behind_master"] = "5"
	vals["something_else_running"] = "Yes"
	returnvals, err := json.Marshal(vals)
	return 1, returnvals, err
}
