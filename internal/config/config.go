package config

import "github.com/spf13/viper"

type Config struct {
	App AppConfig
	DB  DBConfig
}
type AppConfig struct {
	Port string
	Env  string
}
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	cfg := &Config{
		App: AppConfig{
			Port: viper.GetString("APP_PORT"),
			Env:  viper.GetString("APP_ENV"),
		},
		DB: DBConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
	}
	return cfg, nil
}
