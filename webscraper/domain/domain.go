package domain

type Headers struct {
	H1, H2, H3, H4, H5, H6 int
}

type Links struct {
	Internal, External, InAccessible, AllLinks, UniqueLinks int
}

type Result struct {
	HtmlVersion string
	PageTitle   string
	IsLoginPage bool
	Headers
	Links
}
