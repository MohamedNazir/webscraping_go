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
	ctrl := handler.HtmlParserHandler{Service: HtmlParserService}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", ctrl.IndexHandler)
	mux.HandleFunc("/search", ctrl.SearchHandler)
	http.ListenAndServe(":"+port, mux)
}
