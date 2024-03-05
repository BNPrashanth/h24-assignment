package models

import "net/http"

type IndexHTML struct {
	BaseURL string
}

type ParseRequestParams struct {
	URL string `json:"url"`
}

type ParseResponseHeading struct {
	Tag   string   `json:"tag"`
	Count int      `json:"count"`
	List  []string `json:"list"`
}

type ParseResponseHeadings struct {
	Count int                  `json:"count"`
	H1    ParseResponseHeading `json:"h1"`
	H2    ParseResponseHeading `json:"h2"`
	H3    ParseResponseHeading `json:"h3"`
	H4    ParseResponseHeading `json:"h4"`
	H5    ParseResponseHeading `json:"h5"`
	H6    ParseResponseHeading `json:"h6"`
}

type ParseResponseLinkInactive struct {
	Link   string `json:"link"`
	Status int    `json:"status"`
}

type ParseResponseLink struct {
	Count         int                         `json:"count"`
	InternalCount int                         `json:"internal_count"`
	InternalLinks []string                    `json:"internal_links"`
	ExternalCount int                         `json:"external_count"`
	ExternalLinks []string                    `json:"external_links"`
	InactiveCount int                         `json:"inactive_count"`
	InactiveLinks []ParseResponseLinkInactive `json:"inactive_links"`
}

type ParseResponse struct {
	HTMLVersion  string                `json:"html_version"`
	Title        string                `json:"title"`
	Headings     ParseResponseHeadings `json:"headings"`
	Links        ParseResponseLink     `json:"links"`
	HasLoginForm bool                  `json:"has_login_form"`
	Success      bool                  `json:"success"`
}

type ParseResponseError struct {
	Message    string `json:"message"`
	DevMessage string `json:"dev_message"`
	Status     int    `json:"status"`
	Success    bool   `json:"success"`
}

type LinksResponse struct {
	Link string
	Resp *http.Response
	Err  error
}
