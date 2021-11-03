package main

import (
	"context"

	pb "github.com/akileshsethu/grello/proto"
)

func (app *application) CreateBoardHandler(ctx context.Context, boardRequest *pb.BoardRequest) (*pb.Board, error) {
	boardChan := app.models.Board.Insert(ctx, boardRequest.Name, boardRequest.Description, boardRequest.MemberIds)
	usersChan := app.models.User.GetAll(ctx, boardRequest.MemberIds)

	boardResult := <-boardChan
	if boardResult.Err != nil {
		return nil, boardResult.Err
	}

	usersResult := <-usersChan
	if usersResult.Err != nil {
		return nil, usersResult.Err
	}
	var users []*pb.User
	for _, u := range usersResult.Users {
		user := &pb.User{
			Id:   u.ID.Hex(),
			Name: u.Name,
		}
		users = append(users, user)
	}

	boardResponse := &pb.Board{
		Id:          boardResult.Board.ID.Hex(),
		Name:        boardResult.Board.Name,
		Description: boardResult.Board.Description,
		Users:       users,
	}
	return boardResponse, nil
}

func (app *application) GetBoardHandler(ctx context.Context, request *pb.GetBoardRequest) (*pb.Board, error) {
	boardChan := app.models.Board.Get(ctx, request.Id)
	boardResult := <-boardChan
	if boardResult.Err != nil {
		return nil, boardResult.Err
	}

	usersChan := app.models.User.GetAll(ctx, boardResult.Board.UserIds)
	tasksChan := app.models.Task.GetAll(ctx, boardResult.Board.TaskIds)

	usersResult := <-usersChan
	if usersResult.Err != nil {
		return nil, usersResult.Err
	}
	var users []*pb.User
	userMap := make(map[string]*pb.User)
	for _, u := range usersResult.Users {
		user := &pb.User{
			Id:   u.ID.Hex(),
			Name: u.Name,
		}
		users = append(users, user)
		userMap[u.ID.Hex()] = user
	}

	tasksResult := <-tasksChan
	if tasksResult.Err != nil {
		return nil, tasksResult.Err
	}
	var tasks []*pb.Task
	for _, t := range tasksResult.Tasks {
		user, prs := userMap[t.AssignedTo]
		if !prs {
			user = &pb.User{Name: "Unknown User"}
		}
		task := &pb.Task{
			Id:          t.ID.Hex(),
			Name:        t.Name,
			Description: t.Description,
			AssignedTo:  user,
		}
		tasks = append(tasks, task)
	}

	boardResponse := &pb.Board{
		Id:          boardResult.Board.ID.Hex(),
		Name:        boardResult.Board.Name,
		Description: boardResult.Board.Description,
		Users:       users,
		Tasks:       tasks,
	}
	return boardResponse, nil
}
