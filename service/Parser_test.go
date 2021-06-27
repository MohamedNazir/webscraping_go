package service

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/MohamedNazir/webscraper/domain"
	"github.com/MohamedNazir/webscraper/utils/mocks"
	"github.com/MohamedNazir/webscraper/webclient"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

var (
	HTML_FIRST_INPUT                  = `<!doctype html> <html> <head> <title>The Go Playground</title> <link rel="stylesheet" href="/static/style.css"> </head> <body itemscope itemtype="http://schema.org/CreativeWork"> <input type="button" value="Run" id="embedRun"> <div id="banner"> <h1> Sample text </h1> <h2> Sample text </h2> <h3> Sample text </h3> <h4> Sample text </h4> <h5> Sample text </h5> <h6> Sample text </h6> <input type="password" name="password" id ="test" > <input type="button" value="About" id="aboutButton"> </div> <div id="output"></div> <img itemprop="image" src="/static/gopher.png" style="display:none"> <div id="about"> <p><b>About the Playground</b></p> <a href="http://playground"> external link </a> <a href="/pages"> Internal link </a> </p> </div> </body> </html>`
	HTML_SECOND_INPUT                 = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd"> <html> <head> <title>The Test Case 2</title> </head> <body> <div id="banner"> <h1> Sample text </h1> <h3> Sample text </h3> <h5> Sample text </h5> <input type="button" value="About" id="aboutButton"> </div> <img itemprop="image" src="/static/gopher.png" style="display:none"> <div id="about"> <p><b>About the Playground</b></p> <a href="http://playground"> external link </a> <a href="/pages"> Internal link </a> <a href="/"> Broken link 1 </a> <a href="#"> Broken link 2 </a>  </p><p><a href="mailto:someone@example.com">Send email</a></p> </div> </body> </html>`
	DOC_ONE           *html.Tokenizer = html.NewTokenizer(strings.NewReader(HTML_FIRST_INPUT))
	DOC_TWO           *html.Tokenizer = html.NewTokenizer(strings.NewReader(HTML_SECOND_INPUT))
	headers_1                         = domain.Headers{
		H1: 1, H2: 1, H3: 1, H4: 1, H5: 1, H6: 1,
	}
	AllLinks_2 = []string{"http://playground", "http://dummy.com/2/pages", "http://dummy.com/2/", "http://dummy.com/2/#", "mailto:someone@example.com"}
	AllLinks_1 = []string{"http://playground", "http://dummy.com/1/pages"}
	links_1    = domain.Links{
		Internal: 1, External: 1, InAccessible: 0, AllLinks: AllLinks_1,
	}

	expected_1 = &domain.Result{
		Url:         "http://dummy.com/1",
		HtmlVersion: "HTML 5",
		PageTitle:   "The Go Playground",
		IsLoginPage: true,
		Headers:     headers_1,
		Links:       links_1,
	}
	URL1 = "http://dummy.com/1"
	URL2 = "http://dummy.com/2"

	headers_2 = domain.Headers{
		H1: 1, H2: 0, H3: 1, H4: 0, H5: 1, H6: 0,
	}
	links_2 = domain.Links{
		Internal: 3, External: 2, InAccessible: 0, AllLinks: AllLinks_2,
	}
	expected_2 = &domain.Result{
		Url:         "http://dummy.com/2",
		HtmlVersion: "HTML 4.01",
		PageTitle:   "The Test Case 2",
		IsLoginPage: false,
		Headers:     headers_2,
		Links:       links_2,
	}
)

// go test Parser_test.go Parser.go
// go test -race --cover -v Parser_test.go Parser.go

func TestParse_AccessibleHttpMock(t *testing.T) {
	t.Parallel()

	// create a new reader with the HTML
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		r := ioutil.NopCloser(bytes.NewReader([]byte(HTML_INPUT)))
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	client := webclient.Client
	service := NewParserService(client)

	res, err := Parse(*service, DOC_ONE, URL1)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expected_1, res)
}

func TestParse_InAccessibleHttpMock(t *testing.T) {
	t.Parallel()

	// create a new reader with the HTML
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		r := ioutil.NopCloser(bytes.NewReader([]byte(HTML_INPUT)))
		return &http.Response{
			StatusCode: 500,
			Body:       r,
		}, nil
	}

	client := webclient.Client
	service := NewParserService(client)

	res, err := Parse(*service, DOC_TWO, URL2)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expected_2, res)
}
