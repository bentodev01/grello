package data

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TasksResult struct {
	Tasks []Task
	Err   error
}

type TaskResult struct {
	Task Task
	Err  error
}

type Task struct {
	ID          primitive.ObjectID `bson:"_id"`
	BoardId     string             `bson:"board_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	AssignedTo  string             `bson:"assigned_to"`
}

type TaskModel struct {
	DB *mongo.Database
}

// func (m TaskModel) Insert(ctx context.Context, name, description, boardId, assignedTo string) <-chan TaskResult {
// 	resultChan := make(chan TaskResult)

// 	go func() {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				resultChan <- TaskResult{Err: err.(error)}
// 			}
// 			close(resultChan)
// 		}()

// 		select {
// 		case <-ctx.Done():
// 			resultChan <- TaskResult{Err: ctx.Err()}
// 			return
// 		case resultChan <- m.insert(ctx, name, description, boardId, assignedTo):
// 		}
// 	}()

// 	return resultChan
// }

func (m TaskModel) Insert(ctx context.Context, name, description, boardId, assignedTo string) TaskResult {
	res, err := m.DB.Collection("tasks").InsertOne(ctx, bson.D{{"name", name}, {"description", description}, {"board_id", boardId}, {"assigned_to", assignedTo}})
	if err != nil {
		return TaskResult{Err: err}
	}
	id := res.InsertedID.(primitive.ObjectID)
	task := Task{
		ID:          id,
		BoardId:     boardId,
		Name:        name,
		Description: description,
		AssignedTo:  assignedTo,
	}
	return TaskResult{Task: task}
}

func (m TaskModel) GetAll(ctx context.Context, ids []string) <-chan TasksResult {
	resultChan := make(chan TasksResult)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				resultChan <- TasksResult{Err: err.(error)}
			}
			close(resultChan)
		}()

		select {
		case <-ctx.Done():
			resultChan <- TasksResult{Err: ctx.Err()}
			return
		case resultChan <- m.getAll(ctx, ids):
		}
	}()

	return resultChan
}

func (m TaskModel) getAll(ctx context.Context, ids []string) TasksResult {
	objectIds := make([]primitive.ObjectID, len(ids))
	var errorIds []string

	for _, id := range ids {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			errorIds = append(errorIds, id)
		}
		objectIds = append(objectIds, objectId)
	}

	if len(errorIds) > 0 {
		errorMsg := fmt.Sprintf("invalid object ids: %s", strings.Join(errorIds, ","))
		return TasksResult{Err: errors.New(errorMsg)}
	}

	filter := bson.D{{"_id", bson.D{{"$in", objectIds}}}}
	cur, err := m.DB.Collection("tasks").Find(ctx, filter)
	defer cur.Close(ctx)
	if err != nil {
		return TasksResult{Err: err}
	}

	var tasks []Task
	for cur.Next(ctx) {
		task := Task{}
		err = cur.Decode(&task)
		if err != nil {
			break
		}
		tasks = append(tasks, task)
	}

	if err != nil {
		return TasksResult{Err: err}
	}

	return TasksResult{Tasks: tasks, Err: nil}
}
