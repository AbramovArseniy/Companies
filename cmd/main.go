package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/AbramovArseniy/Companies/internal/cfg"
	httphandler "github.com/AbramovArseniy/Companies/internal/handlers/http"
)

func main() {
	cfg := cfg.New()
	handler := httphandler.New(cfg)
	router := handler.Route()
	srv := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	idleConnsClosed := make(chan struct{})
	gr := sync.WaitGroup{}
	sigs := make(chan os.Signal, 1)
	gr.Add(1)
	go func() {
		err := srv.ListenAndServe()
		log.Println("error while starting server:", err)
		gr.Done()
	}()
	log.Println("http server started on:", cfg.Address)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigs
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
		gr.Done()
	}()
	gr.Wait()
}
