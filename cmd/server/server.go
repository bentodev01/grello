package main

import (
	"context"

	pb "github.com/akileshsethu/grello/proto"
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

func (app *application) NewServer() *server {
	return &server{}
}
