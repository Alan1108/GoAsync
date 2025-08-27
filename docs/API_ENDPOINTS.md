# Endpoints de la API GoAsync

## Endpoints de Salud

### GET /health

Verifica el estado básico de salud de la API.

**Respuesta:**

```json
{
  "status": "ok",
  "message": "API funcionando correctamente",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0"
}
```

### GET /health/detailed

Proporciona información detallada del estado de salud de todos los servicios.

**Respuesta:**

```json
{
  "status": "ok",
  "message": "API funcionando correctamente",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0",
  "services": {
    "api": {
      "status": "healthy",
      "uptime": "running"
    },
    "database": {
      "status": "not_configured",
      "message": "Base de datos no configurada aún"
    },
    "redis": {
      "status": "not_configured",
      "message": "Redis no configurado aún"
    }
  }
}
```

## API v1

### GET /api/v1/

Información general de la API.

**Respuesta:**

```json
{
  "message": "Bienvenido a la API GoAsync",
  "version": "1.0.0",
  "status": "active",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### GET /api/v1/status

Estado detallado de la API con información de endpoints y características.

**Respuesta:**

```json
{
  "status": "operational",
  "uptime": "running",
  "version": "1.0.0",
  "endpoints": [
    "GET /health",
    "GET /health/detailed",
    "GET /api/v1/",
    "GET /api/v1/status"
  ],
  "features": [
    "REST API",
    "Health checks",
    "Structured logging",
    "Environment configuration"
  ]
}
```

## Códigos de Estado HTTP

- `200 OK`: Request exitoso
- `204 No Content`: Request exitoso sin contenido (para OPTIONS)
- `400 Bad Request`: Request mal formado
- `404 Not Found`: Recurso no encontrado
- `500 Internal Server Error`: Error interno del servidor

## Headers de Respuesta

- `Content-Type: application/json`
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: POST, OPTIONS, GET, PUT, DELETE, PATCH`
- `Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With`

## Ejemplos de Uso

### cURL

```bash
# Health check básico
curl http://localhost:8080/health

# Health check detallado
curl http://localhost:8080/health/detailed

# Información de la API
curl http://localhost:8080/api/v1/

# Estado de la API
curl http://localhost:8080/api/v1/status
```

### JavaScript (Fetch)

```javascript
// Health check
fetch("http://localhost:8080/health")
  .then((response) => response.json())
  .then((data) => console.log(data));

// API info
fetch("http://localhost:8080/api/v1/")
  .then((response) => response.json())
  .then((data) => console.log(data));
```

## Notas de Implementación

- Todos los endpoints retornan JSON
- Los timestamps están en formato ISO 8601
- La API soporta CORS para desarrollo
- Los logs se generan automáticamente para cada request
- El puerto por defecto es 8080 (configurable via variable de entorno)
