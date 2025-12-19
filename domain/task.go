package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	CollectionTask = "tasks"
)

type Task struct {
	ID     bson.ObjectID `bson:"_id" json:"-"`
	Title  string        `bson:"title" form:"title" binding:"required" json:"title"`
	UserID bson.ObjectID `bson:"userID" json:"-"`
}

type TaskRepository interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
}

type TaskUsecase interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
}
