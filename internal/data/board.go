package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BoardResult struct {
	Board Board
	Err   error
}

type Board struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	UserIds     []string `json:"user_ids"`
	TaskIds     []string `json:"task_ids"`
}

type BoardModel struct {
	DB *mongo.Database
}

func (m BoardModel) Insert(name string, description string, userIds []string) (Board, error) {
	res, err := m.DB.Collection("board").InsertOne(context.Background(), bson.D{{"name", name}, {"description", description}, {"user_ids", userIds}})
	if err != nil {
		return Board{}, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	board := Board{
		ID:          id.String(),
		Name:        name,
		Description: description,
		UserIds:     userIds,
	}
	return board, nil
}

func (m BoardModel) InsertAsync(ctx context.Context, name string, description string, userIds []string) <-chan BoardResult {
	resultChan := make(chan BoardResult)

	go func() {
		defer close(resultChan)
		select {
		case <-ctx.Done():
			// error needs to be returned in resultChan. should it be sent here or should we panic here and handle it with recover in defer statement before closing resultChan?
			return
		case resultChan <- m.insert(ctx, name, description, userIds):
		}
	}()

	return resultChan
}

func (m BoardModel) insert(ctx context.Context, name string, description string, userIds []string) BoardResult {
	res, err := m.DB.Collection("board").InsertOne(ctx, bson.D{{"name", name}, {"description", description}, {"user_ids", userIds}})
	if err != nil {
		return BoardResult{Err: err}
	}
	id := res.InsertedID.(primitive.ObjectID)
	board := Board{
		ID:          id.String(),
		Name:        name,
		Description: description,
		UserIds:     userIds,
	}
	return BoardResult{Board: board}
}
