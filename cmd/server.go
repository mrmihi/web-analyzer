package cmd

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"scraper/common"
	"scraper/middleware"
)

// NewRouter initializes the Gin router with all necessary middleware and routes.
func NewRouter() *gin.Engine {
	router := gin.Default()

	router.Use(otelgin.Middleware(common.ServiceName))
	router.Use(middleware.PrometheusMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LeakBucket())
	router.Use(gin.Recovery())

	return router
}
