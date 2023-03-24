package config

import (
	"errors"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
}

type AppConfig struct {
	Name    string
	Version string
	Port    string
	Env     string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func LoadConfig() (Config, error) {
	// Load Config
	v := viper.New()

	v.SetConfigName("config/config")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return Config{}, errors.New("config file not found")
		}
		return Config{}, err
	}

	// Parse Config
	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return Config{}, err
	}

	return c, nil
}

func LoadConfigPath(path string) (Config, error) {
	// Load Config
	v := viper.New()
	v.SetConfigName(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return Config{}, errors.New("config file not found")
		}
		return Config{}, err
	}

	// Parse Config
	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return Config{}, err
	}

	return c, nil
}
