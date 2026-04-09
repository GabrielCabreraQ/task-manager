package repository

import (
	"context"
	"time"

	"task-backend/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *model.Task) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Task, error)
	MarkCompleted(ctx context.Context, id primitive.ObjectID) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	UpdateTask(ctx context.Context, id primitive.ObjectID, update *model.UpdateTask) error
	FindByTag(ctx context.Context, tag string, page, limit int) ([]model.Task, int64, error)
	FindAll(ctx context.Context, page int, limit int) ([]model.Task, int64, error)
}

type taskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection) *taskRepository {
	return &taskRepository{collection: collection}
}

func (t *taskRepository) CreateTask(ctx context.Context, task *model.Task) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	result, err := t.collection.InsertOne(ctx, task)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		task.ID = oid
	}

	return nil
}

func (t *taskRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Task, error) {
	var task model.Task
	err := t.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (t *taskRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := t.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (t *taskRepository) MarkCompleted(ctx context.Context, id primitive.ObjectID) error {
	_, err := t.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"completed": true}},
	)
	return err
}
func (t *taskRepository) FindAll(ctx context.Context, page, limit int) ([]model.Task, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := int64((page - 1) * limit)

	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := t.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, 0, err
	}

	total, err := t.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (t *taskRepository) UpdateTask(ctx context.Context, id primitive.ObjectID, update *model.UpdateTask) error {
	updateFields := bson.M{
		"updatedAt": time.Now(), // Registramos cuándo se modificó
	}

	if update.Title != "" {
		updateFields["title"] = update.Title
	}
	if update.Description != "" {
		updateFields["description"] = update.Description
	}
	if update.Tags != nil {
		updateFields["tags"] = update.Tags
	}
	if update.Completed != nil {
		updateFields["completed"] = *update.Completed
	}

	_, err := t.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updateFields},
	)
	return err
}

func (t *taskRepository) FindByTag(ctx context.Context, tag string, page, limit int) ([]model.Task, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := int64((page - 1) * limit)

	filter := bson.M{"tags": tag}

	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := t.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, 0, err
	}
	total, _ := t.collection.CountDocuments(ctx, filter)
	return tasks, total, nil
}
