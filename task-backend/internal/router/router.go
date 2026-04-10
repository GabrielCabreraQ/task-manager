package router

import (
	"task-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter configura todas las rutas de la API y habilita CORS básico.
func SetupRouter(h *handler.TaskHandler) *gin.Engine {
	r := gin.Default()

	// Middleware que permite peticiones desde el frontend y responde al OPTIONS.
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Grupo de rutas para la gestión de tareas.
	tasks := r.Group("/tasks")
	{
		tasks.POST("", h.CreateTask)                // Crea una nueva tarea.
		tasks.GET("", h.GetAllTasks)                // Obtiene todas las tareas.
		tasks.GET("/:id", h.GetTaskByID)            // Obtiene una tarea por su ID.
		tasks.GET("/:id/complete", h.MarkCompleted) // para usarlo en dirección 8080
		tasks.PUT("/:id/complete", h.MarkCompleted) // Marca una tarea como completada.
		tasks.DELETE("/:id", h.DeleteTask)          // Elimina una tarea.
		tasks.PUT("/:id", h.UpdateTask)             // Actualiza una tarea existente.
		tasks.GET("/tag/:tag", h.FindByTag)         // Busca tareas por etiqueta.
	}

	return r
}
