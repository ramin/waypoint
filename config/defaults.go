package config

import "time"

const (
	DefaultPort             = "26658"
	DefaultHost             = "http://0.0.0.0"
	DefaultJWT              = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.dr_T72dupkAfSpMI5Db4Qgc3CIdipClgL3Aer7O_40s"
	DefaultSettingsPathName = ".waypoint"
	DefaultMeter            = "celestia/waypoint"
	DefaultReadInterval     = time.Second * 10
)
