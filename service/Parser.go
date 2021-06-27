package service

import (
	"log"
	"net/http"
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
	HTTP     = "http://"
	HTTPS    = "https://"
	SLASH    = "/"
)

type Response struct {
	*http.Response
	err error
}

//Parse to parse the html Tokens and extract the details
func Parse(service DefaultParserService, doc *html.Tokenizer, URL string) (*domain.Result, error) {

	h := &domain.Headers{}
	l := &domain.Links{}
	result := &domain.Result{Url: URL}
	allLinks := []string{}

	headerChan := make(chan string)
	fieldChan := make(chan map[string]interface{})
	linksChan := make(chan string)

	reqChan := make(chan *http.Request)
	respChan := make(chan Response)

	// go-routine to iterate over the html content
	go Iterate(doc, headerChan, linksChan, fieldChan)

	wg := sync.WaitGroup{}
	wg.Add(1)
	// to Handle header channel
	go HandleHeaders(headerChan, &wg, h)

	wg.Add(1)
	// to Handle links channel
	go HandleLinks(linksChan, &wg, l, URL)

	wg.Add(1)
	// to Handle fields channel
	go HandleFields(fieldChan, &wg, result)

	wg.Wait()

	allLinks = l.AllLinks
	go dispatcher(reqChan, allLinks)
	go workerPool(reqChan, respChan, len(allLinks), service)
	fail := consumer(respChan, len(allLinks))

	result.Links = *l
	result.Headers = *h
	result.Links.InAccessible = fail

	return result, nil
}

//Iterate function Iterates over the html Tokens and send the extrated values in the corresponding channels.
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

// getHref function to extrat the href link from anchor tag's attributes.
func getHref(t html.Token) string {
	for _, a := range t.Attr {
		if a.Key == HREF {
			return a.Val
		}
	}
	return ""
}

// isLoginPage function to find whether the page has any password type field.
func isLoginPage(t html.Token) (ok bool) {
	for _, a := range t.Attr {
		if a.Key == TYPE && a.Val == PASSWORD {
			return true
		}
	}
	return
}

// getHtmlVersion to get the html version of the web page.
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

// isInternalLink to find whether the given link is internal.
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

func dispatcher(reqChan chan *http.Request, links []string) {
	defer close(reqChan)
	for _, link := range links {
		req, err := http.NewRequest(http.MethodGet, link, nil)
		if err != nil {
			log.Println(err)
		}
		reqChan <- req
	}
}

// Worker Pool
func workerPool(reqChan chan *http.Request, respChan chan Response, size int, service DefaultParserService) {
	for i := 0; i < size; i++ {
		go service.IsAccessible(reqChan, respChan)
	}
}

//IsReachable func is a worker function which takes request and response and access the web page.
func (s DefaultParserService) IsAccessible(reqChan chan *http.Request, respChan chan Response) {
	for req := range reqChan {
		resp, err := s.webclient.Do(req)
		r := Response{resp, err}
		respChan <- r
	}
}

// Consumer
func consumer(respChan chan Response, size int) int64 {
	var failure int64 = 0
	for i := 0; i < size; i++ {
		r := <-respChan
		if r.err != nil {
			failure++
		}
	}
	return failure
}

// HandleHeaders function to Handle header channel
func HandleHeaders(headerChan <-chan string, wg *sync.WaitGroup, h *domain.Headers) {
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

}

// HandleLinks function to Handle links channel
func HandleLinks(linksChan <-chan string, wg *sync.WaitGroup, l *domain.Links, URL string) {
	defer wg.Done()
	//	allLinks := []string{}
	for link := range linksChan {
		ok := isInternalLink(link)
		if ok {
			l.AddInternal()
			link = strings.TrimSuffix(URL, SLASH) + SLASH + strings.TrimPrefix(link, SLASH)
		} else {
			l.AddExternal()
		}
		l.AddAllLinks(link)
	}
	//	return allLinks
}

func HandleFields(fieldChan <-chan map[string]interface{}, wg *sync.WaitGroup, r *domain.Result) {
	defer wg.Done()
	for field := range fieldChan {
		for k, v := range field {
			if k == LOGIN {
				r.IsLoginPage = v.(bool)
			}
			if k == VERSION {
				r.HtmlVersion = v.(string)
			}
			if k == TITLE {
				r.PageTitle = v.(string)
			}
		}
	}
}
