package router

import (
	"task-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter configura todas las rutas de la API
func SetupRouter(h *handler.TaskHandler) *gin.Engine {
	r := gin.Default()

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

	// Uso de rutas para las consultas
	tasks := r.Group("/tasks")
	{
		tasks.POST("", h.CreateTask)                // POST /tasks
		tasks.GET("", h.GetAllTasks)                // GET /tasks
		tasks.GET("/:id", h.GetTaskByID)            // GET /tasks/{id}
		tasks.PUT("/:id/complete", h.MarkCompleted) // PUT /tasks/{id}/complete
		tasks.DELETE("/:id", h.DeleteTask)          // DELETE /tasks/{id}
		tasks.PUT("/:id", h.UpdateTask)             // PUT /tasks/{id}
		tasks.GET("/tag/:tag", h.FindByTag)         // GET /tasks/tag/{tag}
	}

	return r
}
