package utils

import (
	"fmt"
	"os"
	"reflect"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Username       string `env:"GENOTE_USER" required:"true"`
	Password       string `env:"GENOTE_PASSWORD" required:"true"`
	DiscordWebhook string `env:"DISCORD_WEBHOOK" required:"true"`
}

var (
	instance *Config
	once     sync.Once
	loadErr  error
)

func GetConfig() (*Config, error) {
	once.Do(func() {
		instance, loadErr = loadEnvVariables()
	})
	return instance, loadErr
}

func MustGetConfig() *Config {
	config, err := GetConfig()
	if err != nil {
		panic(err)
	}
	return config
}

func loadEnvVariables() (*Config, error) {
	config := &Config{}

	env_err := godotenv.Load()

	if env_err != nil {
		fmt.Printf("Missing .env file. Make sure to create one or set the environment variables manually.\n")
	}

	t := reflect.TypeOf(config).Elem()
	val := reflect.ValueOf(config).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := val.Field(i)

		key := field.Tag.Get("env")
		required := field.Tag.Get("required") == "true"

		envValue := os.Getenv(key)

		if required && envValue == "" {
			if val.Field(i).String() == "" {
				return nil, fmt.Errorf("missing required environment variable %s", key)
			}
		}

		value.SetString(envValue)
	}

	return config, nil
}
