package config

import (
	"beli-mang/pkg/aws"
	"beli-mang/pkg/postgresql"
	"beli-mang/pkg/str"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// Configs ...
type Configs struct {
	EnvConfig map[string]string
	DB        *sqlx.DB
	S3Config  aws.S3Config
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

	s3Config := aws.S3Config{
		AccessKeyId:     conf.EnvConfig["AWS_ACCESS_KEY_ID"],
		SecretAccessKey: conf.EnvConfig["AWS_SECRET_ACCESS_KEY"],
		BucketName:      conf.EnvConfig["AWS_S3_BUCKET_NAME"],
		Region:          conf.EnvConfig["AWS_REGION"],
	}
	s3Config.Initialize()
	conf.S3Config = s3Config

	return conf, err
}

func GetConfig(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}
