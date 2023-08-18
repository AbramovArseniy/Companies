// package main starts http and grpc services
package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"

	"github.com/AbramovArseniy/Companies/internal/cfg"
	grpchandler "github.com/AbramovArseniy/Companies/internal/handlers/grpc"
	pb "github.com/AbramovArseniy/Companies/internal/handlers/grpc/proto"
	httphandler "github.com/AbramovArseniy/Companies/internal/handlers/http"
	"github.com/AbramovArseniy/Companies/internal/storage/postgres"
)

func main() {
	cfg := cfg.New()
	dbPool, err := postgres.New(cfg.DBAddress)
	if err != nil {
		log.Println("error while connecting to database:", err)
		return
	}
	handler, err := httphandler.New(dbPool)
	if err != nil {
		log.Println("error while creating http handler:", err)
		return
	}
	router := handler.Route()
	httpSrv := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	idleConnsClosed := make(chan struct{})
	gr := sync.WaitGroup{}
	sigs := make(chan os.Signal, 1)
	gr.Add(1)
	go func() {
		err := httpSrv.ListenAndServe()
		log.Println("error on http server:", err)
		gr.Done()
	}()
	log.Println("http server started at:", cfg.Address)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigs
		if err := httpSrv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
		gr.Done()
	}()
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	grpcSrv := grpc.NewServer()
	grpcHandler, err := grpchandler.New(dbPool)
	if err != nil {
		log.Println("error while creating grpc handler:", err)
		return
	}
	pb.RegisterCompaniesServiceServer(grpcSrv, grpcHandler)
	gr.Add(1)
	go func() {
		if err := grpcSrv.Serve(listen); err != nil {
			log.Println("error on grpc server:", err)
		}
		gr.Done()
	}()
	log.Println("gRPC server started at:", listen.Addr())
	gr.Wait()
}
