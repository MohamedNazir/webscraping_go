package service

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MohamedNazir/webscraper/domain"
	"golang.org/x/net/html"
)

var (
	HTML_FIRST_INPUT                  = `<!doctype html> <html> <head> <title>The Go Playground</title> <link rel="stylesheet" href="/static/style.css"> </head> <body itemscope itemtype="http://schema.org/CreativeWork"> <input type="button" value="Run" id="embedRun"> <div id="banner"> <h1> Sample text </h1> <h2> Sample text </h2> <h3> Sample text </h3> <h4> Sample text </h4> <h5> Sample text </h5> <h6> Sample text </h6> <input type="password" name="password" id ="test" > <input type="button" value="About" id="aboutButton"> </div> <div id="output"></div> <img itemprop="image" src="/static/gopher.png" style="display:none"> <div id="about"> <p><b>About the Playground</b></p> <a href="http://playground"> external link </a> <a href="/pages"> Internal link </a> </p> </div> </body> </html>`
	HTML_SECOND_INPUT                 = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd"> <html> <head> <title>The Test Case 2</title> </head> <body> <div id="banner"> <h1> Sample text </h1> <h3> Sample text </h3> <h5> Sample text </h5> <input type="button" value="About" id="aboutButton"> </div> <img itemprop="image" src="/static/gopher.png" style="display:none"> <div id="about"> <p><b>About the Playground</b></p> <a href="http://playground"> external link </a> <a href="/pages"> Internal link </a> <a href="/"> Broken link 1 </a> <a href="#"> Broken link 2 </a> </p> </div> </body> </html>`
	DOC_ONE           *html.Tokenizer = html.NewTokenizer(strings.NewReader(HTML_FIRST_INPUT))
	DOC_TWO           *html.Tokenizer = html.NewTokenizer(strings.NewReader(HTML_SECOND_INPUT))
	headers_1                         = domain.Headers{
		H1: 1, H2: 1, H3: 1, H4: 1, H5: 1, H6: 1,
	}
	links_1 = domain.Links{
		Internal: 1, External: 1, InAccessible: 0, AllLinks: 2, UniqueLinks: 0,
	}
	headers_2 = domain.Headers{
		H1: 1, H2: 0, H3: 1, H4: 0, H5: 1, H6: 0,
	}
	links_2 = domain.Links{
		Internal: 3, External: 1, InAccessible: 0, AllLinks: 4, UniqueLinks: 0,
	}
	expected_1 = &domain.Result{
		HtmlVersion: "HTML 5",
		PageTitle:   "The Go Playground",
		IsLoginPage: true,
		Headers:     headers_1,
		Links:       links_1,
	}
	expected_2 = &domain.Result{
		HtmlVersion: "HTML 4.01",
		PageTitle:   "The Test Case 2",
		IsLoginPage: false,
		Headers:     headers_2,
		Links:       links_2,
	}
)

// go test Parser_test.go Parser.go
// go test -race --cover -v Parser_test.go Parser.go

func TestParse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    *html.Tokenizer
		want    *domain.Result
		wantErr bool
	}{
		{"test case 1", DOC_ONE, expected_1, false},
		{"test cas 2", DOC_TWO, expected_2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

// go test -cpu 1,2,4,8 -benchmem -run=^$ -bench . Parser_test.go
func BenchmarkParse(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Parse(DOC_ONE)
	}

}
