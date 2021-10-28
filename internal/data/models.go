package data

import "go.mongodb.org/mongo-driver/mongo"

type Models struct {
	Board BoardModel
}

func NewModels(db *mongo.Client) Models {
	return Models{
		Board: BoardModel{DB: db.Database("grello")},
	}
}
