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

type UsersResult struct {
	Users []User
	Err   error
}

type User struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type UserModel struct {
	DB *mongo.Database
}

func (m UserModel) Get(ctx context.Context, id string) (User, error) {
	user := User{}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	filter := bson.D{{"_id", objectId}}
	err = m.DB.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return user, nil
	}
	return user, nil
}

func (m UserModel) GetAll(ctx context.Context, ids []string) <-chan UsersResult {
	resultChan := make(chan UsersResult)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				resultChan <- UsersResult{Err: err.(error)}
			}
			close(resultChan)
		}()

		select {
		case <-ctx.Done():
			resultChan <- UsersResult{Err: ctx.Err()}
			return
		case resultChan <- m.getAll(ctx, ids):
		}
	}()

	return resultChan
}

func (m UserModel) getAll(ctx context.Context, ids []string) UsersResult {
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
		return UsersResult{Err: errors.New(errorMsg)}
	}

	filter := bson.D{{"_id", bson.D{{"$in", objectIds}}}}
	cur, err := m.DB.Collection("users").Find(ctx, filter)
	defer cur.Close(ctx)
	if err != nil {
		return UsersResult{Err: err}
	}

	var users []User
	for cur.Next(ctx) {
		user := User{}
		err = cur.Decode(&user)
		if err != nil {
			break
		}
		users = append(users, user)
	}

	if err != nil {
		return UsersResult{Err: err}
	}

	return UsersResult{Users: users, Err: nil}
}
