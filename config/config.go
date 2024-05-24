package config

import (
	"beli-mang/pkg/postgresql"
	"beli-mang/pkg/str"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// Configs ...
type Configs struct {
	EnvConfig map[string]string
	DB        *sqlx.DB
}

func LoadConfigs() (conf Configs, err error) {
	conf.EnvConfig, err = godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error loading ..env file")
	}

	// postgresql conn
	sqlConfig := postgresql.DBConfig{
		DBName:                 conf.EnvConfig["DB_NAME"],
		Host:                   conf.EnvConfig["DB_HOST"],
		Port:                   conf.EnvConfig["DB_PORT"],
		User:                   conf.EnvConfig["DB_USERNAME"],
		Password:               conf.EnvConfig["DB_PASSWORD"],
		DBParams:               conf.EnvConfig["DB_PARAMS"],
		DBMaxConnection:        str.StringToInt(conf.EnvConfig["DB_MAX_CONNECTION"]),
		DBMaxIdleConnection:    str.StringToInt(conf.EnvConfig["DB_MAX_IDLE_CONNECTION"]),
		DBMinuteTimeConnection: str.StringToInt(conf.EnvConfig["DB_MAX_LIFETIME_CONNECTION"]),
	}
	conf.DB = sqlConfig.Connect()

	return conf, err
}
