package app

import (
	"crypto/tls"
	"log"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	HtmlParserService := service.NewHtmlParserService(client)
	hh := handler.HtmlParserHandler{Service: HtmlParserService}

	mux := http.NewServeMux()
	mux.HandleFunc("/", hh.IndexHandler)
	mux.HandleFunc("/search", hh.SearchHandler)

	log.Printf("Server starting at port :%s \n", port)
	http.ListenAndServe(":"+port, mux)
}
