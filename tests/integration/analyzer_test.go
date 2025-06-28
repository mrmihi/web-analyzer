package integration

import (
	"context"
	"os"
	"path/filepath"
	"scraper/dto"
	"scraper/internal/scraper/rodAnalyzer"
	"scraper/services"
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	. "github.com/smartystreets/goconvey/convey"
)

// MockAnalyzer implements the rodAnalyzer.PageAnalyzer interface
type MockAnalyzer struct {
	browser *rod.Browser
}

// New creates a new MockAnalyzer
func NewMockAnalyzer() (*MockAnalyzer, error) {
	u := launcher.New().Headless(true).Leakless(false).NoSandbox(true).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()

	return &MockAnalyzer{
		browser: browser,
	}, nil
}

// Analyze implements the PageAnalyzer interface
func (m *MockAnalyzer) Analyze(context.Context, string) (dto.AnalyzeWebsiteRes, error) {
	page := m.browser.MustPage("")
	defer func(page *rod.Page) {
		err := page.Close()
		if err != nil {

		}
	}(page)

	mockPath, err := filepath.Abs("mocks/sample.html")
	if err != nil {
		return dto.AnalyzeWebsiteRes{}, err
	}

	htmlContent, err := os.ReadFile(mockPath)
	if err != nil {
		return dto.AnalyzeWebsiteRes{}, err
	}

	page.MustSetDocumentContent(string(htmlContent))

	extendedPage := &rodAnalyzer.ExtendedPage{Page: page}

	result := dto.AnalyzeWebsiteRes{}
	result.HTMLVersion = extendedPage.HTMLVersion()
	result.Title = extendedPage.MustInfo().Title
	result.Headings.H1 = extendedPage.ElementCount("h1")
	result.Headings.H2 = extendedPage.ElementCount("h2")
	result.Headings.H3 = extendedPage.ElementCount("h3")
	result.Headings.H4 = extendedPage.ElementCount("h4")
	result.Headings.H5 = extendedPage.ElementCount("h5")
	result.Headings.H6 = extendedPage.ElementCount("h6")
	result.LoginForm = extendedPage.ContainsLoginForm()

	result.InternalLinks = 2
	result.ExternalLinks = 1
	result.InaccessibleLinks = 1

	return result, nil
}

// Close implements the PageAnalyzer interface
func (m *MockAnalyzer) Close() error {
	return m.browser.Close()
}

func TestWebAnalyzer(t *testing.T) {
	Convey("Given a web analyzer service with a mock analyzer", t, func() {
		mockAnalyzer, err := NewMockAnalyzer()
		So(err, ShouldBeNil)
		defer func(mockAnalyzer *MockAnalyzer) {
			err := mockAnalyzer.Close()
			if err != nil {
			}
		}(mockAnalyzer)

		service := services.NewWebAnalysisService(mockAnalyzer)

		Convey("When analyzing a mock webpage", func() {
			result, err := service.AnalyseWebPage(context.Background(), "mock-url")

			Convey("Then the analysis should complete without errors", func() {
				So(err, ShouldBeNil)

				Convey("And the HTML version should be detected correctly", func() {
					So(result.HTMLVersion, ShouldEqual, "HTML5")
				})

				Convey("And the title should be extracted correctly", func() {
					So(result.Title, ShouldEqual, "Sample Page for Testing")
				})

				Convey("And the heading counts should be correct", func() {
					So(result.Headings.H1, ShouldEqual, 1)
					So(result.Headings.H2, ShouldEqual, 2)
					So(result.Headings.H3, ShouldEqual, 2)
					So(result.Headings.H4, ShouldEqual, 1)
					So(result.Headings.H5, ShouldEqual, 1)
					So(result.Headings.H6, ShouldEqual, 1)
				})

				Convey("And the link counts should be correct", func() {
					So(result.InternalLinks, ShouldEqual, 2)
					So(result.ExternalLinks, ShouldEqual, 1)
					So(result.InaccessibleLinks, ShouldEqual, 1)
				})

				Convey("And the login form detection should be correct", func() {
					So(result.LoginForm, ShouldBeTrue)
				})
			})
		})
	})
}
