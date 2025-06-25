package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"scraper/handlers"
)

func AddAnalyzerRoutes(group *gin.RouterGroup) {
	group.GET("/analyze", handlers.AnalyzerHandler)
}

func AddMetricsRoutes(group *gin.RouterGroup) {
	group.GET("/system/metrics", gin.WrapH(promhttp.Handler()))
}
