package database

import (
	"fmt"
	"mceasy/service-demo/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetPostgreConnection(config config.Config) (*sqlx.DB, error) {
	dbHost := config.Postgres.Host
	dbPort := config.Postgres.Port
	dbUser := config.Postgres.User
	dbPassword := config.Postgres.Password
	dbName := config.Postgres.DBName

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sqlx.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(20)

	return db, nil
}
