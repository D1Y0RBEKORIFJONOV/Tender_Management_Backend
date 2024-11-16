package databaseconnection

import (
	"awesomeProject/internal/config"
	"awesomeProject/logger"
	"database/sql"
	"fmt"
)

type Database struct {
	Db *sql.DB
}

func Connect() (*Database, error) {
	c := config.New()
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Dbname))
	if err != nil {
		logger.SetupLogger(err.Error())
	}
	if err := db.Ping(); err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}
	return &Database{Db: db}, nil
}
