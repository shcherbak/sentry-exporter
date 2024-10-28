package config

import (
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type Config struct {
	APP_ENV           string
	API_TOKEN         string
	AUTH_TOKEN        string
	BASE_URL          string
	ORGANIZATION_SLUG string
	ROUTINE_MAX       string
	SLEEP_SEC         string
	LISTEN_PORT       string
	REDIS_ADDR        string
	REDIS_PORT        string
	REDIS_DBNO        string
	TTL_SECONDS       string
}

var config *Config

func GetConfig() Config {
	return *config
}

func loadEnvFile() {
	log.Println("Loading .env file.")
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file.")
	}
}

func init() {
	config = &Config{}

	_, found := os.LookupEnv("APP_ENV")
	if !found {
		loadEnvFile()
	}

	_, found = os.LookupEnv("LISTEN_PORT")
	if !found {
		os.Setenv("LISTEN_PORT", "80")
	}

	_, found = os.LookupEnv("AUTH_TOKEN")
	if !found {
		os.Setenv("AUTH_TOKEN", "")
	}

	_, found = os.LookupEnv("ROUTINE_MAX")
	if !found {
		os.Setenv("ROUTINE_MAX", "3")
	}

	_, found = os.LookupEnv("SLEEP_SEC")
	if !found {
		os.Setenv("SLEEP_SEC", "45")
	}

	_, found = os.LookupEnv("TTL_SECONDS")
	if !found {
		os.Setenv("TTL_SECONDS", "300")
	}

	_, found = os.LookupEnv("ORGANIZATION_SLUG")
	if !found {
		os.Setenv("ORGANIZATION_SLUG", "sentry")
	}

	refl := reflect.ValueOf(config).Elem()
	numFields := refl.NumField()
	for i := 0; i < numFields; i++ {
		envName := refl.Type().Field(i).Name
		envVal, foud := os.LookupEnv(envName)
		if !foud {
			panic("Environment [" + envName + "] not found.")
		}
		refl.Field(i).SetString(envVal)
	}
}
