package handlers

import (
	"context"
	"net/http"
	"scraper/common"
	"scraper/config"
	"scraper/internal/logger"
	"scraper/services"
	"time"

	"github.com/gin-gonic/gin"
)

// AnalysisController holds the dependencies for the analysis handlers.
type AnalysisController struct {
	AnalysisService *services.WebAnalysisService
}

// NewAnalysisController creates a new handler with its dependencies.
func NewAnalysisController(service *services.WebAnalysisService) *AnalysisController {
	return &AnalysisController{
		AnalysisService: service,
	}
}

func (ac *AnalysisController) Analyze(c *gin.Context) {
	ctx := c.Request.Context()
	logger.InfoCtx(ctx, "Received a request to analyze a webpage")

	url := c.Query("url")
	if url == "" {
		logger.InfoCtx(ctx, "Missing URL in the request")
		urlError := common.NewGinError(common.RequestFail, "URL is required", nil)
		c.JSON(http.StatusBadRequest, urlError)
		return
	}

	logger.InfoCtx(ctx, "Analyzing webpage", logger.Field{Key: "url", Value: url})

	analysisCtx, cancel := context.WithTimeout(ctx, config.Config.AnalyzeTimeOut*time.Minute)
	defer cancel()

	result, err := ac.AnalysisService.AnalyseWebPage(analysisCtx, url)
	if err != nil {
		logger.ErrorCtx(ctx, "Analysis failed", logger.Field{Key: "url", Value: url}, logger.Field{Key: "error", Value: err})
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
