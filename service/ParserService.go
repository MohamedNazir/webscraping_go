package service

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/MohamedNazir/webscraper/domain"
	"github.com/MohamedNazir/webscraper/webclient"
	"golang.org/x/net/html"
)

const (
	CONN_FAILED = "failed to connect to the URL -"
	INVALID_URL = "invalid url http/https schema or host missing"
)

type ParserService interface {
	ParseHtml(URL string) (*domain.Result, error)
	IsAccessible(reqChan chan *http.Request, respChan chan Response)
}

type DefaultParserService struct {
	webclient webclient.HTTPClient
}

func NewParserService(webclient webclient.HTTPClient) *DefaultParserService {
	return &DefaultParserService{webclient: webclient}
}

func (s DefaultParserService) ParseHtml(URL string) (*domain.Result, error) {

	res, _ := url.Parse(URL)

	if isEmpty(res.Scheme) || isEmpty(res.Host) {
		return nil, fmt.Errorf(INVALID_URL)
	}

	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.webclient.Do(request)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf(CONN_FAILED+" %s", URL)
	}

	doc := html.NewTokenizer(resp.Body)
	return Parse(s, doc, URL)
}

func isEmpty(s string) bool {
	return s == ""
}
