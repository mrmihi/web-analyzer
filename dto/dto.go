package dto

type AnalyzeWebsiteReq struct {
	URL string `json:"url" validate:"required,url" messages:"Please provide a valid url to analyse"`
}

type AnalyzeWebsiteRes struct {
	HTMLVersion       string   `json:"html_version"`
	Title             string   `json:"title"`
	Headings          Headings `json:"headings"`
	InternalLinks     int      `json:"internal_links"`
	ExternalLinks     int      `json:"external_links"`
	InaccessibleLinks int      `json:"inaccessible_links"`
	LoginForm         bool     `json:"login_form"`
}

type Headings struct {
	H1 int `json:"h1"`
	H2 int `json:"h2"`
	H3 int `json:"h3"`
	H4 int `json:"h4"`
	H5 int `json:"h5"`
	H6 int `json:"h6"`
}
