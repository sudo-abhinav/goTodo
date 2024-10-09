package Database

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	DBconn *sqlx.DB
)

//feat/comments
// 6.todo db credentials takes from env : = working on that
// 7.todo write migration : done

type SSLMode string

const (
	SSLModeDisable SSLMode = "disable"
)

//func init() {
//	dbs := "host=localhost port=5433 user=local password=local dbname=todo sslmode=ss"
//	//fmt.Println("line printing")
//	db, err := sqlx.Open("postgres", dbs)
//
//	if err != nil {
//		fmt.Println("Error in DB connection", err)
//	}
//	DBconn = db
//	log.Println("database connected..")
//}

func ConnectAndMigrate(host, port, databaseName, user, password string, sslMode SSLMode) error {

	connectionSTR := fmt.Sprintf("host=%s port=%s user=%s password=%s  dbname=%s sslmode=%s", host, port, user, password, databaseName, sslMode)
	db, err := sqlx.Open("postgres", connectionSTR)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {

		return err
	}
	DBconn = db

	return migrateUp(db)

}

func ShutDownDB() error {
	return DBconn.Close()
}

func migrateUp(db *sqlx.DB) error {
	driver, driErr := postgres.WithInstance(db.DB, &postgres.Config{})
	if driErr != nil {
		return driErr
	}
	m, migErr := migrate.NewWithDatabaseInstance(
		"file://Database/migrations/", // Path to migration files
		"postgres", driver)            // Database driver and name
	if migErr != nil {
		// Return error if migration instance creation fails
		return migErr
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		// If an error occurs, but it's not "ErrNoChange" (no changes detected), return the error
		return err
	}

	return nil
}
