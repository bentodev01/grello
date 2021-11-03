package main

import (
	"context"

	pb "github.com/akileshsethu/grello/proto"
)

type server struct {
	pb.UnimplementedBoardServiceServer
	app *application
}

func (s *server) GetBoard(ctx context.Context, req *pb.GetBoardRequest) (*pb.Board, error) {
	board, err := s.app.GetBoardHandler(ctx, req)
	return board, err
}

func (s *server) CreateBoard(ctx context.Context, req *pb.BoardRequest) (*pb.Board, error) {
	board, err := s.app.CreateBoardHandler(ctx, req)
	return board, err
}

func (s *server) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.TaskResponse, error) {
	task, err := s.app.CreateTaskHandler(ctx, req)
	return &pb.TaskResponse{Task: task, BoardId: req.BoardId}, err
}

func NewServer(app *application) *server {
	return &server{app: app}
}
