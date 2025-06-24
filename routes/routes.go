package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"scraper/handlers"
)

var metricsHandler = promhttp.Handler()

func AddAnalyzerRoutes(group *gin.RouterGroup) {
	group.POST("/", handlers.AnalyzerHandler)
}

func AddMetricsRoutes(group *gin.RouterGroup) {
	group.GET("/system/metrics", gin.WrapH(metricsHandler))
}
