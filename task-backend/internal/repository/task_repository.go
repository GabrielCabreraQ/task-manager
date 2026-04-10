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

// TaskRepository define las operaciones de persistencia para las tareas.
type TaskRepository interface {
	CreateTask(ctx context.Context, task *model.Task) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Task, error)
	MarkCompleted(ctx context.Context, id primitive.ObjectID) (*model.Task, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	UpdateTask(ctx context.Context, id primitive.ObjectID, update *model.UpdateTask) error
	FindByTag(ctx context.Context, tag string, page, limit int) ([]model.Task, int64, error)
	FindAll(ctx context.Context, page, limit int, completed *bool) ([]model.Task, int64, error)
}

type taskRepository struct {
	collection *mongo.Collection
}

// NewTaskRepository inicializa una nueva implementación de TaskRepository.
func NewTaskRepository(collection *mongo.Collection) *taskRepository {
	return &taskRepository{collection: collection}
}

// CreateTask inserta una nueva tarea en la colección, establece las marcas de tiempo y actualiza el ID generado.
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

// FindByID busca y devuelve la tarea con el ObjectID en la base de datos.
func (t *taskRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Task, error) {
	var task model.Task
	err := t.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Delete elimina la tarea que coincide con el ObjectID proporcionado.
func (t *taskRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := t.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// MarkCompleted marca la tarea como completada y devuelve la tarea actualizada.
func (t *taskRepository) MarkCompleted(ctx context.Context, id primitive.ObjectID) (*model.Task, error) {
	_, err := t.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"completed": true, "updatedAt": time.Now()}},
	)
	if err != nil {
		return nil, err
	}
	return t.FindByID(ctx, id)
}

// FindAll devuelve un listado paginado de tareas con filtro opcional por estado.
func (t *taskRepository) FindAll(ctx context.Context, page, limit int, completed *bool) ([]model.Task, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := int64((page - 1) * limit)

	// Construir filtro — si completed es nil se traen todas las tareas
	filter := bson.D{}
	if completed != nil {
		filter = bson.D{{Key: "completed", Value: *completed}}
	}

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

	// El total respeta el mismo filtro para que la paginación sea correcta
	total, err := t.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// FindByTag busca tareas que contienen la etiqueta indicada y devuelve resultados paginados junto con el total.
func (t *taskRepository) FindByTag(ctx context.Context, tag string, page, limit int) ([]model.Task, int64, error) {

	// Normaliza los parámetros de paginación.
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Calcula el desplazamiento para la página solicitada.
	skip := int64((page - 1) * limit)

	filter := bson.M{"tags": tag}

	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	// Obtiene las tareas que contienen la etiqueta dada.
	cursor, err := t.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, 0, err
	}

	// Devuelve el total de tareas que coinciden con la etiqueta para paginación.
	total, _ := t.collection.CountDocuments(ctx, filter)
	return tasks, total, nil
}

// UpdateTask actualiza los campos proporcionados de una tarea existente.
func (t *taskRepository) UpdateTask(ctx context.Context, id primitive.ObjectID, update *model.UpdateTask) error {
	updateFields := bson.M{
		"updatedAt": time.Now(),
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
