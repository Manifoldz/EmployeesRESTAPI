package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Manifoldz/EmployeesRESTAPI/internal/server"
	"github.com/Manifoldz/EmployeesRESTAPI/pkg/handler"
	"github.com/Manifoldz/EmployeesRESTAPI/pkg/repository"
	"github.com/Manifoldz/EmployeesRESTAPI/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @author Manifoldz
// @title Employees REST API
// @version 1.0
// @description This is a REST API for a Employees list
// @host localhost:8000
// @BasePath /

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error while reading config file: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error while loading.env file:  %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("error while connecting to database: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	go func() {
		if err := srv.Start(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error while starting http server: %s", err.Error())
		}
	}()
	logrus.Printf("Server started on port:  %s", viper.GetString("port"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Printf("Server shutting down on port:  %s", viper.GetString("port"))

	if err := srv.Stop(context.Background()); err != nil {
		logrus.Fatalf("error while stopping http server:  %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Fatalf("error while closing database connection:  %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
