package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"scraper/cmd/app"
	"scraper/config"
	"scraper/internal/logger"
	"syscall"
	"time"
)

func main() {

	stopCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctx := context.Background()

	service, cleanup := app.New()

	defer cleanup()

	go func() {
		logger.InfoCtx(ctx, fmt.Sprintf("Starting server on port %d", config.Config.Port))

		err := service.Server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-stopCtx.Done()

	stop()

	logger.InfoCtx(ctx, "Shutdown signal received, starting graceful shutdown!")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()

	if err := service.Server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	logger.InfoCtx(ctx, "Server exited gracefully")
}
