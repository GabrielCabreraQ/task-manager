package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Representación de una tarea en la base de datos
type Task struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" validate:"required" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Completed   bool               `json:"completed" bson:"completed"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdateAt    time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	Tags        []string           `json:"tags,omitempty" bson:"tags,omitempty"`
}

// Para la creación de tareas
type CreateTask struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// Para la paginación de tareas
type TaskList struct {
	Tasks []Task `json:"tasks"`
	Total int64  `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

// Para la actualización de tareas
type UpdateTask struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Completed   *bool    `json:"completed,omitempty"`
}
