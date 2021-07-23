package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "tida_ng"
	password = "huong"
	dbname   = "mydatabase"
)

func createConnection() *sql.DB {
	// connection string
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", conn)

	if err != nil {
		panic(err)
	}

	// close database
	defer db.Close()

	// check db
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	return db
}

