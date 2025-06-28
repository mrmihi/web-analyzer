package app

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cache/persistence"
	"log"
	"net/http"
	"scraper/api"
	"scraper/cmd"
	"scraper/config"
	"scraper/handlers"
	"scraper/internal/logger"
	"scraper/internal/scraper"
	"scraper/internal/scraper/htmlAnalyzer"
	"scraper/internal/scraper/rodAnalyzer"
	"scraper/services"
	"time"
)

var Service *App

// App holds all the core components of the application
type App struct {
	Config             *config.Cfg
	Logger             logger.Logger
	AnalysisController *handlers.AnalysisController
	Server             *http.Server
}

// New creates and wires up all the application components.
func New() (*App, func()) {
	ctx := context.Background()
	var analyzer scraper.PageAnalyzer
	var err error
	appLogger := logger.NewZapLogger(true)
	logger.SetLogger(appLogger)

	appConfig := config.GetConfig()

	switch appConfig.AnalyzerType {
	case "rod":
		analyzer, err = rodAnalyzer.New()
		if err != nil {
			log.Fatalf("FATAL: Failed to create rod analyzer: %s\n", err)
		}
		appLogger.InfoCtx(ctx, "Using 'rod' page analyzer.")
	case "html":
		analyzer, err = htmlAnalyzer.New()
		if err != nil {
			log.Fatalf("FATAL: Failed to create html parser: %s\n", err)
		}
		appLogger.InfoCtx(ctx, "Using 'html' page analyzer.")
	default:
		log.Fatalf("FATAL: Invalid analyzer type specified: %s\n", appConfig.AnalyzerType)
	}

	analysisService := services.NewWebAnalysisService(analyzer)

	analysisController := handlers.NewAnalysisController(analysisService)

	router := cmd.NewRouter()
	cacheStore := persistence.NewInMemoryStore(appConfig.InMemStoreTTL * time.Minute)

	apiGroup := router.Group("/api")
	v1 := apiGroup.Group("/v1")

	api.AddAnalyzeRoutes(v1, cacheStore, appConfig.InMemStoreTTL*time.Minute, analysisController)
	api.AddMetricsRoutes(v1)

	server := &http.Server{
		Addr:    appConfig.Host + ":" + appConfig.Port,
		Handler: router,
	}

	Service = &App{
		Config:             appConfig,
		Logger:             appLogger,
		AnalysisController: analysisController,
		Server:             server,
	}

	cleanup := func() {
		fmt.Println("Running cleanup tasks...")
		if err := analyzer.Close(); err != nil {
			Service.Logger.ErrorCtx(context.Background(), "Error closing analyzer", logger.Field{Key: "error", Value: err})
		}
	}

	return Service, cleanup
}
