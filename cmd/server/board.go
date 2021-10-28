package main

import (
	"context"

	pb "github.com/akileshsethu/grello/proto"
)

func (app *application) CreateBoardHandler(ctx context.Context, boardRequest *pb.BoardRequest) (*pb.Board, error) {
	boardResult := <-app.models.Board.InsertAsync(ctx, boardRequest.Name, boardRequest.Description, boardRequest.MemberIds)
	if boardResult.Err != nil {
		return nil, boardResult.Err
	}
	boardResponse := &pb.Board{
		Id:          boardResult.Board.ID,
		Name:        boardResult.Board.Name,
		Description: boardResult.Board.Description,
		Users:       make([]*pb.User, 0),
		Tasks:       make([]*pb.Task, 0),
	}
	return boardResponse, nil
}
