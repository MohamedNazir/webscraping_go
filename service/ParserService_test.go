package service

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/MohamedNazir/webscraper/utils/mocks"
	wc "github.com/MohamedNazir/webscraper/webclient"

	"github.com/stretchr/testify/assert"
)

const (
	// build response HTML
	HTML_INPUT  = `<!doctype html> <html> <head> <title>Mock Test case</title></head></body> </html>`
	URL         = "http://google.com"
	HtmlVersion = "HTML 5"
	Title       = "Mock Test case"
	EXP_ERR     = "failed to connect to the URL - http://google.com"
)

func init() {
	wc.Client = &mocks.MockClient{}
}

func TestParseHtml(t *testing.T) {

	// create a new reader with the HTML
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		r := ioutil.NopCloser(bytes.NewReader([]byte(HTML_INPUT)))
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	service := NewParserService(wc.Client)
	res, err := service.ParseHtml(URL)

	assert.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, res.HtmlVersion, HtmlVersion)
	assert.Equal(t, res.PageTitle, Title)
}

func TestParseHtml_InvalidUrl(t *testing.T) {
	URL_INVALID := "abc"
	service := NewParserService(wc.Client)
	res, err := service.ParseHtml(URL_INVALID)

	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.Equal(t, INVALID_URL, err.Error())
}

func TestParseHtml_PageNotfound(t *testing.T) {

	// create a 404 response
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		r := ioutil.NopCloser(bytes.NewReader([]byte(HTML_INPUT)))
		return &http.Response{
			StatusCode: 404,
			Body:       r,
		}, nil
	}

	service := NewParserService(wc.Client)
	res, err := service.ParseHtml(URL)

	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.Equal(t, EXP_ERR, err.Error())
}
