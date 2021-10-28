package main

import (
	"context"

	pb "github.com/akileshsethu/grello/proto"
)

func (app *application) CreateBoardHandler(ctx context.Context, boardRequest *pb.BoardRequest) (*pb.Board, error) {
	boardResult, _ := <-app.models.Board.InsertAsync(ctx, boardRequest.Name, boardRequest.Description, boardRequest.MemberIds)
	//do you use the false from channel alone to determine if its closed and throw an error or explicitly get an error from request as in comment in InsertAsync
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
