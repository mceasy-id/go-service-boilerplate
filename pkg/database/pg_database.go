package database

import (
	"database/sql"
	"fmt"
	"mceasy/service-demo/config"

	_ "github.com/lib/pq"
)

func GetDatabaseConnection(config *config.Config) (*sql.DB, error) {
	dbHost := config.Postgres.Host
	dbPort := config.Postgres.Port
	dbUser := config.Postgres.User
	dbPassword := config.Postgres.Password
	dbName := config.Postgres.DBName

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dsn)

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
