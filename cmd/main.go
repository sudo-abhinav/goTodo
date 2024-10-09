package main

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/sudo-abhinav/go-todo/Database"
	"github.com/sudo-abhinav/go-todo/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutDownTime = 10 * time.Second

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := routes.SetupRoutes()

	fmt.Println(os.Getenv("DB_HOST"), "hello")

	if err := Database.ConnectAndMigrate(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		Database.SSLModeDisable); err != nil {
		logrus.Panicf("Failed to initialize and migrate database with error = %+v", err)
	}
	logrus.Print("migration successfully..")

	go func() {
		if err := server.Run(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Panicf("Failed to run server with error: %+v", err)
		}
	}()
	logrus.Print("Server started at :8080")
	<-done

	logrus.Info("shutting down server")

	if err := Database.ShutDownDB(); err != nil {
		logrus.WithError(err).Error("failed to close database connection")
	}

	if err := server.Shutdown(shutDownTime); err != nil {
		logrus.WithError(err).Panic("failed to gracefully shutdown server")
	}
}
