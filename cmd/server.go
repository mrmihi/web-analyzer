package cmd

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
	api2 "scraper/api"
	"scraper/common"
	"scraper/config"
	"scraper/handlers"
	"scraper/middleware"
	"time"
)

// NewServer creates a new HTTP server with the configured routes and middleware.
func NewServer(analysisController *handlers.AnalysisController) *http.Server {
	router := gin.Default()

	router.Use(otelgin.Middleware(common.ServiceName))
	router.Use(middleware.PrometheusMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LeakBucket())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	v1 := api.Group("/v1")

	api2.AddAnalyzerRoutes(v1, analysisController)
	api2.AddMetricsRoutes(v1)

	server := &http.Server{
		Addr:         config.Config.Host + ":" + config.Config.Port,
		Handler:      router,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		IdleTimeout:  1 * time.Minute,
	}
	return server
}
