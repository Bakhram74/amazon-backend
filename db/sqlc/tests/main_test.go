package tests

import (
	"context"
	"fmt"
	"github.com/Bakhram74/amazon.git/internal/config"
	"github.com/Bakhram74/amazon.git/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

var testStore repository.Store

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("../../../configs")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), postgresUrl(config.Storage))
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = repository.NewStore(connPool)
	os.Exit(m.Run())
}

func postgresUrl(cfg config.StorageConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)
}
