package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BoardResult struct {
	Board Board
	Err   error
}

type Board struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	UserIds     []string           `bson:"user_ids"`
	TaskIds     []string           `bson:"task_ids"`
}

func (b Board) ToJson() (string, error) {
	id := b.ID.Hex()
	js, err := json.Marshal(struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		UserIds     []string `json:"user_ids"`
		TaskIds     []string `json:"task_ids"`
	}{
		ID:          id,
		Name:        b.Name,
		Description: b.Description,
		UserIds:     b.UserIds,
		TaskIds:     b.TaskIds,
	})
	return string(js[:]), err
}

func (b *Board) FromJson(jsonStr string) error {
	var input struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		UserIds     []string `json:"user_ids"`
		TaskIds     []string `json:"task_ids"`
	}

	jsonBytes := []byte(jsonStr)
	err := json.Unmarshal(jsonBytes, &input)
	if err != nil {
		return err
	}

	objectId, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return err
	}

	b.ID = objectId
	b.Name = input.Name
	b.Description = input.Description
	b.UserIds = input.UserIds
	b.TaskIds = input.TaskIds

	return nil
}

type BoardModel struct {
	DB *mongo.Database
}

func (m BoardModel) Insert(ctx context.Context, name string, description string, userIds []string) <-chan BoardResult {
	resultChan := make(chan BoardResult)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				resultChan <- BoardResult{Err: err.(error)}
			}
			close(resultChan)
		}()

		select {
		case <-ctx.Done():
			resultChan <- BoardResult{Err: ctx.Err()}
			return
		case resultChan <- m.insert(ctx, name, description, userIds):
		}
	}()

	return resultChan
}

func (m BoardModel) insert(ctx context.Context, name string, description string, userIds []string) BoardResult {
	res, err := m.DB.Collection("boards").InsertOne(ctx, bson.D{{"name", name}, {"description", description}, {"user_ids", userIds}})
	if err != nil {
		return BoardResult{Err: err}
	}
	id := res.InsertedID.(primitive.ObjectID)
	board := Board{
		ID:          id,
		Name:        name,
		Description: description,
		UserIds:     userIds,
	}
	return BoardResult{Board: board}
}

func (m BoardModel) Get(ctx context.Context, id string) BoardResult {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return BoardResult{Err: err}
	}
	board := Board{}
	err = m.DB.Collection("boards").FindOne(ctx, bson.D{{"_id", objectId}}).Decode(&board)
	if err != nil {
		return BoardResult{Err: err}
	}
	return BoardResult{Board: board}
}

func (m BoardModel) AddTask(ctx context.Context, id string, taskId string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$push", bson.D{{"task_ids", taskId}}}}
	res, err := m.DB.Collection("boards").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount != 1 {
		return errors.New(fmt.Sprintf("Wrong number of updated columns %d", res.MatchedCount))
	}

	return nil
}
