package config

import "time"

const (
	DefaultPort             = "26658"
	DefaultHost             = "ws://0.0.0.0"
	DefaultJWT              = "use: celestia light auth admin"
	DefaultSettingsPathName = ".waypoint"
	DefaultMeter            = "celestia/waypoint"
	DefaultReadInterval     = time.Second * 10
	DefaultInfoInterval     = time.Second * 20
)
