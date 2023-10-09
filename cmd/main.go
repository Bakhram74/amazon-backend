package main

import (
	"context"
	"github.com/Bakhram74/amazon-backend.git/internal/config"
	"github.com/Bakhram74/amazon-backend.git/internal/handler"
	"github.com/Bakhram74/amazon-backend.git/internal/repository"
	"github.com/Bakhram74/amazon-backend.git/internal/server"
	"github.com/Bakhram74/amazon-backend.git/internal/service"
	"github.com/Bakhram74/amazon-backend.git/pkg/client/postgresql"
	"github.com/Bakhram74/amazon-backend.git/pkg/logging"

	_ "github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logging.GetLogger()
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		logger.Fatal("cannot load config")
	}

	storage := cfg.Storage
	connPool, err := postgresql.NewClient(context.TODO(), 3, postgresql.Config{
		Host:     storage.Host,
		Port:     storage.Port,
		Username: storage.Username,
		Password: storage.Password,
		DBName:   storage.Database,
		SSLMode:  storage.SSLMode,
	})
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}
	store := repository.NewStore(connPool)
	repos := repository.NewRepository(store)
	services := service.NewService(repos)
	handlers, err := handler.NewHandler(cfg, services)
	if err != nil {
		logger.Fatal(err.Error())
	}
	svr := new(server.Server)
	logger.Infof("server is listening address %s", cfg.HttpAddress)
	if err = svr.Run(cfg.HttpAddress, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}

}
