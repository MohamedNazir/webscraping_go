package service

import (
	"github.com/MohamedNazir/webscraper/domain"
	"github.com/MohamedNazir/webscraper/network"
)

type ParserService interface {
	ParseHtml(url string) (*domain.Result, error)
}

type DefaultParserService struct {
	httpComm network.HttpCommunicator
}

func NewParserService(httpComm network.HttpCommunicator) *DefaultParserService {
	return &DefaultParserService{httpComm: httpComm}
}

func (s DefaultParserService) ParseHtml(u string) (*domain.Result, error) {

	doc, err := s.httpComm.GetWebPageContent(u)
	if err != nil {
		return nil, err
	}
	return Parse(doc, u)

}
