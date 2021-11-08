package main

import (
	"context"

	pb "github.com/bentodev01/grello/proto"
)

func (app *application) CreateTaskHandler(ctx context.Context, taskRequest *pb.AddTaskRequest) (*pb.Task, error) {
	taskResult := app.models.Task.Insert(ctx, taskRequest.Name, taskRequest.Description, taskRequest.BoardId, taskRequest.AssignedTo)
	if taskResult.Err != nil {
		return nil, taskResult.Err
	}

	err := app.models.Board.AddTask(ctx, taskRequest.BoardId, taskResult.Task.ID.Hex())
	if err != nil {
		return nil, err
	}

	user, err := app.models.User.Get(ctx, taskResult.Task.AssignedTo)
	if err != nil {
		return nil, err
	}

	task := &pb.Task{
		Id:          taskResult.Task.ID.Hex(),
		Name:        taskResult.Task.Name,
		Description: taskResult.Task.Description,
		AssignedTo: &pb.User{
			Id:   user.ID.Hex(),
			Name: user.Name,
		},
	}
	return task, nil
}
