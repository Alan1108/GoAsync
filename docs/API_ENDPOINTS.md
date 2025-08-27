# API Endpoints - GoAsync

Esta documentación describe todos los endpoints disponibles en la API de GoAsync.

## Base URL

```
http://localhost:8080/api/v1
```

## Autenticación

Actualmente la API no requiere autenticación. En futuras versiones se implementará JWT.

## Endpoints

### Health Check

- **GET** `/health` - Verifica el estado de la aplicación y la base de datos

### Usuarios

#### Obtener usuarios

- **GET** `/users` - Lista todos los usuarios con paginación
  - Query params:
    - `page` (int, default: 1) - Número de página
    - `per_page` (int, default: 10, max: 10000) - Elementos por página

#### Obtener usuario específico

- **GET** `/users/{id}` - Obtiene un usuario por su ID
- **GET** `/users/{id}/profile` - Obtiene un usuario con su perfil completo
- **GET** `/users/{id}/stats` - Obtiene estadísticas de un usuario
- **GET** `/users/{id}/activity` - Obtiene actividad reciente de un usuario
  - Query params:
    - `limit` (int, default: 10, max: 10000) - Número de actividades a retornar

#### Estadísticas de usuarios

- **GET** `/users/stats` - Obtiene estadísticas de todos los usuarios

#### Crear y gestionar usuarios

- **POST** `/users` - Crea un nuevo usuario
- **PUT** `/users/{id}` - Actualiza un usuario existente
- **DELETE** `/users/{id}` - Elimina un usuario

### Posts

#### Obtener posts

- **GET** `/posts` - Lista todos los posts con filtros y paginación

  - Query params:
    - `page` (int, default: 1) - Número de página
    - `per_page` (int, default: 10, max: 10000) - Elementos por página
    - `status` (string) - Filtrar por estado (draft, published, archived)
    - `search` (string) - Buscar en título, contenido y extracto
    - `category_id` (uuid) - Filtrar por categoría
    - `author_id` (uuid) - Filtrar por autor
    - `tag_id` (uuid) - Filtrar por tag

- **GET** `/posts/published` - Lista solo posts publicados
- **GET** `/posts/{id}` - Obtiene un post por su ID
- **GET** `/posts/slug/{slug}` - Obtiene un post por su slug
- **GET** `/posts/{id}/with-tags` - Obtiene un post con sus tags

#### Crear y gestionar posts

- **POST** `/posts` - Crea un nuevo post
- **PUT** `/posts/{id}` - Actualiza un post existente
- **DELETE** `/posts/{id}` - Elimina un post

### Categorías

#### Obtener categorías

- **GET** `/categories` - Lista todas las categorías activas
- **GET** `/categories/{id}` - Obtiene una categoría por su ID
- **GET** `/categories/slug/{slug}` - Obtiene una categoría por su slug
- **GET** `/categories/{id}/with-posts` - Obtiene una categoría con sus posts

#### Crear y gestionar categorías

- **POST** `/categories` - Crea una nueva categoría
- **PUT** `/categories/{id}` - Actualiza una categoría existente
- **DELETE** `/categories/{id}` - Elimina una categoría

### Tags

#### Obtener tags

- **GET** `/tags` - Lista todos los tags
- **GET** `/tags/popular` - Lista los tags más populares
  - Query params:
    - `limit` (int, default: 10, max: 10000) - Número de tags a retornar
- **GET** `/tags/{id}` - Obtiene un tag por su ID
- **GET** `/tags/slug/{slug}` - Obtiene un tag por su slug
- **GET** `/tags/{id}/with-posts` - Obtiene un tag con sus posts

#### Crear y gestionar tags

- **POST** `/tags` - Crea un nuevo tag
- **PUT** `/tags/{id}` - Actualiza un tag existente
- **DELETE** `/tags/{id}` - Elimina un tag

### Comentarios

#### Obtener comentarios

- **GET** `/comments` - Lista todos los comentarios
  - Query params:
    - `page` (int, default: 1) - Número de página
    - `per_page` (int, default: 10, max: 10000) - Elementos por página
    - `approved_only` (bool, default: true) - Solo comentarios aprobados
- **GET** `/comments/{id}` - Obtiene un comentario por su ID
- **GET** `/posts/{post_id}/comments` - Obtiene comentarios de un post específico

#### Crear y gestionar comentarios

- **POST** `/comments` - Crea un nuevo comentario
- **PUT** `/comments/{id}` - Actualiza un comentario existente
- **DELETE** `/comments/{id}` - Elimina un comentario
- **PATCH** `/comments/{id}/approve` - Aprueba un comentario

### Estadísticas

#### Estadísticas generales

- **GET** `/stats/database` - Obtiene estadísticas generales de la base de datos
- **GET** `/stats/posts` - Obtiene estadísticas de posts

#### Logs de actividad

- **GET** `/stats/activity` - Obtiene logs de actividad con filtros

  - Query params:
    - `page` (int, default: 1) - Número de página
    - `per_page` (int, default: 10, max: 10000) - Elementos por página
    - `user_id` (uuid) - Filtrar por usuario
    - `action` (string) - Filtrar por acción
    - `resource_type` (string) - Filtrar por tipo de recurso
    - `start_date` (date) - Fecha de inicio (YYYY-MM-DD)
    - `end_date` (date) - Fecha de fin (YYYY-MM-DD)

- **GET** `/stats/activity/recent` - Obtiene actividad reciente
  - Query params:
    - `limit` (int, default: 10, max: 10000) - Número de actividades a retornar
- **GET** `/stats/activity/user/{user_id}` - Obtiene actividad de un usuario específico
- **GET** `/stats/daily` - Obtiene estadísticas diarias
  - Query params:
    - `days` (int, default: 7) - Número de días a consultar

## Códigos de Respuesta

- **200** - OK - Operación exitosa
- **201** - Created - Recurso creado exitosamente
- **400** - Bad Request - Datos de entrada inválidos
- **404** - Not Found - Recurso no encontrado
- **409** - Conflict - Conflicto (ej: slug duplicado)
- **500** - Internal Server Error - Error interno del servidor
- **503** - Service Unavailable - Servicio no disponible

## Formato de Respuesta

Todas las respuestas siguen el siguiente formato:

```json
{
  "data": {...},
  "message": "Mensaje descriptivo",
  "pagination": {
    "page": 1,
    "per_page": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

## Paginación

Los endpoints que soportan paginación incluyen información de paginación en la respuesta:

```json
{
  "pagination": {
    "page": 1,
    "per_page": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

## Filtros

Muchos endpoints soportan filtros a través de query parameters:

- **Búsqueda**: `?search=texto`
- **Filtros específicos**: `?status=published&category_id=uuid`
- **Ordenamiento**: Implementado internamente (más reciente primero)

## Ejemplos de Uso

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

### Obtener estadísticas de la base de datos

```bash
curl -X GET "http://localhost:8080/api/v1/stats/database"
```

## Notas

- Todos los IDs son UUIDs
- Las fechas están en formato ISO 8601
- Los slugs son URLs amigables generados automáticamente
- Los comentarios requieren aprobación antes de ser visibles públicamente
- Los logs de actividad se generan automáticamente para todas las operaciones CRUD
