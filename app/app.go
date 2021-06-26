package app

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		osCall := <-c
		log.Printf("system call:%+v", osCall)
		cancel()
	}()

	if err := serve(ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
}

func serve(ctx context.Context) (err error) {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	hh := handler.HtmlParserHandler{Service: service.NewHtmlParserService(client)}

	mux := http.NewServeMux()
	mux.HandleFunc("/", hh.IndexHandler)
	mux.HandleFunc("/search", hh.SearchHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	log.Printf("server started")

	<-ctx.Done()

	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
