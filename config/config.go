package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App            AppConfig
	Authentication AuthenticationConfig
	Logger         LoggerConfig
	Postgres       PostgresConfig
	Observability  ObservabilityConfig
	ExternalURI    ExternalURIConfig
}

type AppConfig struct {
	Name    string
	Version string
	Port    string
	Env     string
	Key     string
}

type LoggerConfig struct {
	Mode  string
	Level string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type ObservabilityConfig struct {
	Enable       bool
	Mode         string
	OtlpEndpoint string
}

type AuthenticationConfig struct {
	Key string
}

type ExternalURIConfig struct {
	// Config for external URI
	MasterData MasterData
	Identity   string
	Scope      string
}

type MasterData struct {
	Address        string
	Vehicle        string
	User           string
	CompanyProfile string
	Token          string
}

func (m MasterData) GetBearerToken() string {
	return "Bearer " + m.Token
}

func LoadConfig(env string) (Config, error) {
	// Load Config
	v := viper.New()

	v.SetConfigName(fmt.Sprintf("config/config-%s", env))
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
