package app

import (
	"context"
	"fmt"
	"github.com/dzhordano/avito-bootcamp2024/internal/config"
	"github.com/dzhordano/avito-bootcamp2024/internal/delivery/http"
	"github.com/dzhordano/avito-bootcamp2024/internal/repository"
	"github.com/dzhordano/avito-bootcamp2024/internal/server"
	"github.com/dzhordano/avito-bootcamp2024/internal/service"
	"github.com/dzhordano/avito-bootcamp2024/pkg/auth"
	"github.com/dzhordano/avito-bootcamp2024/pkg/databases/postgres"
	"github.com/dzhordano/avito-bootcamp2024/pkg/logger"
	"github.com/dzhordano/avito-bootcamp2024/pkg/notifications/sender"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Run initializes whole application.
func Run() {
	cfg := config.MustLoad()

	log := logger.NewLogger("debug")

	tokenManager := auth.NewJWTManager(cfg.Auth.SecretKey, cfg.Auth.TokenTTL)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
		cfg.Postgres.SSLMode,
	)

	pool, err := postgres.NewClient(dsn)
	if err != nil {
		log.Error("failed to connect to postgres: " + err.Error())

		return
	}
	defer pool.Close()

	notificationSender := sender.New()

	toWaitTasks := &sync.WaitGroup{}

	repo := repository.New(pool)
	svc := service.New(service.Deps{
		Repos:         repo,
		TokensManager: tokenManager,
		Notifications: notificationSender,
		WaitGroup:     toWaitTasks,
		Logger:        log,
	})

	handler := http.NewHandler(svc, tokenManager)
	srv := server.NewServer(cfg, handler.Init())

	go func() {
		if err := srv.Run(); err != nil {
			log.Error(err.Error())
		}
	}()

	toWaitTasks.Wait()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server: " + err.Error())

		return
	}
}
