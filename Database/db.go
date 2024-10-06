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

// 6.todo db credentials takes from env
// 7.todo write migration
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
