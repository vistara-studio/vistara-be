package db

import (
	"fmt"
	"time"

	"github.com/vistara-studio/vistara-be/internal/infra/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(conf *config.Env) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conf.PostgresHost,
		conf.PostgresPort,
		conf.PostgresUsername,
		conf.PostgresPassword,
		conf.PostgresDB,
		conf.PostgresSSL,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(10 * time.Second)
	db.SetConnMaxLifetime(1 * time.Hour)

	return db, nil
}
