package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	cfg  Config
	once sync.Once
)

// Config stores all service configuration that can be injected via
// ENV var or falls back to defaults
type Config struct {
	Verbosity    string `mapstructure:"verbosity"`
	LogFormatter string `mapstructure:"logFormatter"`
	Environment  string `mapstructure:"environment"`

	Listen string `mapstructure:"listen"`

	JWT string `mapstructure:"jwt"`

	// address of light node (ie: for receiving UTIA)
	Address string `mapstructure:"address"`

	Blocks       int           `mapstructure:"blocks"`
	Delay        time.Duration `mapstructure:"delay"`
	Namespace    string        `mapstructure:"namespace"`
	ReadInterval time.Duration `mapstructure:"readInterval"`

	// an internal with which to output some debug data
	// about the running processes
	LogInterval time.Duration `mapstructure:"readInterval"`
	DisplayInfo bool          `mapstructure:"displayInfo"`

	Meter string `mapstructure:"meter"`

	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// Read returns an instance of config, initializing
// it only once ever, handling settings defaults and
// binding ENV vars to structural config values
func Read() *Config {
	once.Do(func() {
		path, err := settingsPath()
		if err != nil {
			log.Error(err)
			return
		}

		viper.SetConfigName("config") // name of config file (without extension)
		viper.SetConfigType("toml")   // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath(path)     // call multiple times to add many search paths

		viper.SetDefault("verbosity", "INFO")
		viper.SetDefault("logFormatter", "")
		viper.SetDefault("settingsPath", "")
		viper.SetDefault("host", DefaultHost)
		viper.SetDefault("port", DefaultPort)
		viper.SetDefault("jwt", DefaultJWT)
		viper.SetDefault("listen", "0.0.0.0:8080")
		viper.SetDefault("meter", DefaultMeter)
		viper.SetDefault("readInterval", DefaultReadInterval)
		viper.SetDefault("infoInterval", DefaultInfoInterval)

		_ = viper.BindEnv("jwt", "JWT")

		_ = viper.BindEnv("blocks", "BLOCKS")
		_ = viper.BindEnv("delay", "DELAY")
		_ = viper.BindEnv("readInterval", "READ_INTERNAL")
		_ = viper.BindEnv("infoInterval", "INFO_INTERNAL")
		_ = viper.BindEnv("host", "HOST")
		_ = viper.BindEnv("port", "PORT")
		_ = viper.BindEnv("namespace", "NAMESPACE")
		_ = viper.BindEnv("listen", "LISTEN")
		_ = viper.BindEnv("meter", "METER")
		_ = viper.BindEnv("displayInfo", "DISPLAY_INFO")

		_ = viper.BindEnv("verbosity", "VERBOSITY")
		_ = viper.BindEnv("logFormatter", "LOG_FORMATTER")
		_ = viper.BindEnv("settingsPath", "SETTINGS_PATH")

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error if desired
				log.Info("no config found, consider writing defaults")

				exists, err := Exists(path)
				if err != nil {
					log.Error(err)
					return
				}

				if !exists {
					err = os.Mkdir(path, 0700)
					if err != nil {
						log.Error(err)
						return
					}
				}

				err = viper.SafeWriteConfig()
				if err != nil {
					log.Error(err)
					return
				}
			} else {
				// Config file was found but another error was produced
				log.Error(err)
				return
			}
		}

		cfg = Config{}
		err = viper.Unmarshal(&cfg)

		if err != nil {
			log.Error(err)
			return
		}
	})

	return &cfg
}

// Exists test if any settings file exists
func Exists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	}

	return false, nil
}

// SavePath is a helper for setting
// up the user's home path
func settingsPath() (string, error) {
	// if path := config.Read().SettingsPath; path != "" {
	// 	return fmt.Sprintf("%s/%s", path, SettingsDirName), nil
	// }

	return defaultSettingsPath()
}

func defaultSettingsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", home, DefaultSettingsPathName), nil
}
