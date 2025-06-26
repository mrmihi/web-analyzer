package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func AddMetricsRoutes(group *gin.RouterGroup) {
	group.GET("/system/metrics", gin.WrapH(promhttp.Handler()))
}
