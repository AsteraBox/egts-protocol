package config

/*
Описание конфигурационного файла
*/

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Settings struct {
	Host     string                       `yaml:"host"`
	Port     string                       `yaml:"port"`
	ConnTTl  int                          `yaml:"conn_ttl"`
	LogLevel string                       `yaml:"log_level"`
	Storage  map[string]map[string]string `yaml:"storage"`
}

func (s *Settings) GetEmptyConnTTL() time.Duration {
	return time.Duration(s.ConnTTl) * time.Second
}
func (s *Settings) GetListenAddress() string {
	return s.Host + ":" + s.Port
}

func (s *Settings) GetLogLevel() log.Level {
	var lvl log.Level

	switch strings.ToUpper(s.LogLevel) {
	case "DEBUG":
		lvl = log.DebugLevel
	case "INFO":
		lvl = log.InfoLevel
	case "WARNING", "WARN":
		lvl = log.WarnLevel
	case "ERROR":
		lvl = log.ErrorLevel
	default:
		lvl = log.InfoLevel
	}
	return lvl
}

func New(configPath string) (*Settings, error) {
	viper.SetDefault("Host", "0.0.0.0")
	viper.SetDefault("Port", "6000")
	viper.SetDefault("ConnTTl", 10)
	viper.SetDefault("LogLevel", "DEBUG")

	viper.RegisterAlias("ConnTTl", "con_live_sec")
	viper.RegisterAlias("LogLevel", "log_level")

	viper.AddConfigPath(filepath.Dir(configPath))
	viper.AddConfigPath(".")
	viper.SetConfigName(filepath.Base(configPath))
	viper.SetConfigType(filepath.Ext(configPath)[1:])

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("could not read the config file: %v", err)
	}

	var config *Settings
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal: %v", err)
	}

	return config, nil
}
