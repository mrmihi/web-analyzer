package rodAnalyzer

import (
	"context"
	"errors"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"scraper/common"
	"scraper/config"
	"scraper/internal/logger"
)

var (
	ErrConnectionError      = errors.New("failed to connect to the webpage")
	ErrTargetIsNotValidHTML = errors.New("target resource is not a valid HTML document")
)

// PageAnalyzer is an interface that defines the methods for analyzing web pages.
type PageAnalyzer interface {
	Analyze(ctx context.Context, url string, handler func(p *ExtendedPage) error) error
	GetBrowser() *rod.Browser
	Close() error
}

// RodAnalyzer is the concrete implementation of PageAnalyzer using the rod library.
type RodAnalyzer struct {
	Browser *rod.Browser
}

// NewRodAnalyzer creates and configures a new rod-based analyzer.
func NewRodAnalyzer() (*RodAnalyzer, error) {
	var l *launcher.Launcher
	if config.Env.ChromeSetup != "" {
		l = launcher.New().Bin(config.Env.ChromeSetup)
	} else {
		path, exists := launcher.LookPath()
		if !exists {
			return nil, errors.New("cannot find a browser binary")
		}
		l = launcher.New().Bin(path)
	}

	u := l.Headless(true).NoSandbox(true).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()

	router := browser.HijackRequests()
	router.MustAdd("*", func(ctx *rod.Hijack) {
		switch ctx.Request.Type() {
		case proto.NetworkResourceTypeImage,
			proto.NetworkResourceTypeStylesheet,
			proto.NetworkResourceTypeFont,
			proto.NetworkResourceTypeMedia:
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
		default:
			ctx.ContinueRequest(&proto.FetchContinueRequest{})
		}
	})
	go router.Run()

	return &RodAnalyzer{Browser: browser}, nil
}

// Analyze implements the PageAnalyzer interface. It encapsulates the session logic.
func (r *RodAnalyzer) Analyze(ctx context.Context, url string, handler func(p *ExtendedPage) error) error {
	logger.InfoCtx(ctx, "Visiting page", logger.Field{Key: "url", Value: url})

	var e proto.NetworkResponseReceived
	page := r.Browser.MustPage("")
	defer func(page *rod.Page) {
		err := page.Close()
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to close page", logger.Field{Key: "error", Value: err})
		}
	}(page)

	page = page.Context(ctx)

	wait := page.WaitEvent(&e)
	if err := page.Navigate(url); err != nil {
		logger.ErrorCtx(ctx, "Failed to retrieve webpage", logger.Field{Key: "error", Value: err})
		return ErrConnectionError
	}
	wait()

	page.MustWaitLoad()

	//contentType, ok := e.Response.Headers["Content-Type"]
	//if !ok || !strings.Contains(contentType.String(), "text/html") {
	//	logger.ErrorCtx(ctx, "Invalid content type", logger.Field{Key: "content-type", Value: contentType})
	//	return ErrTargetIsNotValidHTML
	//}

	if e.Response.Status < 200 || e.Response.Status >= 300 {
		logger.ErrorCtx(ctx, "Invalid response status", logger.Field{Key: "status", Value: e.Response.Status})
		return common.NewGinError(422, "Failed to analyze webpage", e.Response.Status)
	}

	return handler(&ExtendedPage{page})
}

// GetBrowser provides access to the underlying browser for pool creation.
func (r *RodAnalyzer) GetBrowser() *rod.Browser {
	return r.Browser
}

// Close cleans up the browser instance.
func (r *RodAnalyzer) Close() error {
	return r.Browser.Close()
}
