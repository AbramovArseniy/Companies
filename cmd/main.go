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
	"time"

	"google.golang.org/grpc"

	"github.com/AbramovArseniy/Companies/internal/cfg"
	grpchandler "github.com/AbramovArseniy/Companies/internal/handlers/grpc"
	pb "github.com/AbramovArseniy/Companies/internal/handlers/grpc/proto"
	httphandler "github.com/AbramovArseniy/Companies/internal/handlers/http"
	"github.com/AbramovArseniy/Companies/internal/kafka"
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
	kafka1, err := kafka.New(dbPool, *cfg)
	if err != nil {
		log.Println("cannot create kafka consumer:", err)
	}
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
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
	gr := sync.WaitGroup{}
	pb.RegisterCompaniesServiceServer(grpcSrv, grpcHandler)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	gr.Add(1)
	go func() {
		<-sigint
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		err := kafka1.Close()
		if err != nil {
			log.Println("cannot close kafka consumer:", err)
		}
		grpcSrv.GracefulStop()
		close(idleConnsClosed)
	}()
	go func() {
		if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
		<-idleConnsClosed
		gr.Done()
	}()
	log.Println("http server started at:", cfg.Address)
	gr.Add(1)
	go func() {
		if err := grpcSrv.Serve(listen); err != nil {
			log.Println("error on grpc server:", err)
		}
		gr.Done()
	}()
	err = kafka1.ListenTagChanges(cfg.ChangesTopic)
	if err != nil {
		log.Println("cannot listen to tag changes:", err)
	}
	log.Println("gRPC server started at:", listen.Addr())
	gr.Wait()
}
