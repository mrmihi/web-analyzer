package rodAnalyzer

import (
	"github.com/go-rod/rod"
	"strings"
)

type ExtendedPage struct {
	*rod.Page
}

func (ep *ExtendedPage) ElementCount(selector string) int {
	elements, err := ep.Elements(selector)
	if err != nil {
		return 0
	}
	return len(elements)
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

func RunWithNewPagePool(b *rod.Browser, limit int, fn func(rod.Pool[rod.Page])) {
	pool := rod.NewPagePool(limit)
	defer pool.Cleanup(func(p *rod.Page) { p.Close() })
	fn(pool)
}
