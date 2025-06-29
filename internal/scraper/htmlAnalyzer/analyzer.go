package htmlAnalyzer

import (
	"context"
	"scraper/dto"
	"scraper/internal/logger"
)

type HTMLParse struct {
}

func New() (*HTMLParse, error) {
	return &HTMLParse{}, nil
}

func (r *HTMLParse) Analyze(ctx context.Context, targetUrl string) (dto.AnalyzeWebsiteRes, error) {
	logger.InfoCtx(ctx, "Visiting page", logger.Field{Key: "url", Value: targetUrl})
	return dto.AnalyzeWebsiteRes{}, nil
}

func (r *HTMLParse) Close() error {
	return nil
}
