package rodAnalyzer

import (
	"context"
	"github.com/go-rod/rod"
	"io"
	"net/http"
	"net/url"
	"scraper/internal/logger"
	"strings"
)

type ExtendedPage struct {
	*rod.Page
}

func (ep *ExtendedPage) ElementCount(selector string) int {
	if selector == "h1" || selector == "h2" || selector == "h3" || selector == "h4" || selector == "h5" || selector == "h6" {
		result, err := ep.Eval(`(selector) => {
			return document.getElementsByTagName(selector).length;
		}`, selector)

		if err != nil {
			return 0
		}

		return result.Value.Int()
	}

	result, err := ep.Eval(`(selector) => {
		return document.querySelectorAll(selector).length;
	}`, selector)

	if err != nil {
		return 0
	}

	return result.Value.Int()
}

func (ep *ExtendedPage) ContainsLoginForm() bool {
	return ep.MustHas("form input[type=password]")
}

func (ep *ExtendedPage) HTMLVersion() string {
	doctype, err := ep.Eval(`() => document.doctype ? new XMLSerializer().serializeToString(document.doctype) : ''`)
	if err != nil || doctype.Value.Str() == "" {
		return "Unknown"
	}
	if strings.Contains(strings.ToLower(doctype.Value.Str()), "html") {
		return "HTML5"
	}
	return "HTML4 or older"
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
