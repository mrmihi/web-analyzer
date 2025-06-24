package routes

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
	"scraper/consts"
	"scraper/middleware"
	"time"
)

// NewServer creates a new HTTP server with the configured routes and middleware.
func NewServer() *http.Server {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	router.Use(otelgin.Middleware(consts.ServiceName))
	router.Use(middleware.PrometheusMiddleware())
	router.Use(middleware.LeakBucket())
	// TODO: gin-helmet installation issues
	//router.Use(ginhelmet.Default())

	api := router.Group("/api")
	v1 := api.Group("/v1")

	AddAnalyzerRoutes(v1)
	AddMetricsRoutes(v1)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		IdleTimeout:  1 * time.Minute,
	}
	return server
}
