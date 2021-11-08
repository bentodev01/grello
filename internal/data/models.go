package data

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	Board interface {
		Insert(ctx context.Context, name string, description string, userIds []string) <-chan BoardResult
		Get(ctx context.Context, id string) BoardResult
		AddTask(ctx context.Context, id string, taskId string) error
	}
	User interface {
		Get(ctx context.Context, id string) (User, error)
		GetAll(ctx context.Context, ids []string) <-chan UsersResult
	}
	Task interface {
		Insert(ctx context.Context, name, description, boardId, assignedTo string) TaskResult
		GetAll(ctx context.Context, ids []string) <-chan TasksResult
	}
}

func NewModels(db *mongo.Client) Models {
	DB := db.Database("grello")
	return Models{
		Board: BoardModel{DB: DB},
		User:  UserModel{DB: DB},
		Task:  TaskModel{DB: DB},
	}
}
