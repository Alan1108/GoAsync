# GoAsync - API de Blog y Gestión de Contenido

GoAsync es una API RESTful completa desarrollada en Go para la gestión de blogs, usuarios, posts, categorías, tags y comentarios. La aplicación utiliza PostgreSQL como base de datos y Gin como framework web.

## 🚀 Características

- **API RESTful completa** con endpoints para todas las entidades
- **Base de datos PostgreSQL** con esquema optimizado
- **Paginación** en todos los endpoints de listado
- **Filtros avanzados** para búsqueda y filtrado
- **Logs de actividad** automáticos para auditoría
- **Estadísticas** en tiempo real
- **Documentación completa** de la API
- **Docker** para desarrollo y despliegue
- **Middleware** para CORS y logging

## 📋 Entidades Principales

- **Usuarios**: Gestión completa de usuarios y perfiles
- **Posts**: Artículos con soporte para categorías y tags
- **Categorías**: Organización de contenido por temas
- **Tags**: Etiquetado flexible de contenido
- **Comentarios**: Sistema de comentarios con moderación
- **Estadísticas**: Métricas y logs de actividad

## 🛠️ Tecnologías

- **Go 1.21+**
- **Gin** - Framework web
- **PostgreSQL** - Base de datos
- **Docker & Docker Compose** - Contenedores
- **Logrus** - Logging estructurado
- **UUID** - Identificadores únicos

## 📦 Instalación

### Prerrequisitos

- Go 1.21 o superior
- Docker y Docker Compose
- Git

### Clonar el repositorio

```bash
git clone <repository-url>
cd GoAsync
```

### Configuración con Docker (Recomendado)

1. **Iniciar la base de datos:**

```bash
docker-compose up -d postgres
```

2. **Configurar variables de entorno:**

```bash
cp .env.example .env
# Editar .env con tus configuraciones
```

3. **Ejecutar migraciones y seeders:**

```bash
# Las migraciones se ejecutan automáticamente al iniciar el contenedor
# Para ejecutar seeders manualmente:
docker-compose exec postgres psql -U postgres -d goasync -f /docker-entrypoint-initdb.d/seed.sql
```

4. **Compilar y ejecutar:**

```bash
go mod tidy
go build -o goasync .
./goasync
```

### Configuración Manual

1. **Instalar dependencias:**

```bash
go mod tidy
```

2. **Configurar PostgreSQL:**

- Crear base de datos `goasync`
- Ejecutar scripts en `docker/postgres/init.sql`
- Ejecutar seeders en `docker/postgres/seed.sql`

3. **Configurar variables de entorno:**

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=goasync
export DB_USER=postgres
export DB_PASSWORD=password
export PORT=8080
```

4. **Ejecutar la aplicación:**

```bash
go run main.go
```

## 🌐 Endpoints de la API

La API está disponible en `http://localhost:8080/api/v1`

### Endpoints Principales

#### Health Check

- `GET /health` - Estado de la aplicación

#### Usuarios

- `GET /users` - Listar usuarios
- `GET /users/{id}` - Obtener usuario
- `POST /users` - Crear usuario
- `PUT /users/{id}` - Actualizar usuario
- `DELETE /users/{id}` - Eliminar usuario

#### Posts

- `GET /posts` - Listar posts
- `GET /posts/published` - Posts publicados
- `GET /posts/{id}` - Obtener post
- `POST /posts` - Crear post
- `PUT /posts/{id}` - Actualizar post
- `DELETE /posts/{id}` - Eliminar post

#### Categorías

- `GET /categories` - Listar categorías
- `GET /categories/{id}` - Obtener categoría
- `POST /categories` - Crear categoría
- `PUT /categories/{id}` - Actualizar categoría
- `DELETE /categories/{id}` - Eliminar categoría

#### Tags

- `GET /tags` - Listar tags
- `GET /tags/popular` - Tags populares
- `POST /tags` - Crear tag
- `PUT /tags/{id}` - Actualizar tag
- `DELETE /tags/{id}` - Eliminar tag

#### Comentarios

- `GET /comments` - Listar comentarios
- `GET /posts/{post_id}/comments` - Comentarios de un post
- `POST /comments` - Crear comentario
- `PATCH /comments/{id}/approve` - Aprobar comentario

#### Estadísticas

- `GET /stats/database` - Estadísticas generales
- `GET /stats/activity` - Logs de actividad
- `GET /stats/posts` - Estadísticas de posts

