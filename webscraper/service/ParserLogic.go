package service

import (
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/MohamedNazir/webscraper/webscraper/domain"
	"golang.org/x/net/html"
)

const (
	TITLE    = "title"
	A        = "a"
	H1       = "h1"
	H2       = "h2"
	H3       = "h3"
	H4       = "h4"
	H5       = "h5"
	H6       = "h6"
	LOGIN    = "login"
	VERSION  = "version"
	HTML     = "html"
	DTD      = "dtd"
	HREF     = "href"
	TYPE     = "type"
	INPUT    = "input"
	PASSWORD = "password"
	HTML5    = "HTML 5"
)

func Parse(doc *html.Tokenizer) (*domain.Result, error) {

	res := &domain.Result{}
	var wg sync.WaitGroup
	wg.Add(4)
	hMap := make(map[string]int)
	intLink := []string{}
	extLink := []string{}
	fieldMap := make(map[string]interface{})

	headerChan := make(chan string)
	fieldChan := make(chan map[string]interface{})
	linksChan := make(chan string)

	isDone := make(chan bool, 1)
	go Iterate(doc, isDone, headerChan, linksChan, fieldChan)

	go func() {
		defer wg.Done()
		for hTag := range headerChan {
			hMap[hTag] = hMap[hTag] + 1
		}

	}()
	go func() {
		defer wg.Done()
		for link := range linksChan {
			ok := isInternalLink(link)
			if ok {
				intLink = append(intLink, link)
			} else {
				extLink = append(extLink, link)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for field := range fieldChan {
			for k, v := range field {
				fieldMap[k] = v
			}
		}
	}()

	<-isDone

	res.Headers.H1 = hMap[H1]
	res.Headers.H2 = hMap[H2]
	res.Headers.H3 = hMap[H3]
	res.Headers.H4 = hMap[H4]
	res.Headers.H5 = hMap[H5]
	res.Headers.H6 = hMap[H6]

	res.Links.AllLinks = len(intLink) + len(extLink)
	res.Links.Internal = len(intLink)
	res.Links.External = len(extLink)

	for k, v := range fieldMap {
		if k == LOGIN {
			res.IsLoginPage = v.(bool)
		}
		if k == VERSION {
			res.HtmlVersion = v.(string)
		}
		if k == TITLE {
			res.PageTitle = v.(string)
		}
	}

	return res, nil
}

func Iterate(doc *html.Tokenizer, isDone chan bool, headerChan chan string, linksChan chan string, fieldChan chan map[string]interface{}) {

	for tokenType := doc.Next(); tokenType != html.ErrorToken; {
		token := doc.Token()
		if tokenType == html.StartTagToken {
			if token.Data == A {
				url := getHref(token)
				linksChan <- url
				tokenType = doc.Next()
				continue
			}
			if token.Data == TITLE {
				tokenType = doc.Next()
				if tokenType == html.TextToken {
					m := make(map[string]interface{})
					m[TITLE] = string(doc.Token().Data)
					fieldChan <- m
				}
				tokenType = doc.Next()
				continue
			}
			if token.Data == INPUT {
				ok := isLoginPage(token)
				if ok {
					m := make(map[string]interface{})
					m[LOGIN] = true
					fieldChan <- m
				}
				tokenType = doc.Next()
				continue
			}
			if token.Data == H1 || token.Data == H2 || token.Data == H3 || token.Data == H4 || token.Data == H5 || token.Data == H6 {
				headerChan <- token.Data
				tokenType = doc.Next()
				continue
			}
		}

		if tokenType == html.DoctypeToken {
			version := getHtmlVersion(token)
			m := make(map[string]interface{})
			m[VERSION] = version
			fieldChan <- m
			tokenType = doc.Next()
			continue
		}
		tokenType = doc.Next()
	}

	isDone <- true

}

func getHref(t html.Token) string {
	for _, a := range t.Attr {
		if a.Key == HREF {
			return a.Val
		}
	}
	return ""
}

func isLoginPage(t html.Token) (ok bool) {
	for _, a := range t.Attr {
		if a.Key == TYPE && a.Val == PASSWORD {
			return true
		}
	}
	return
}

func getHtmlVersion(t html.Token) string {

	if strings.EqualFold(t.Data, HTML) {
		// because html 5 will have only HTML type as data. <!DOCTYPE html>
		return HTML5
	}

	// doc type for lower version of html
	//<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">

	parts := strings.Split(t.Data, "//")
	for _, v := range parts {
		if strings.HasPrefix(v, DTD) {
			version := strings.Split(v, " ")
			return version[1] + " " + version[2]
		}
	}
	return ""
}

func isInternalLink(link string) bool {

	u, err := url.Parse(link)
	if err != nil {
		log.Fatal(err)
	}
	if u.Scheme == "" {
		return true
	}
	return false
}
