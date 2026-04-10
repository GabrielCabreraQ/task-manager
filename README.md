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

## Ejecutar el proyecto localmente solo con Docker

1. Desde la dirección (`task-manager/deployment`):

```bash
docker-compose up --build
```

2. Accede a la aplicación en el navegador:

- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`

## Ejecución local sin docker

Se necesitan 3 terminales abiertas en paralelo: Una para MongoDB, una para el backend y otra para el frontend.

### Terminal cmd para MongoDB

Ejecutar MongoDB desde la carpeta de instalación, por ejemplo:

```
"C:\Program Files\MongoDB\Server\8.0\bin\mongod.exe" --dbpath "C:\data\db"
```

Verifica que MongoDB está corriendo conectándote con `mongosh` en otra terminal, esto debería mostrar la consola de MongoDB sin errores. MongoDB queda escuchando en: 
```
mongodb://localhost:27017
```

### Terminal powershell para backend

1. Entra en la carpeta del backend desde la raiz `task-manager/`
```
cd task-backend
```

2. Crea el archivo `.env` y copia lo siguiente:

```env
PORT=8080
DB_URL=localhost:27017
DB_NAME=taskmanager
DB_USER=
DB_PASSWORD=
```
Guarda el archivo en la dirección: `task-manager/task-backend/`

3. Descarga las dependencias de go:
```
go mod tidy
```

4. Ejecuta el backend desde `task-backend/`:

```
go run ./cmd/main.go
```

El backend quedará corriendo en la url: `http://localhost:8080`

### Terminal powershell para frontend

1. Entra en la carpeta del frontend desde la raiz `task-manager/`

```
cd task-frontend
```

2. Instala dependencias desde `task-frontend/`:

```
npm install
```

3. Inicia el frontend:

```
npm run dev
```

El frontend quedará corriendo en la url `http://localhost:3000` en el navegador. La aplicación debe cargar y mostrar las tareas desde MongoDB.

## Endpoints de la API

| Método     | Endpoint                     | Descripción                                      | Ejemplo |
|------------|------------------------------|--------------------------------------------------|---------|
| `POST`     | `/tasks`                     | Crear una nueva tarea                            | Ver ejemplos abajo |
| `GET`      | `/tasks`                     | Obtener todas las tareas (paginadas)             | `?page=1&limit=50` |
| `GET`      | `/tasks/:id`                 | Obtener una tarea por ID                         | `/tasks/69d87f428eb9fcdb6041d329` |
| `PUT`      | `/tasks/:id`                 | Actualizar una tarea                             |  |
| `PUT`      | `/tasks/:id/complete`        | Marcar tarea como completada                     | Solo método PUT |
| `DELETE`   | `/tasks/:id`                 | Eliminar una tarea                               | `/tasks/69d87f428eb9fcdb6041d329` |
| `GET`      | `/tasks/tag/:tag`            | Buscar tareas por etiqueta                       | `/tasks/tag/trabajo` |

Ejemlo de creación de tarea
```
'{
    "title": "Terminar informe mensual",
    "description": "Incluye gráficos y resumen ejecutivo",
    "tags": ["trabajo", "urgente"]
  }'

```