package service

import (
	"context"
	"errors"

	"task-backend/internal/model"
	"task-backend/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService struct {
	repository repository.TaskRepository
}

func NewTaskService(repository repository.TaskRepository) *TaskService {
	return &TaskService{repository: repository}
}

// CreateTask crea una nueva tarea
func (s *TaskService) CreateTask(ctx context.Context, newTask model.CreateTask) (*model.Task, error) {
	if len(newTask.Title) < 3 {
		return nil, errors.New("el título debe tener al menos 3 caracteres")
	}

	task := &model.Task{
		Title:       newTask.Title,
		Description: newTask.Description,
		Tags:        newTask.Tags,
		Completed:   false,
	}

	err := s.repository.CreateTask(ctx, task)
	if err != nil {
		return nil, errors.New("error al guardar la tarea en la base de datos")
	}

	return task, nil
}

// UpdateTask actualiza una tarea existente
func (s *TaskService) UpdateTask(ctx context.Context, id string, update model.UpdateTask) (*model.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID de tarea inválido")
	}

	err = s.repository.UpdateTask(ctx, objectID, &update)
	if err != nil {
		return nil, errors.New("error al actualizar la tarea")
	}

	// Retornamos la tarea actualizada
	return s.FindByID(ctx, id)
}

// FindByTag busca tareas por una etiqueta específica
func (s *TaskService) FindByTag(ctx context.Context, tag string, page, limit int) (model.TaskList, error) {
	tasks, total, err := s.repository.FindByTag(ctx, tag, page, limit)
	if err != nil {
		return model.TaskList{}, errors.New("error al buscar tareas por etiqueta")
	}

	return model.TaskList{
		Tasks: tasks,
		Total: total,
	}, nil
}

// FindAll obtiene todas las tareas con paginación
func (s *TaskService) FindAll(ctx context.Context, page, limit int) (model.TaskList, error) {
	tasks, total, err := s.repository.FindAll(ctx, page, limit)
	if err != nil {
		return model.TaskList{}, errors.New("error al obtener las tareas")
	}

	return model.TaskList{
		Tasks: tasks,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// FindByID obtiene una tarea por ID
func (s *TaskService) FindByID(ctx context.Context, id string) (*model.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID de tarea inválido")
	}

	task, err := s.repository.FindByID(ctx, objectID)
	if err != nil {
		return nil, errors.New("tarea no encontrada")
	}

	return task, nil
}

// MarkCompleted marca una tarea como completada
func (s *TaskService) MarkCompleted(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID de tarea inválido")
	}

	err = s.repository.MarkCompleted(ctx, objectID)
	if err != nil {
		return errors.New("error al actualizar la tarea")
	}

	return nil
}

// Delete elimina una tarea
func (s *TaskService) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID de tarea inválido")
	}

	err = s.repository.Delete(ctx, objectID)
	if err != nil {
		return errors.New("error al eliminar la tarea")
	}

	return nil
}
