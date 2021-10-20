package data

import "go.mongodb.org/mongo-driver/mongo"

type BoardModel struct {
	DB *mongo.Client
}
