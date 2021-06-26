package service

import (
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/MohamedNazir/webscraper/domain"
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
	DTD      = "DTD"
	HREF     = "href"
	TYPE     = "type"
	INPUT    = "input"
	PASSWORD = "password"
	HTML5    = "HTML 5"
)

func Parse(doc *html.Tokenizer, url string) (*domain.Result, error) {

	h := &domain.Headers{}
	l := &domain.Links{}
	res := &domain.Result{Url: url}

	headerChan := make(chan string)
	fieldChan := make(chan map[string]interface{}, 3)
	linksChan := make(chan string)

	// go-routine to iterate over the html content
	go Iterate(doc, headerChan, linksChan, fieldChan)

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		for hTag := range headerChan {
			if hTag == H1 {
				h.AddH1()
			}
			if hTag == H2 {
				h.AddH2()
			}
			if hTag == H3 {
				h.AddH3()
			}
			if hTag == H4 {
				h.AddH4()
			}
			if hTag == H5 {
				h.AddH5()
			}
			if hTag == H6 {
				h.AddH6()
			}
		}

	}()

	go func() {
		defer wg.Done()
		for link := range linksChan {
			ok := isInternalLink(link)
			if ok {
				l.AddInternal()
			} else {
				l.AddExternal()
			}
			yes := isAccessible(link)
			if !yes {
				l.AddInAccessible()
			}
		}
	}()

	go func() {
		defer wg.Done()
		for field := range fieldChan {
			for k, v := range field {
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
		}
	}()

	wg.Wait()

	res.Links = *l
	res.Headers = *h

	return res, nil
}

func Iterate(doc *html.Tokenizer, headerChan chan string, linksChan chan string, fieldChan chan map[string]interface{}) {

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

	close(fieldChan)
	close(headerChan)
	close(linksChan)
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

func isAccessible(link string) bool {
	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") || (strings.HasPrefix(link, "/") && len(link) > 1) {
		return true
	}
	return false
}
