package Database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var (
	DBconn *sqlx.DB
)

/*
	1) Create a separate function for opening database connection
	2) Create a separate function for migration up
	3) Create a separate function for handling Transaction
	4) Keep db credential in environment variables and access them using os.Getenv()
	5) Write every db query in a migration file
*/

// No need to write this in init method

func init() {
	dbs := "host=localhost port=5433 user=local password=local dbname=todo sslmode=disable"
	//fmt.Println("line printing")
	db, err := sqlx.Open("postgres", dbs)

	if err != nil {
		fmt.Println("Error in DB connection", err)
	}
	DBconn = db
	log.Println("database connected..")
}

func DbConnectionClose() error {
	return DBconn.Close()
}
