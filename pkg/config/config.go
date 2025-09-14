package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var (
	config *viper.Viper
)

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Load(env string, configPaths ...string) {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("config/")
	config.AddConfigPath("../../")
	config.AddConfigPath(".")
	if len(configPaths) != 0 {
		for _, path := range configPaths {
			config.AddConfigPath(path)
		}
	}
	err = config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file", err)
		return
	}

	if env == CONFIG_SERVER {
		// existing behavior: server.yaml requires env vars for keys
		for _, v := range config.AllKeys() {
			if strings.ToLower(v) == "version" {
				continue // skipping first line
			}
			key := config.GetString(v)
			if key == "" {
				resp, err := setStringSlice(v)
				if err != nil {
					log.Fatal(err)
				}
				config.Set(v, resp)
			} else {
				key = strings.ReplaceAll(key, "$", "")
				if ev, ok := os.LookupEnv(key); ok {
					config.Set(v, ev)
				} else {
					log.Fatal("env value for key [", key, "] is missing")
				}
			}
		}
		log.Println("application running with server.yaml")
	} else {
		// For non-server (local/dev) loads, prefer common DB_* env vars when present so
		// docker-compose can pass DB_HOST=db etc. But fall back to YAML values if env not set.
		// This block does not require all env variables to be present.
		if ev := os.Getenv("DB_HOST"); ev != "" {
			config.Set("databases.postgres.host", ev)
		}
		if ev := os.Getenv("DB_PORT"); ev != "" {
			config.Set("databases.postgres.port", ev)
		}
		if ev := os.Getenv("DB_USER"); ev != "" {
			config.Set("databases.postgres.user", ev)
		}
		if ev := os.Getenv("DB_PASSWORD"); ev != "" {
			config.Set("databases.postgres.password", ev)
		}
		if ev := os.Getenv("DB_NAME"); ev != "" {
			config.Set("databases.postgres.db", ev)
		}
		log.Println("application running with ", env, ".yaml")
	}
}

func GetConfig() *viper.Viper {
	return config
}

func setStringSlice(key string) ([]string, error) {
	var (
		resp []string
		err  error
	)

	keys := config.GetStringSlice(key)

	for _, k := range keys {
		k = strings.ReplaceAll(k, "$", "")
		if ev, ok := os.LookupEnv(k); ok {
			resp = append(resp, ev)
		} else {
			return resp, fmt.Errorf("env value for key [ %v ] is missing", k)
		}
	}
	return resp, err
}
