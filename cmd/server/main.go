package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/akileshsethu/grello/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type config struct {
	host string
	port int
	env  string
}

type application struct {
	config config
}

func main() {
	log.Println("Starting server..")
	var cfg config
	flag.StringVar(&cfg.host, "host", "localhost", "Service host")
	flag.IntVar(&cfg.port, "port", 50051, "Service port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

	app := &application{
		config: cfg,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.host, cfg.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grpcServer := app.NewServer()
	pb.RegisterBoardServiceServer(s, grpcServer)
	reflection.Register(s)
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
