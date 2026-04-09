package main

import (
	"context"
	"log"

	"task-backend/internal/config"
	"task-backend/internal/handler"
	"task-backend/internal/repository"
	"task-backend/internal/router"
	"task-backend/internal/service"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// Conectar a MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatalf("Error conectando a MongoDB: %v", err)
	}

	// Verificar que la conexión esté activa
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("No se pudo conectar a MongoDB: %v", err)
	}

	log.Println("Conectado exitosamente a MongoDB")

	// Inicializar base de datos y colección
	db := client.Database(cfg.MongoDB.DBName)
	collection := db.Collection("tasks") // ← colección más clara

	// Inyectar dependencias
	repo := repository.NewTaskRepository(collection)
	svc := service.NewTaskService(repo)
	h := handler.NewTaskHandler(svc)

	// Configurar router
	r := router.SetupRouter(h)

	// Iniciar servidor
	log.Printf("Servidor corriendo en http://localhost:%s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
