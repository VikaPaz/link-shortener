package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Config struct{ Host, Port, User, Password, Dbname string }

func Connection(conf Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, nil
}

func NewRepository(conn *sql.DB) *repositoryImpl {
	// probably, add logic
	return &repositoryImpl{
		conn: conn,
	}
}
