package config

import "time"

const (
	DefaultHost             = "ws://0.0.0.0"
	DefaultPort             = "26658"
	DefaultJWT              = "use: celestia light auth admin"
	DefaultSettingsPathName = ".waypoint"
	DefaultMeter            = "celestia/waypoint"
	DefaultReadInterval     = 1 * time.Minute
	DefaultInfoInterval     = 20 * time.Second
)
