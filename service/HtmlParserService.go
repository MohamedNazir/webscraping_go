package service

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/MohamedNazir/webscraper/domain"
	"golang.org/x/net/html"
)

const (
	CONN_FAILED = "failed to connect to the URL - %s"
	INVALID_URL = "invalid url http/https schema or host missing"
)

type HtmlParserService interface {
	ParseHtml(url string) (*domain.Result, error)
}

type DefaultParaserService struct {
	client *http.Client
}

func NewHtmlParserService(client *http.Client) *DefaultParaserService {
	return &DefaultParaserService{client: client}
}

func (s DefaultParaserService) ParseHtml(u string) (*domain.Result, error) {
	res, _ := url.Parse(u)

	if res.Scheme == "" || res.Host == "" {
		return nil, fmt.Errorf(INVALID_URL)
	}

	response, err := s.client.Get(u)
	if err != nil || response.StatusCode != 200 {
		return nil, fmt.Errorf(CONN_FAILED, u)
	}
	defer response.Body.Close()
	doc := html.NewTokenizer(response.Body)

	return Parse(doc, u)

}
