# Task Manager

Aplicación web para la gestión de tareas diseñada para un entorno full-stack.

## Descripción del proyecto

Task Manager es una aplicación que permite crear, consultar, actualizar, eliminar y marcar tareas como completadas. El proyecto cuenta con una API REST en Go que persiste datos en MongoDB, y una interfaz de usuario construida con Nuxt 3 y Vue 3.

Este repositorio está preparado para ejecución local, desarrollo con Docker y despliegue en entornos controlados.

## Tecnologías utilizadas

- **Backend**: Go, Gin
- **Frontend**: Nuxt 3, Vue 3
- **Base de datos**: MongoDB
- **Contenedores**: Docker, Docker Compose
- **CI/CD**: Jenkins
- **Validación**: validator.v10
- **Configuración**: godotenv

## Características principales

- Creación de tareas con título, descripción y etiquetas
- Obtención de tareas paginadas
- Búsqueda de tareas por etiqueta
- Consulta de tarea individual por ID
- Actualización de datos de tareas
- Eliminación de tareas
- Marcado de tareas como completadas
- Configuración por variables de entorno
- Estructura separada de backend, frontend y base de datos

## Estructura del proyecto

```text
task-manager/
├── README.md
├── deployment/
│   ├── docker-compose.yml
│   └── Jenkinsfile
├── task-backend/
│   ├── .env.example
│   ├── .gitignore
│   ├── Dockerfile
│   ├── LICENSE
│   ├── README.md
│   ├── dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── config/
│       │   └── config.go
│       ├── handler/
│       │   └── task_handler.go
│       ├── model/
│       │   └── task.go
│       ├── repository/
│       │   └── task_repository.go
│       ├── router/
│       │   └── router.go
│       └── service/
│           └── task_service.go
├── task-frontend/
│   ├── .gitignore
│   ├── Dockerfile
│   ├── README.md
│   ├── nuxt.config.js
│   ├── package-lock.json
│   ├── package.json
│   ├── tsconfig.json
│   ├── assets/
│   │   └── css/
│   │       ├── main.scss
│   │       └── variables.scss
│   ├── components/
│   │   ├── ConfirmModal.vue
│   │   ├── TagInput.vue
│   │   ├── TagsSidebar.vue
│   │   ├── TaskCard.vue
│   │   ├── TaskModal.vue
│   │   └── ToastContainer.vue
│   ├── layouts/
│   │   └── default.vue
│   ├── pages/
│   │   ├── index.vue
│   │   └── tasks/
│   │       └── [id].vue
│   ├── plugins/
│   │   └── axios.js
│   └── store/
│       ├── tasks.js
│       └── toast.js
```

## Requisitos previos

- Docker y Docker Compose instalados
- Go 1.XX instalado (solo si ejecutas backend local sin Docker)
- Node.js instalado (solo si ejecutas frontend local sin Docker)
- MongoDB disponible localmente o por contenedor

## Clonar el repositorio

Si todavía no tienes el proyecto en tu máquina, clona el repositorio y entra en la carpeta raíz:

```bash
git clone https://github.com/GabrielCabreraQ/task-manager
cd task-manager
```

Si ya tienes el repositorio clonado, ve a `task-manager`.

## Ejecutar el proyecto localmente con Docker

1. Desde la dirección (`task-manager/deployment`):

```bash
docker-compose up --build
```

2. Accede a la aplicación en el navegador:

- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`

## Ejecutar backend localmente

1. Copia o crea un archivo `.env` dentro de `task-backend/` con las variables necesarias:

```env
PORT=8080
DB_URL=localhost:27017
DB_NAME=taskmanager
DB_USER=
DB_PASSWORD=
```
2. Se debe tener la aplicación de MongoDB abierto para realizar la conección directa entre el backend y la base de datos.

3. Ejecuta desde `task-backend/`:

```bash
go run ./cmd/main.go
```

## Ejecutar frontend localmente

1. Instala dependencias desde `task-frontend/`:

```bash
npm install
```

2. Inicia el frontend:

```bash
npm run dev
```

## Notas de implementación

- El backend organiza el código en capas: controlador, modelo, repositorio y servicio.
- El router usa middleware básico de CORS para permitir peticiones desde el frontend.
- La configuración se carga desde variables de entorno con soporte `.env`.
- Los modelos usan MongoDB `ObjectID` para identificar tareas.

