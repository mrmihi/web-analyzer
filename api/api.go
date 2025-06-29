package api

import (
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"scraper/handlers"
	"time"
)

func AddMetricsRoutes(group *gin.RouterGroup) {
	group.GET("/system/metrics", gin.WrapH(promhttp.Handler()))
}

func AddAnalyzeRoutes(group *gin.RouterGroup, store persistence.CacheStore, ttl time.Duration, controller *handlers.AnalysisController) {
	group.GET("/analyze/", cache.CachePage(store, ttl, controller.Analyze))
}
