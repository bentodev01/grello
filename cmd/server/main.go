package main

import (
	"context"
	"log"
	"net"

	pb "github.com/akileshsethu/grello/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedBoardServiceServer
}

func (s *server) GetBoard(ctx context.Context, req *pb.GetBoardRequest) (*pb.Board, error) {
	return &pb.Board{Name: "test"}, nil
}

func (s *server) CreateBoard(ctx context.Context, req *pb.BoardRequest) (*pb.Board, error) {
	return &pb.Board{Name: "test"}, nil
}

func (s *server) UpdateBoard(ctx context.Context, req *pb.UpdateBoardRequest) (*pb.Board, error) {
	return &pb.Board{Name: "test"}, nil
}

func (s *server) DeleteBoard(ctx context.Context, req *pb.DeleteBoardRequest) (*pb.DeleteBoardResponse, error) {
	return &pb.DeleteBoardResponse{Message: "Deleted"}, nil
}

func (s *server) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.Board, error) {
	return &pb.Board{Name: "test"}, nil
}

func main() {
	log.Println("Starting server..")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBoardServiceServer(s, &server{})
	reflection.Register(s)
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
