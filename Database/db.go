package Database

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var (
	DBconn *sqlx.DB
)

 feat/comments
// 6.todo db credentials takes from env
// 7.todo write migration
=======
type SSLMode string

const (
	SSLModeDisable SSLMode = "disable"
)

 master
func init() {
	dbs := "host=localhost port=5433 user=local password=local dbname=todo sslmode=ss"
	//fmt.Println("line printing")
	db, err := sqlx.Open("postgres", dbs)

	if err != nil {
		fmt.Println("Error in DB connection", err)
	}
	DBconn = db
	log.Println("database connected..")
}

func ConnectAndMigrate(host, port, databaseName, user, password string, sslMode SSLMode) error {

	connectionSTR := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, databaseName, password, sslMode)
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

func DbConnectionClose() error {
	return DBconn.Close()
}

func migrateUp(db *sqlx.DB) error {
	// Create a new migration driver for PostgresSQL using the database instance
	driver, driErr := postgres.WithInstance(db.DB, &postgres.Config{})
	if driErr != nil {
		// Return error if migration driver creation fails
		return driErr
	}
	// Set up the migration instance with the file path for migrations and the database driver
	m, migErr := migrate.NewWithDatabaseInstance(
		"file://database/migrations", // Path to migration files
		"postgres", driver)           // Database driver and name

	if migErr != nil {
		// Return error if migration instance creation fails
		return migErr
	}
	// Run the migrations (updating the database schema)
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		// If an error occurs, but it's not "ErrNoChange" (no changes detected), return the error
		return err
	}
	// Return nil if migration was successful or no changes were needed
	return nil
}
