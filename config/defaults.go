package config

import "time"

const (
	DefaultPort             = "26658"
	DefaultHost             = "http://0.0.0.0"
	DefaultJWT              = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.z6M-aDkCsWMaaQrhK6JVpRPwdAwOBUoruhjieFd5iCE"
	DefaultSettingsPathName = ".waypoint"
	DefaultMeter            = "celestia/waypoint"
	DefaultReadInterval     = time.Second * 10
	DefaultInfoInterval     = time.Second * 20
)
