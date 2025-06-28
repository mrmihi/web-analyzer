package scraper

import (
	"context"
	"scraper/dto"
)

// PageAnalyzer defines the common interface for any web page analyzer.
type PageAnalyzer interface {
	Analyze(ctx context.Context, url string) (dto.AnalyzeWebsiteRes, error)
	Close() error
}