## 📖 Documentación Completa

Consulta la documentación detallada de la API en:

- [API Endpoints](docs/API_ENDPOINTS.md)
- [Base de Datos](docs/DATABASE.md)

## 🧪 Ejemplos de Uso

### Obtener posts publicados

```bash
curl -X GET "http://localhost:8080/api/v1/posts/published?page=1&per_page=5"
```

### Crear un nuevo post

```bash
curl -X POST "http://localhost:8080/api/v1/posts" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Mi primer post",
    "content": "Contenido del post...",
    "excerpt": "Resumen del post",
    "category_id": "uuid-de-categoria",
    "status": "published",
    "tag_ids": ["uuid-tag-1", "uuid-tag-2"]
  }'
```

### Obtener estadísticas

```bash
curl -X GET "http://localhost:8080/api/v1/stats/database"
```

## 🔧 Desarrollo

### Estructura del Proyecto

```
GoAsync/
├── cmd/                    # Comandos de la aplicación
├── docker/                 # Configuración de Docker
├── docs/                   # Documentación
├── internal/               # Código interno de la aplicación
│   ├── config/            # Configuración
│   ├── handlers/          # Manejadores HTTP
│   ├── models/            # Modelos de datos
│   └── services/          # Lógica de negocio
├── pkg/                   # Paquetes públicos
│   ├── database/          # Utilidades de base de datos
│   ├── logger/            # Sistema de logging
│   └── middleware/        # Middleware HTTP
├── scripts/               # Scripts de utilidad
├── main.go               # Punto de entrada
├── go.mod                # Dependencias Go
└── docker-compose.yml    # Configuración de Docker Compose
```

### Comandos Útiles

```bash
# Compilar
go build -o goasync .

# Ejecutar tests
go test ./...

# Ejecutar con hot reload (requiere air)
air

# Generar documentación
godoc -http=:6060

# Limpiar
go clean
```

### Variables de Entorno

```bash
# Servidor
PORT=8080
GIN_MODE=debug

# Base de datos
DB_HOST=localhost
DB_PORT=5432
DB_NAME=goasync
DB_USER=postgres
DB_PASSWORD=password
DB_SSLMODE=disable

# Logging
LOG_LEVEL=info
```

## 🐳 Docker

### Construir imagen

```bash
docker build -t goasync .
```

### Ejecutar con Docker Compose

```bash
docker-compose up -d
```

### Ver logs

```bash
docker-compose logs -f
```

## 📊 Base de Datos

La aplicación incluye:

- **Esquema completo** con todas las tablas necesarias
- **Índices optimizados** para consultas rápidas
- **Triggers** para actualización automática de timestamps
- **Funciones** para estadísticas y utilidades
- **Vistas** para consultas complejas
- **Datos de ejemplo** para testing

### Migraciones

Las migraciones se ejecutan automáticamente al iniciar el contenedor de PostgreSQL.

### Seeders

Los datos de ejemplo incluyen:

- 5 usuarios con perfiles
- 8 categorías
- 10 tags
- 4 posts con contenido completo
- Comentarios de ejemplo
- Logs de actividad

## 🔒 Seguridad

- **Validación de entrada** en todos los endpoints
- **Sanitización** de datos
- **Logs de auditoría** para todas las operaciones
- **Manejo de errores** estructurado
- **CORS** configurado para desarrollo

## 🚀 Despliegue

### Producción

1. **Configurar variables de entorno de producción**
2. **Usar PostgreSQL de producción**
3. **Configurar logging apropiado**
4. **Implementar autenticación JWT**
5. **Configurar HTTPS**

### Monitoreo

- **Health checks** automáticos
- **Logs estructurados** con Logrus
- **Métricas** de base de datos
- **Estadísticas** de actividad

## 🤝 Contribuir

1. Fork el proyecto
2. Crear una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir un Pull Request

## 📝 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 🆘 Soporte

Si tienes problemas o preguntas:

1. Revisa la documentación
2. Busca en los issues existentes
3. Crea un nuevo issue con detalles del problema

## 🎯 Roadmap

- [ ] Autenticación JWT
- [ ] Subida de archivos
- [ ] API GraphQL
- [ ] Cache con Redis
- [ ] Tests automatizados
- [ ] CI/CD pipeline
- [ ] Documentación con Swagger
- [ ] Dashboard de administración
