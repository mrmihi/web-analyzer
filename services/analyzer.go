package services

import (
	"context"
	"github.com/go-rod/rod"
	"github.com/samber/lo"
	"net/url"
	"scraper/dto"
	"scraper/internal/logger"
	"scraper/internal/scraper"
	"sync"
)

// WebAnalysisService contains the business logic for analyzing a webpage.
// It holds its dependencies as interface types.
type WebAnalysisService struct {
	Analyzer rodAnalyzer.PageAnalyzer
}

func NewWebAnalysisService(analyzer rodAnalyzer.PageAnalyzer) *WebAnalysisService {
	return &WebAnalysisService{Analyzer: analyzer}
}

func (s *WebAnalysisService) AnalyseWebPage(ctx context.Context, targetUrl string) (dto.AnalyzeWebsiteRes, error) {
	var result dto.AnalyzeWebsiteRes

	handler := func(p *rodAnalyzer.ExtendedPage) error {
		result.HTMLVersion = p.HTMLVersion()
		result.PageTitle = p.MustInfo().Title
		result.HeadingCounts.H1 = p.ElementCount("h1")
		result.HeadingCounts.H2 = p.ElementCount("h2")
		result.HeadingCounts.H3 = p.ElementCount("h3")
		result.HeadingCounts.H4 = p.ElementCount("h4")
		result.HeadingCounts.H5 = p.ElementCount("h5")
		result.HeadingCounts.H6 = p.ElementCount("h6")
		result.ContainsLoginForm = p.ContainsLoginForm()

		analyzeLinks := func(pp rod.Pool[rod.Page]) {
			baseURL := lo.FromPtr(ok(url.Parse(p.MustInfo().URL)))
			allLinks := ok(p.Elements(`a[href]:not([href^="mailto:"]):not([href^="tel:"])`))

			var wg sync.WaitGroup
			for _, a := range allLinks {
				wg.Add(1)
				go func(el *rod.Element) {
					defer wg.Done()
					href := ok(el.Property("href")).String()
					if isExternal(href, &baseURL) {
						result.ExternalLinkCount++
					} else {
						result.InternalLinkCount++
					}
				}(a)
			}
			wg.Wait()
		}

		rodAnalyzer.RunWithNewPagePool(s.Analyzer.GetBrowser(), 5, analyzeLinks)
		return nil
	}

	err := s.Analyzer.Analyze(ctx, targetUrl, handler)
	if err != nil {
		return dto.AnalyzeWebsiteRes{}, err
	}

	return result, nil
}

func isExternal(link string, base *url.URL) bool {
	linkURL, err := url.Parse(link)
	if err != nil {
		return false
	}
	return linkURL.IsAbs() && linkURL.Hostname() != "" && linkURL.Hostname() != base.Hostname()
}

func ok[T any](v T, err error) T {
	if err != nil {
		logger.WarnCtx(context.Background(), "Ok function received an error", logger.Field{Key: "error", Value: err})
	}
	return v
}
