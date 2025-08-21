package config

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DBDriver       string `mapstructure:"DB_DRIVER"`
	DBSource       string `mapstructure:"DB_SOURCE"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	JWTSecret      string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn   string `mapstructure:"JWT_EXPIRES_IN"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
