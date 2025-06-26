package rodAnalyzer

import (
	"github.com/go-rod/rod"
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
