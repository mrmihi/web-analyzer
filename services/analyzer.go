package services

import (
	"context"
	"scraper/dto"
	"scraper/internal/scraper"
)

// WebAnalysisService contains the business logic for analyzing a webpage.
type WebAnalysisService struct {
	Analyzer rodAnalyzer.PageAnalyzer
}

// NewWebAnalysisService creates a new WebAnalysisService.
func NewWebAnalysisService(analyzer rodAnalyzer.PageAnalyzer) *WebAnalysisService {
	return &WebAnalysisService{Analyzer: analyzer}
}

// AnalyseWebPage performs the analysis of a web page given its URL.
func (s *WebAnalysisService) AnalyseWebPage(ctx context.Context, targetUrl string) (dto.AnalyzeWebsiteRes, error) {
	return s.Analyzer.Analyze(ctx, targetUrl)
}
