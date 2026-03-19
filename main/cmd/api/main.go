package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	sqlcdb "main-private/main/internal/infrastructure/db/sqlc"
	"main-private/main/internal/infrastructure/repository"
	itemuc "main-private/main/internal/usecase/item"
	statsuc "main-private/main/internal/usecase/stats"
	"main-private/main/pkg/config"
	"main-private/main/pkg/logger"
	"main-private/main/pkg/server"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Env)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Error("failed to ping database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("connected to database")

	queries := sqlcdb.New(pool)
	itemRepo := repository.NewItemRepository(queries)
	itemInteractor := itemuc.NewInteractor(itemRepo)

	statsRepo := repository.NewStatsRepository(queries)
	statsInteractor := statsuc.NewInteractor(statsRepo)

	handler := server.New(log, itemInteractor, statsInteractor)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("server starting", slog.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("server stopped")
}
