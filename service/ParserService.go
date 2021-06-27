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

type ParserService interface {
	ParseHtml(url string) (*domain.Result, error)
}

type DefaultParserService struct {
	client *http.Client
}

func NewParserService(client *http.Client) *DefaultParserService {
	return &DefaultParserService{client: client}
}

func (s DefaultParserService) ParseHtml(u string) (*domain.Result, error) {

	res, _ := url.Parse(u)

	if isEmpty(res.Scheme) || isEmpty(res.Host) {
		return nil, fmt.Errorf(INVALID_URL)
	}

	resp, err := s.client.Get(u)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf(CONN_FAILED, u)
	}

	doc := html.NewTokenizer(resp.Body)
	return Parse(doc, u)

}

func isEmpty(s string) bool {
	return s == ""
}
