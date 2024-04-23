package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"2sem/hw1/homework/internal/config"
	"2sem/hw1/homework/internal/repo"
	"2sem/hw1/homework/internal/service"
	"2sem/hw1/homework/internal/transport"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
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
	chatRepo := repo.NewRepo(pool)
	chatService := service.NewChatService(chatRepo)
	chatHandler := transport.NewChatHandler(chatService, upgrader)

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

	runMigrations(connString)

	return pool
}

func runMigrations(connString string) {
	log.Printf("Run migrations on %s\n", connString)
	m, err := migrate.New("file://migrations", connString)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("Migrate no change")
	} else if err != nil {
		log.Fatal(err)
	}
	log.Println("Migrate ran successfully")
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
