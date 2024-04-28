package app

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"hw3/storage/internal/config"
	"hw3/storage/internal/converter"
	"hw3/storage/internal/repository"
	"hw3/storage/internal/repository/redis"
	"hw3/storage/internal/service"
	"hw3/storage/internal/transport/kafka"
)

const (
	serverConfigPath = "server_config.yaml"
)

func RunStorage() {
	cfg, err := config.ParseServerConfig(serverConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	dbConnString := setupDBConnString(cfg)
	pool := connectToDB(dbConnString)
	defer pool.Close()

	storageRepo := repository.NewStorageRepo(pool)
	cache := redis.NewStorageCache(cfg.Redis)
	storageService := service.NewStorageService(storageRepo, cache)
	messageConverter := converter.MessageConverter{}

	storageHandler, err := kafka.NewStorageHandler(cfg.Kafka.Addrs, storageService, messageConverter)
	if err != nil {
		msg := fmt.Sprintf("Error while init StorageHandler: %v", err)
		log.Fatal(msg)
	}

	log.Println(cfg)
	log.Println(storageHandler)
	err = storageHandler.ConsumeMessages(cfg.Kafka.Topic)
	if err != nil {
		log.Fatal(err)
	}
}

func setupDBConnString(cfg *config.ServerConfig) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)
}

func connectToDB(connString string) *pgxpool.Pool {
	poolCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		log.Fatal(err)
	}
	return pool
}
