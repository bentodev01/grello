package data

import "go.mongodb.org/mongo-driver/mongo"

type Models struct {
	Board BoardModel
	User  UserModel
	Task  TaskModel
}

func NewModels(db *mongo.Client) Models {
	DB := db.Database("grello")
	return Models{
		Board: BoardModel{DB: DB},
		User:  UserModel{DB: DB},
		Task:  TaskModel{DB: DB},
	}
}
