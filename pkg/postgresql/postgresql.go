package postgresql

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	// sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	// sqlxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
)

type DBConfig struct {
	Host                   string
	Port                   string
	User                   string
	Password               string
	DBName                 string
	DBParams               string
	DBMaxConnection        int
	DBMaxIdleConnection    int
	DBMinuteTimeConnection int
}

func (conf DBConfig) Connect() *sqlx.DB {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Errors")
			fmt.Println("Recovered from panic:", r)
		}
	}()

	// Register PostgreSQL driver for tracing
	// sqltrace.Register("postgres", &pq.Driver{}, sqltrace.WithServiceName(servicename))

	// Connect to PostgreSQL database with tracing
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName, conf.DBParams)
	log.Println("Connecting to PostgreSQL: ", dsn)

	db, err := sqlx.ConnectContext(context.Background(), "postgres", dsn)
	if err != nil {
		msg := fmt.Sprintf("Cannot connect to PostgreSQL: %s, %v", dsn, err)
		slog.Error(msg)
		panic(msg)
	}

	// Set database connection pool settings
	if conf.DBMaxConnection == 0 {
		conf.DBMaxConnection = 20
	}

	if conf.DBMaxIdleConnection == 0 {
		conf.DBMaxIdleConnection = 20
	}

	if conf.DBMinuteTimeConnection == 0 {
		conf.DBMinuteTimeConnection = 3
	}

	db.SetMaxOpenConns(conf.DBMaxConnection)
	db.SetMaxIdleConns(conf.DBMaxIdleConnection)
	db.SetConnMaxLifetime(time.Duration(conf.DBMinuteTimeConnection) * time.Minute)

	return db
}

func Close(db *sqlx.DB) error {
	return db.Close()
}
