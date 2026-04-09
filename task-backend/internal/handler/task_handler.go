package handler

import (
	"task-backend/internal/service"

	"github.com/go-playground/validator/v10"

	"net/http"
	"task-backend/internal/model"

	"github.com/gin-gonic/gin"

	"strconv"
)

// TaskHandler expone los endpoints HTTP y usa TaskService para la lógica de negocio.
type TaskHandler struct {
	service   *service.TaskService
	validator *validator.Validate
}

// NewTaskHandler crea un controlador de tareas con la validación configurada.
func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{
		service:   service,
		validator: validator.New(),
	}
}

// CreateTask recibe una petición para crear una nueva tarea y valida el payload.
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req model.CreateTask

	// Intenta convertir el JSON recibido en el modelo de creación.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato JSON inválido"})
		return
	}

	// Valida las reglas definidas en el modelo.
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Llama al servicio para guardar la tarea.
	task, err := h.service.CreateTask(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetAllTasks devuelve la lista paginada de tareas.
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Seguridad: evitar valores inválidos.
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 { // máximo 100 tareas por página
		limit = 10
	}

	// Solicita al servicio todos los registros con paginación.
	result, err := h.service.FindAll(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTaskByID obtiene una tarea por su identificador.
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	id := c.Param("id")

	// Busca la tarea usando el servicio.
	task, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// MarkCompleted cambia el estado de una tarea a completada.
func (h *TaskHandler) MarkCompleted(c *gin.Context) {
	id := c.Param("id")

	// Marca la tarea como completada a través del servicio.
	err := h.service.MarkCompleted(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tarea marcada como completada correctamente"})
}

// DeleteTask elimina una tarea por ID.
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	// Elimina la tarea en el servicio.
	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tarea eliminada correctamente"})
}

// UpdateTask recibe los datos para modificar una tarea existente.
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var update model.UpdateTask

	// Vincula el JSON recibido con el modelo de actualización.
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato JSON inválido"})
		return
	}

	// Llama al servicio para aplicar los cambios.
	task, err := h.service.UpdateTask(c.Request.Context(), id, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// FindByTag busca tareas filtradas por etiqueta y devuelve resultados paginados.
func (h *TaskHandler) FindByTag(c *gin.Context) {
	tag := c.Param("tag")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Obtiene las tareas que coinciden con la etiqueta.
	result, err := h.service.FindByTag(c.Request.Context(), tag, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
