package handlers

import (
	"context"
	"fmt"
	"net/http"
	"scraper/dto"
	"scraper/internal/logger"
	"scraper/internal/scraper"
	"scraper/services"
	"time"

	"github.com/gin-gonic/gin"
)

// AnalyzerHandler handles requests to the root path ("/")
func AnalyzerHandler(c *gin.Context) {
	ctx := context.Background()

	analyzer, err := rodAnalyzer.NewRodAnalyzer()
	if err != nil {
		logger.ErrorCtx(context.Background(), "Failed to create analyzer", logger.Field{Key: "error", Value: err})
		return
	}
	defer func(analyzer *rodAnalyzer.RodAnalyzer) {
		err := analyzer.Close()
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to close analyzer", logger.Field{Key: "error", Value: err})
		}
	}(analyzer)

	analysisService := services.NewWebAnalysisService(analyzer)

	var request dto.AnalyzeWebsiteReq
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.ErrorCtx(ctx, "Invalid request body", logger.Field{Key: "error", Value: err})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body. Please provide a valid URL in the request body."})
		return
	}

	if request.URL == "" {
		logger.InfoCtx(ctx, "Missing URL in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	logger.InfoCtx(ctx, "Analyzing webpage", logger.Field{Key: "url", Value: request.URL})

	timer := time.NewTimer(120 * time.Second)

	done := make(chan dto.AnalyzeWebsiteRes, 1)
	errChan := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- fmt.Errorf("internal server error: %v", r)
			}
		}()

		result, _ := analysisService.AnalyseWebPage(context.Background(), request.URL)
		done <- result
	}()

	select {
	case result := <-done:
		timer.Stop()
		c.JSON(http.StatusOK, result)
	case err := <-errChan:
		timer.Stop()
		logger.ErrorCtx(ctx, "Analysis failed", logger.Field{Key: "url", Value: request.URL}, logger.Field{Key: "error", Value: err})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case <-timer.C:
		logger.ErrorCtx(ctx, "Request timed out", logger.Field{Key: "url", Value: request.URL})
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out after 120 seconds. The analysis is taking longer than expected."})
	}
}
