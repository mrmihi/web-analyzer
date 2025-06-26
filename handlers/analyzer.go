package handlers

import (
	"context"
	"fmt"
	"net/http"
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
	logger.InfoCtx(ctx, "Received request to analyze a webpage")

	url := c.Query("url")
	if url == "" {
		logger.InfoCtx(ctx, "Missing URL in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	logger.InfoCtx(ctx, "Analyzing webpage", logger.Field{Key: "url", Value: url})

	analysisCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	result, err := ac.AnalysisService.AnalyseWebPage(analysisCtx, url)
	if err != nil {
		logger.ErrorCtx(ctx, "Analysis failed", logger.Field{Key: "url", Value: url}, logger.Field{Key: "error", Value: err})
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to analyze the webpage: %v", err)})
		return
	}
	c.JSON(http.StatusOK, result)
}
