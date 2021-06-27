package network

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

const (
	CONN_FAILED = "failed to connect to the URL - %s"
	INVALID_URL = "invalid url http/https schema or host missing"
)

type HttpCommunicator interface {
	GetWebPageContent(u string) (*html.Tokenizer, error)
}

type WebPageDownloader struct {
	client *http.Client
}

func NewWebPageDownloader(client *http.Client) *WebPageDownloader {
	return &WebPageDownloader{client: client}
}

func (d WebPageDownloader) GetWebPageContent(u string) (*html.Tokenizer, error) {
	res, _ := url.Parse(u)

	if isEmpty(res.Scheme) || isEmpty(res.Host) {
		return nil, fmt.Errorf(INVALID_URL)
	}

	resp, err := d.client.Get(u)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf(CONN_FAILED, u)
	}

	//	defer resp.Body.Close()
	doc := html.NewTokenizer(resp.Body)

	return doc, nil
}

func isEmpty(s string) bool {
	return s == ""
}
