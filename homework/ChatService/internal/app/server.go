package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"hw3/chat-service/internal/config"
	"hw3/chat-service/internal/converter"
	"hw3/chat-service/internal/repo/redis"
	"hw3/chat-service/internal/service"
	"hw3/chat-service/internal/transport/kafka"
	"hw3/chat-service/internal/transport/rest"
)

const (
	serverConfigPath = "server_config.yaml"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

func RunServer() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ParseServerConfig(serverConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	dbConnString := setupDBConnString(cfg)

	pool := connectToDB(dbConnString)
	defer pool.Close()

	mux := http.NewServeMux()
	storageCache := redis.NewStorageCache(cfg.Redis)

	chatService := service.NewChatService(storageCache)
	messageConverter := converter.MessageConverter{}

	addrs := []string{"kafka1:9092"}
	kafkaHandler, err := kafka.NewChatHandler(addrs)
	if err != nil {
		log.Fatal(err)
	}

	chatHandler := rest.NewChatHandler(chatService, upgrader, messageConverter, kafkaHandler)

	mux.HandleFunc("/chat", chatHandler.Chat)

	runServer(ctx, mux, cfg)
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

func runServer(ctx context.Context, mux *http.ServeMux, cfg *config.ServerConfig) {
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	server := http.Server{
		Addr:    port,
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	log.Printf("Run server on %s\n", port)
	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Println("Server shutdown gracefully")
}
