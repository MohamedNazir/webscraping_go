package service

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MohamedNazir/webscraper/webscraper/domain"
	"golang.org/x/net/html"
)

var (
	HTMLString string = `
<!doctype html>
<html>
	<head>
		<title>The Go Playground</title>
		<link rel="stylesheet" href="/static/style.css">
	</head>
	<body itemscope itemtype="http://schema.org/CreativeWork">
		<input type="button" value="Run" id="embedRun">
		<div id="banner">
			<h1> Sample text </h1>
			<h2> Sample text </h2>
			<h3> Sample text </h3>
			<h4> Sample text </h4>
			<h5> Sample text </h5>
			<h6> Sample text </h6>
			<input type="password" name="password" id ="test" >
			<input type="button" value="About" id="aboutButton">
		</div>
		<div id="output"></div>
		<img itemprop="image" src="/static/gopher.png" style="display:none">
		<div id="about">
<p><b>About the Playground</b></p>
<a href="http://playground"> external link </a>
<a href="/pages"> Internal link </a>
</p>
		</div>
	</body>
</html>
`
)

func TestParse(t *testing.T) {

	z := html.NewTokenizer(strings.NewReader(HTMLString))

	headers := domain.Headers{
		H1: 1, H2: 1, H3: 1, H4: 1, H5: 1, H6: 1,
	}
	links := domain.Links{
		Internal: 1, External: 1, InAccessible: 0, AllLinks: 2, UniqueLinks: 0,
	}
	wanted := &domain.Result{
		HtmlVersion: "HTML 5",
		PageTitle:   "The Go Playground",
		IsLoginPage: true,
		Headers:     headers,
		Links:       links,
	}
	type args struct {
		doc *html.Tokenizer
	}
	tests := []struct {
		name    string
		args    *html.Tokenizer
		want    *domain.Result
		wantErr bool
	}{
		{"Parser Logic Test", z, wanted, false},
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
