package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	h "github.com/MohamedNazir/webscraper/handler"
	s "github.com/MohamedNazir/webscraper/service"
	wc "github.com/MohamedNazir/webscraper/webclient"
)

var (
	timeout time.Duration = 5 * time.Second
)

const (
	SERVER_STARTED   = "Server started"
	SERVER_STOPPED   = "Server stopped"
	SHUTDOWN_SUCCESS = "Server exited properly"
	SHUTDOWN_FAILED  = "server Shutdown Failed:%s"
)

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

	client := wc.Client
	hh := h.ParserHandler{Service: s.NewParserService(client)}

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

	log.Println(SERVER_STARTED)

	<-ctx.Done()

	log.Println(SERVER_STOPPED)

	ctxShutDown, cancel := context.WithTimeout(context.Background(), timeout)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf(SHUTDOWN_FAILED, err)
	}

	log.Println(SHUTDOWN_SUCCESS)

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
