package service

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/MohamedNazir/webscraper/domain"
	wc "github.com/MohamedNazir/webscraper/webclient"
	"golang.org/x/net/html"
)

const (
	CONN_FAILED = "failed to connect to the URL -"
	INVALID_URL = "invalid url http/https schema or host missing"
)

type ParserService interface {
	ParseHtml(url string) (*domain.Result, error)
	IsReachable(url string) bool
}

type DefaultParserService struct {
	webclient wc.HTTPClient
}

func NewParserService(webclient wc.HTTPClient) *DefaultParserService {
	return &DefaultParserService{webclient: webclient}
}

func (s DefaultParserService) ParseHtml(u string) (*domain.Result, error) {

	res, _ := url.Parse(u)

	if isEmpty(res.Scheme) || isEmpty(res.Host) {
		return nil, fmt.Errorf(INVALID_URL)
	}

	request, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.webclient.Do(request)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf(CONN_FAILED+" %s", u)
	}

	doc := html.NewTokenizer(resp.Body)
	return Parse(doc, u)

}

func isEmpty(s string) bool {
	return s == ""
}
