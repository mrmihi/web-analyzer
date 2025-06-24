package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"scraper/config"
	"scraper/consts"
	"scraper/internal/logger"
	"scraper/routes"
	"syscall"
	"time"
)

type Service struct {
	Name   string
	Server *http.Server
	Logger logger.Logger
	Config *config.Config
}

func NewService(serviceName string) *Service {
	zapLogger := logger.NewZapLogger(true)
	logger.SetLogger(zapLogger)
	return &Service{
		Name:   serviceName,
		Server: routes.NewServer(),
		Logger: logger.GetLogger(),
		Config: config.GetConfig(),
	}
}

func main() {

	stopCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctx := context.Background()

	service := NewService(consts.ServiceName)

	go func() {
		logger.InfoCtx(ctx, fmt.Sprintf("Starting server on port %d", config.Env.Port))

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
