package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type DBConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Secrets struct {
	DBName     string
	DBUser     string
	DBPassword string
}

type Config struct {
	Env      string `yaml:"env"`
	DBConfig `yaml:"db"`
	Secrets
}

func MustLoad(pathToConfig string, pathToSecret string) Config {
	cfg := mustLoad[Config](pathToConfig)
	switch cfg.Env {
	case EnvLocal, EnvDev, EnvProd:
		err := godotenv.Load(pathToSecret)
		if err != nil {
			log.Fatal("secrets load failed")
		}
		cfg.DBName = os.Getenv("POSTGRES_DB")
		cfg.DBUser = os.Getenv("POSTGRES_USER")
		cfg.DBPassword = os.Getenv("POSTGRES_PASSWORD")
	}
	return cfg
}

func mustLoad[T any](pathToConfig string) T {
	if pathToConfig == "" {
		log.Fatal("Empty path to config file")
	}

	if _, err := os.Stat(pathToConfig); os.IsNotExist(err) {
		log.Fatalf("Incorrect path to config file: '%v'", pathToConfig)
	}

	var cfg T

	cfgData, err := os.ReadFile(pathToConfig)
	if err != nil {
		log.Fatalf("Error reading file: %v", err.Error())
	}

	err = yaml.Unmarshal(cfgData, &cfg)
	if err != nil {
		log.Fatalf("Error unmarshaling file: %v", err.Error())
	}

	return cfg
}
