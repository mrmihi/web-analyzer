package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"scraper/dto"
	"scraper/handlers"
	"scraper/services"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAnalyzerHandler(t *testing.T) {
	Convey("Given an analyzer handler with a mock service", t, func() {
		mockAnalyzer, err := NewMockAnalyzer()
		So(err, ShouldBeNil)
		defer func(mockAnalyzer *MockAnalyzer) {
			err := mockAnalyzer.Close()
			if err != nil {
			}
		}(mockAnalyzer)

		service := services.NewWebAnalysisService(mockAnalyzer)

		handler := handlers.NewAnalysisController(service)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/analyze", handler.Analyze)

		Convey("When sending a valid request", func() {
			url := "https://example.com"
			req, err := http.NewRequest("GET", "/analyze/?url="+url, nil)
			So(err, ShouldBeNil)

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			Convey("Then the response should be successful", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)

				Convey("And the response body should contain the analysis results", func() {
					var result dto.AnalyzeWebsiteRes
					err := json.Unmarshal(resp.Body.Bytes(), &result)
					So(err, ShouldBeNil)

					So(result.HTMLVersion, ShouldEqual, "HTML5")
					So(result.Title, ShouldEqual, "Sample Page for Testing")
					So(result.Headings.H1, ShouldEqual, 1)
					So(result.Headings.H2, ShouldEqual, 2)
					So(result.Headings.H3, ShouldEqual, 2)
					So(result.Headings.H4, ShouldEqual, 1)
					So(result.Headings.H5, ShouldEqual, 1)
					So(result.Headings.H6, ShouldEqual, 1)
					So(result.InternalLinks, ShouldEqual, 2)
					So(result.ExternalLinks, ShouldEqual, 1)
					So(result.InaccessibleLinks, ShouldEqual, 1)
					So(result.LoginForm, ShouldBeTrue)
				})
			})
		})

		Convey("When sending a request with an empty URL", func() {
			req, err := http.NewRequest("GET", "/analyze/", nil)
			So(err, ShouldBeNil)

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			Convey("Then the response should be a bad request", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)

				Convey("And the response body should contain an error message", func() {
					var errorResponse map[string]string
					err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					So(err, ShouldBeNil)
					So(errorResponse["error"], ShouldEqual, "URL is required")
				})
			})
		})

		Convey("When sending a request with an invalid URL format", func() {
			invalidURL := "invalid-url-format"

			req, err := http.NewRequest("GET", "/analyze/?url="+invalidURL, nil)
			So(err, ShouldBeNil)

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			Convey("Then the response should be a bad request", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)

				Convey("And the response body should contain an error message", func() {
					var errorResponse map[string]string
					err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					So(err, ShouldBeNil)
					So(errorResponse["error"], ShouldContainSubstring, "Invalid URL format")
				})
			})
		})
	})
}
