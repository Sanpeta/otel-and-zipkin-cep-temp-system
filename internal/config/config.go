package config

import (
	"github.com/spf13/viper"
)

// Config is a configuration for the application
type Config struct {
	ENV               string `mapstructure:"ENV"`
	SERVER_ADDRESS    string `mapstructure:"SERVER_ADDRESS"`
	TOKEN_WEATHER_API string `mapsctructure:"TOKEN_WEATHER_API"`
}

// LoadConfig loads configuration from file
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
