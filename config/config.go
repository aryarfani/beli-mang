package config

import (
	"beli-mang/pkg/aws"
	"beli-mang/pkg/postgresql"
	"beli-mang/pkg/str"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// Configs ...
type Configs struct {
	DB       *sqlx.DB
	S3Config aws.S3Config
}

func LoadConfigs() (conf Configs, err error) {
	// postgresql conn
	sqlConfig := postgresql.DBConfig{
		DBName:                 GetConfig("DB_NAME"),
		Host:                   GetConfig("DB_HOST"),
		Port:                   GetConfig("DB_PORT"),
		User:                   GetConfig("DB_USERNAME"),
		Password:               GetConfig("DB_PASSWORD"),
		DBParams:               GetConfig("DB_PARAMS"),
		DBMaxConnection:        str.StringToInt(GetConfig("DB_MAX_CONNECTION")),
		DBMaxIdleConnection:    str.StringToInt(GetConfig("DB_MAX_IDLE_CONNECTION")),
		DBMinuteTimeConnection: str.StringToInt(GetConfig("DB_MAX_LIFETIME_CONNECTION")),
	}
	conf.DB = sqlConfig.Connect()

	s3Config := aws.S3Config{
		AccessKeyId:     GetConfig("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: GetConfig("AWS_SECRET_ACCESS_KEY"),
		BucketName:      GetConfig("AWS_S3_BUCKET_NAME"),
		Region:          GetConfig("AWS_REGION"),
	}
	s3Config.Initialize()
	conf.S3Config = s3Config

	return conf, err
}

func GetConfig(key string) string {
	// load .env file
	godotenv.Load(".env")

	return strings.TrimSpace(os.Getenv(key))
}
