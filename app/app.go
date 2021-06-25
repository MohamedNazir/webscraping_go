package app

import (
	"crypto/tls"
	"net/http"
	"os"

	"github.com/MohamedNazir/webscraper/handler"
	"github.com/MohamedNazir/webscraper/service"
)

var (
	config    = &tls.Config{InsecureSkipVerify: true}
	transport = &http.Transport{
		TLSClientConfig: config,
	}
	client *http.Client
)

func init() {
	client = &http.Client{
		Transport: transport,
	}
}

func StartApplication() {

	HtmlParserService := service.NewHtmlParserService(client)
	hh := handler.HtmlParserHandler{Service: HtmlParserService}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", hh.IndexHandler)
	mux.HandleFunc("/search", hh.SearchHandler)
	http.ListenAndServe(":"+port, mux)
}
