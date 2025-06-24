package dto

// TODO: Refactor this file to use the global response structure

type AnalyzeWebsiteReq struct {
	URL string `json:"url" validate:"required,url" messages:"Please provide a valid url to analyse"`
}

type AnalyzeWebsiteRes struct {
	HTMLVersion           string        `json:"html_version"`
	PageTitle             string        `json:"page_title"`
	HeadingCounts         HeadingCounts `json:"heading_counts"`
	InternalLinkCount     int           `json:"internal_link_count"`
	ExternalLinkCount     int           `json:"external_link_count"`
	InaccessibleLinkCount int           `json:"inaccessible_link_count"`
	ContainsLoginForm     bool          `json:"contains_login_form"`
}

type HeadingCounts struct {
	H1 int `json:"h1"`
	H2 int `json:"h2"`
	H3 int `json:"h3"`
	H4 int `json:"h4"`
	H5 int `json:"h5"`
	H6 int `json:"h6"`
}
