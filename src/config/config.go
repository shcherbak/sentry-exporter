package config

import (
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type Config struct {
	API_TOKEN         string
	APP_ENV           string
	AUTH_TOKEN        string
	BASE_URL          string
	LISTEN_PORT       string
	ORGANIZATION_SLUG string
	PROJECTS_EXCLUDE  string
	PROJECTS_INCLUDE  string
	REDIS_ADDR        string
	REDIS_DBNO        string
	REDIS_PORT        string
	ROUTINE_MAX       string
	SLEEP_SEC         string
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
