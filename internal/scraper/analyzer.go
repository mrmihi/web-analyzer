package rodAnalyzer

import (
	"context"
	"errors"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"io"
	"net/http"
	"net/url"
	"scraper/common"
	"scraper/config"
	"scraper/dto"
	"scraper/internal/logger"
	"sync"
)

var (
	ErrConnectionError = errors.New("failed to connect to the webpage")
)

// PageAnalyzer defines the interface for analyzing web pages.
type PageAnalyzer interface {
	Analyze(ctx context.Context, url string) (dto.AnalyzeWebsiteRes, error)
	Close() error
}

// RodAnalyzer is the concrete implementation of PageAnalyzer using the rod library.
type RodAnalyzer struct {
	Browser *rod.Browser
}

// New creates and configures a new rod-based analyzer.
func New() (*RodAnalyzer, error) {
	var l *launcher.Launcher
	if config.Config.ChromeSetup != "" {
		l = launcher.New().Bin(config.Config.ChromeSetup)
	} else {
		path, exists := launcher.LookPath()
		if !exists {
			return nil, errors.New("cannot find a browser binary")
		}
		l = launcher.New().Bin(path)
	}

	u := l.Headless(true).NoSandbox(true).Leakless(false).Set("no-sandbox").Set("disable-gpu").MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()

	router := browser.HijackRequests()
	router.MustAdd("*", func(ctx *rod.Hijack) {
		switch ctx.Request.Type() {
		case proto.NetworkResourceTypeImage,
			proto.NetworkResourceTypeStylesheet,
			proto.NetworkResourceTypeFont,
			proto.NetworkResourceTypeMedia,
			proto.NetworkResourceTypeTextTrack,
			proto.NetworkResourceTypeManifest,
			proto.NetworkResourceTypeEventSource,
			proto.NetworkResourceTypeWebSocket:
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
		default:
			ctx.ContinueRequest(&proto.FetchContinueRequest{})
		}
	})
	go router.Run()

	return &RodAnalyzer{Browser: browser}, nil
}

// Analyze fetches and analyzes the webpage at the given URL, returning the analysis results.
func (r *RodAnalyzer) Analyze(ctx context.Context, targetUrl string) (dto.AnalyzeWebsiteRes, error) {
	logger.InfoCtx(ctx, "Visiting page", logger.Field{Key: "url", Value: targetUrl})

	var result dto.AnalyzeWebsiteRes
	var e proto.NetworkResponseReceived

	page := r.Browser.MustPage("").Context(ctx)
	defer func() {
		if err := page.Close(); err != nil {
			logger.ErrorCtx(ctx, "Failed to close page", logger.Field{Key: "error", Value: err})
		}
	}()

	wait := page.WaitEvent(&e)
	if err := page.Navigate(targetUrl); err != nil {
		logger.ErrorCtx(ctx, "Failed to retrieve webpage", logger.Field{Key: "error", Value: err})
		return result, ErrConnectionError
	}
	wait()
	page.MustWaitLoad()

	if e.Response.Status < 200 || e.Response.Status >= 300 {
		logger.ErrorCtx(ctx, "Invalid response status", logger.Field{Key: "status", Value: e.Response.Status})
		return result, common.NewGinError(422, "Failed to analyze webpage", e.Response.Status)
	}

	extendedPage := &ExtendedPage{page}

	result.HTMLVersion = extendedPage.HTMLVersion()
	result.Title = extendedPage.MustInfo().Title
	result.Headings.H1 = extendedPage.ElementCount("h1")
	result.Headings.H2 = extendedPage.ElementCount("h2")
	result.Headings.H3 = extendedPage.ElementCount("h3")
	result.Headings.H4 = extendedPage.ElementCount("h4")
	result.Headings.H5 = extendedPage.ElementCount("h5")
	result.Headings.H6 = extendedPage.ElementCount("h6")
	result.LoginForm = extendedPage.ContainsLoginForm()

	baseURL, err := url.Parse(extendedPage.MustInfo().URL)
	if err != nil {
		logger.WarnCtx(ctx, "Could not parse base URL", logger.Field{Key: "error", Value: err})
		return result, err
	}

	allLinkElements, err := extendedPage.Elements(`a[href]:not([href^="mailto:"]):not([href^="tel:"])`)
	if err != nil {
		logger.WarnCtx(ctx, "Could not get link elements", logger.Field{Key: "error", Value: err})
		return result, err
	}

	var wg sync.WaitGroup
	for _, a := range allLinkElements {
		wg.Add(1)
		go func(el *rod.Element) {
			defer wg.Done()
			href, _ := el.Property("href")
			if !isLinkAccessible(href.String()) {
				logger.InfoCtx(ctx, "Link is inaccessible", logger.Field{Key: "link", Value: href.String()})
				result.InaccessibleLinks++
			} else if isExternal(href.String(), baseURL) {
				result.ExternalLinks++
			} else {
				result.InternalLinks++
			}
		}(a)
	}
	wg.Wait()

	return result, nil
}

// Close cleans up the browser instance.
func (r *RodAnalyzer) Close() error {
	return r.Browser.Close()
}

// isExternal is a helper function to check if a link is external.
func isExternal(link string, base *url.URL) bool {
	linkURL, err := url.Parse(link)
	if err != nil {
		return false // Or handle as an inaccessible link
	}
	return linkURL.IsAbs() && linkURL.Hostname() != "" && linkURL.Hostname() != base.Hostname()
}

func isLinkAccessible(link string) bool {
	resp, err := http.Head(link)
	if err != nil {
		logger.InfoCtx(context.Background(), "Failed to check link accessibility", logger.Field{Key: "link", Value: link})
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.ErrorCtx(context.Background(), "Failed to close response body", logger.Field{Key: "error", Value: err})
		}
	}(resp.Body)
	return true
}
