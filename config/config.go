package config

import (
	"encoding/json"

	"github.com/maxiiot/devicebridge/backend/mqtt"
)

// Cfg configuration entity
var Cfg Configuration

// Configuration vbasebridge app's configurations
type Configuration struct {
	General struct {
		Port     int    `mapstructure:"port" json:"port"`
		LogLevel string `mapstructure:"log_level" json:"log_level"`
	} `mapstructure:"general" json:"general"`

	LoraServer struct {
		Server   string `mapstructure:"server" json:"server"`
		UserName string `mapstructure:"username" json:"username"`
		Password string `mapstructure:"password" json:"password"`
	} `mapstructure:"loraserver" json:"loraserver"`

	LoraBackend struct {
		Mqtt     mqtt.Config `mapstructure:"mqtt" json:"mqtt"`
		HTTPPort int         `mapstructure:"http_port" json:"http_port"`
	} `mapstructure:"lora_backend" json:"lora_backend"`

	Postgres struct {
		AutoMigrate bool   `mapstructure:"auto_migrate" json:"auto_migrate"`
		DSN         string `mapstructure:"dsn" json:"dsn"`
	} `mapstructure:"postgres" json:"postgres"`

	Publisher struct {
		Mqtt mqtt.Config `mapstructure:"mqtt" json:"mqtt"`
	} `mapstructure:"publisher" json:"publisher"`
}

func (c Configuration) String() string {
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err.Error()
	}

	return string(b)
}
