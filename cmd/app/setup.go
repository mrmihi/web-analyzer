package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"scraper/cmd"
	"scraper/config"
	"scraper/handlers"
	"scraper/internal/logger"
	"scraper/internal/scraper"
	"scraper/services"
)

// App holds all the core components of the application
type App struct {
	Config          *config.Cfg
	Logger          logger.Logger
	Analyzer        *rodAnalyzer.RodAnalyzer
	AnalysisService *services.WebAnalysisService
	Server          *http.Server
}

// New creates and wires up all the application components.
func New() (*App, func()) {

	appLogger := logger.NewZapLogger(true)
	logger.SetLogger(appLogger)

	appConfig := config.GetConfig()

	analyzer, err := rodAnalyzer.NewRodAnalyzer()
	if err != nil {
		log.Fatalf("FATAL: Failed to create rod analyzer: %s\n", err)
	}

	analysisService := services.NewWebAnalysisService(analyzer)

	analysisController := handlers.NewAnalysisController(analysisService)

	server := cmd.NewServer(analysisController)

	app := &App{
		Config:          appConfig,
		Logger:          appLogger,
		Analyzer:        analyzer,
		AnalysisService: analysisService,
		Server:          server,
	}

	cleanup := func() {
		fmt.Println("Running cleanup tasks...")
		if err := app.Analyzer.Close(); err != nil {
			app.Logger.ErrorCtx(context.Background(), "Error closing analyzer", logger.Field{Key: "error", Value: err})
		}
	}

	return app, cleanup
}
