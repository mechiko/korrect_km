package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type IConfig interface {
	Configuration() *Configuration
	SetInConfig(key string, value interface{}, save ...bool) error
	GetKeyString(name string) string
	GetByName(name string) interface{}
	DatabaseByKey(key string) *DatabaseConfiguration
}

type Config struct {
	*viper.Viper
	configuration  *Configuration
	configFileName string
	warning        string
}

var _ IConfig = (*Config)(nil)

const modError = "pkg:config"

var UserHomeDir string

func New(cfgName string, userHome bool) (cfg *Config, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s panic %v", modError, r)
		}
	}()

	if userHome {
		UserHomeDir = userHomeDir()
	}
	configName := cfgName
	if cfgName == "" {
		configName = "config"
	}

	viperOrigin := viper.GetViper()
	configFileName := configName + ".toml"
	configFileName = filepath.Join(UserHomeDir, ConfigPath, configFileName)
	configPath := filepath.Join(UserHomeDir, ConfigPath)

	viperOrigin.SetConfigName(configName)
	viperOrigin.SetConfigType("toml")
	viperOrigin.AddConfigPath(configPath)

	if err := viperOrigin.MergeConfig(strings.NewReader(string(TomlConfig))); err != nil {
		return nil, fmt.Errorf("%s %w", modError, err)
	}
	warn := ""
	if err := viperOrigin.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("%s %w", modError, err)
		} else {
			warn = fmt.Sprintf("%s Config file ('%s') not found", modError, configFileName)
		}
	}

	conf := &Configuration{}
	if err := viperOrigin.Unmarshal(conf); err != nil {
		return nil, fmt.Errorf("%s %w", modError, err)
	}

	cfg = &Config{
		Viper:          viperOrigin,
		configuration:  conf,
		configFileName: configFileName,
		warning:        warn,
	}
	viperOrigin.SafeWriteConfig()
	return cfg, nil
}

func (c *Config) Configuration() *Configuration {
	return c.configuration
}

func userHomeDir() string {
	switch runtime.GOOS {
	case "windows":
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	case "linux":
		home := os.Getenv("HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}
