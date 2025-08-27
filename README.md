# GoAsync API

Una API REST moderna construida en Go con Gin framework, diseñada para ser escalable y fácil de mantener.

## 🚀 Características

- **Framework**: Gin para routing HTTP
- **Logging**: Logrus para logging estructurado
- **Configuración**: Variables de entorno con godotenv
- **Base de Datos**: PostgreSQL 15 con esquema completo
- **Containerización**: Docker y Docker Compose
- **Hot Reload**: Air para desarrollo con recarga automática
- **Estructura**: Arquitectura limpia y organizada
- **Seeder**: Datos de ejemplo incluidos

## 📋 Prerrequisitos

- Go 1.21 o superior
- Docker y Docker Compose (opcional)
- PostgreSQL 15 (para instalación local)
- Make (opcional, para usar comandos del Makefile)

## 🛠️ Instalación

### Opción 1: Instalación local

1. Clona el repositorio:

```bash
git clone <tu-repositorio>
cd goasync
```

2. Instala las dependencias:

```bash
make deps
# o manualmente:
go mod tidy
go mod download
```

3. Copia el archivo de variables de entorno:

```bash
cp env.example .env
```

4. Configura la base de datos PostgreSQL local o usa Docker:

```bash
# Con Docker (recomendado)
make db-up

# O instala PostgreSQL localmente y crea la base de datos
```

5. Ejecuta el seeder para poblar la base de datos:

```bash
make seed
```

6. Ejecuta la aplicación:

```bash
make run
# o manualmente:
go run main.go
```

### Opción 2: Con Docker

1. Construye y ejecuta con Docker Compose:

```bash
make docker-compose-up
```

2. O solo la aplicación:

```bash
make docker-build
make docker-run
```

## 🗄️ Base de Datos

### Estructura Incluida

- **8 categorías** (Tecnología, Ciencia, Salud, etc.)
- **5 usuarios** con perfiles completos
- **10 tags** relacionados con tecnología
- **4 posts** de ejemplo con contenido completo
- **Sistema de comentarios** y relaciones
- **Logs de actividad** de ejemplo

### Comandos de Base de Datos

```bash
make db-up          # Levantar servicios de BD
make db-down        # Detener servicios
make db-reset       # Reiniciar BD
make db-logs        # Ver logs
make db-connect     # Conectar a PostgreSQL
```

### Seeder

```bash
make seed           # Ejecutar seeder local
make seed-docker    # Ejecutar seeder en Docker
make seed-clean     # Limpiar y ejecutar seeder
```

Para más detalles sobre la base de datos, consulta [DATABASE.md](docs/DATABASE.md).

## 🌱 Seeder de Base de Datos

El proyecto incluye un seeder completo para poblar la base de datos con datos de ejemplo:

### Modos de Seeder

#### 🔥 Seeder Masivo (Recomendado para pruebas de rendimiento)

```bash
# Ejecutar seeder masivo localmente
make seed-massive

# Ejecutar seeder masivo en Docker
make seed-massive-docker

# Usar script directamente
./scripts/seed-db.sh --massive
```

**Genera aproximadamente 60,000+ registros:**

- 1,000 usuarios con nombres realistas
- 15 categorías de contenido
- 100+ tags tecnológicos y generales
- 5,000 posts de ejemplo
- 15,000 comentarios
- 25,000 relaciones post-tag
- 10,000 logs de actividad

#### 📝 Seeder Básico (Para desarrollo rápido)

```bash
# Ejecutar seeder básico localmente
make seed-small

# Ejecutar seeder básico en Docker
make seed-small-docker

# Usar script directamente
./scripts/seed-db.sh --small
```

**Genera aproximadamente 30 registros:**

- 5 usuarios básicos
- 8 categorías principales
- 10 tags esenciales
- 4 posts de ejemplo
- Comentarios básicos

#### ⚡ Seeder por Defecto

```bash
# Ejecutar seeder por defecto
make seed

# Ejecutar en Docker
make seed-docker

# Usar script directamente
./scripts/seed-db.sh
```

### Comandos del Seeder

```bash
# Comandos principales
make seed              # Seeder básico local
make seed-docker       # Seeder básico en Docker
make seed-massive      # Seeder masivo local
make seed-massive-docker # Seeder masivo en Docker
make seed-small        # Seeder pequeño local
make seed-small-docker # Seeder pequeño en Docker

# Con limpieza de base de datos
make seed-clean        # Limpiar DB y ejecutar seeder básico

# Script bash con opciones avanzadas
./scripts/seed-db.sh --help           # Ver todas las opciones
./scripts/seed-db.sh --massive        # Seeder masivo
./scripts/seed-db.sh --small          # Seeder pequeño
./scripts/seed-db.sh --docker         # Ejecutar en Docker
./scripts/seed-db.sh --clean          # Limpiar DB primero
./scripts/seed-db.sh --verbose        # Modo verbose
```

### Características del Seeder

- **Inserción en lotes**: Optimizado para insertar miles de registros eficientemente
- **Datos realistas**: Nombres, emails y contenido que simulan un entorno real
- **Relaciones coherentes**: Mantiene integridad referencial entre entidades
- **Logging detallado**: Progreso en tiempo real con emojis y estadísticas
- **Manejo de errores**: Recuperación robusta y mensajes informativos
- **Flexibilidad**: Múltiples modos para diferentes necesidades

### Casos de Uso

#### 🚀 Para Pruebas de Rendimiento

```bash
make seed-massive
```

Ideal para probar:

- Consultas complejas con grandes volúmenes de datos
- Rendimiento de índices y optimizaciones
- Escalabilidad de la aplicación
- Estrategias de paginación

#### 🧪 Para Desarrollo Rápido

```bash
make seed-small
```

Perfecto para:

- Desarrollo y debugging
- Pruebas unitarias
- Demostraciones
- Entornos de staging

#### 🔄 Para Reinicio Limpio

```bash
make seed-clean
```

Útil cuando:

- Cambias el esquema de la base de datos
- Quieres empezar desde cero
- Hay inconsistencias en los datos
- Cambias entre diferentes modos de seeder

## 🚀 Desarrollo

### Comandos útiles

```bash
make help          # Muestra todos los comandos disponibles
make dev           # Ejecuta con hot reload (Air)
make dev-full      # Levanta BD + seeder + aplicación
make dev-docker    # Todo en Docker
make test          # Ejecuta los tests
make lint          # Ejecuta el linter
make format        # Formatea el código
make clean         # Limpia archivos generados
```

### Hot Reload

Para desarrollo con recarga automática:

```bash
make dev
```

Esto instalará Air automáticamente y ejecutará la aplicación con hot reload.

## 📁 Estructura del proyecto

```
goasync/
├── main.go                    # Punto de entrada de la aplicación
├── go.mod                     # Dependencias de Go
├── Makefile                   # Comandos útiles
├── Dockerfile                 # Configuración de Docker
├── docker-compose.yml         # Orquestación de servicios
├── .air.toml                  # Configuración de Air (hot reload)
├── env.example                # Variables de entorno de ejemplo
├── docker/                    # Configuración de Docker
│   └── postgres/             # Dockerfile y scripts de PostgreSQL
│       ├── Dockerfile        # Dockerfile personalizado para PostgreSQL
│       ├── init.sql          # Script de inicialización de la BD
│       └── seed.sql          # Script de datos de ejemplo
├── cmd/                       # Comandos de la aplicación
│   └── seeder/               # Seeder de base de datos
│       └── main.go           # Seeder principal
├── scripts/                   # Scripts de utilidad
│   └── seed-db.sh            # Script bash para ejecutar seeder
├── internal/                  # Código interno de la aplicación
│   ├── config/               # Configuración y variables de entorno
│   ├── handlers/             # Handlers HTTP con tests
│   ├── models/               # Modelos de datos
│   └── services/             # Lógica de negocio
├── pkg/                       # Paquetes reutilizables
│   ├── database/             # Conexiones de base de datos
│   ├── logger/               # Sistema de logging personalizado
│   └── middleware/           # Middleware HTTP
├── docs/                      # Documentación
│   ├── API_ENDPOINTS.md      # Documentación de endpoints
│   └── DATABASE.md           # Documentación de base de datos
└── README.md                  # Este archivo
```

## 🌐 Endpoints disponibles

- `GET /health` - Verificación de salud de la API
- `GET /health/detailed` - Verificación detallada de salud
- `GET /api/v1/` - Información de la API v1
- `GET /api/v1/status` - Estado detallado de la API

## 🔧 Configuración

### Variables de entorno

Crea un archivo `.env` basado en `env.example`:

```bash
# Configuración del servidor
PORT=8080
GIN_MODE=debug

# Configuración de la base de datos
DB_HOST=localhost
DB_PORT=5432
DB_NAME=goasync
DB_USER=postgres
DB_PASSWORD=password
DB_SSLMODE=disable

# Configuración de Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Configuración de logs
LOG_LEVEL=info
```

### Puertos

- **API**: 8080 (configurable via PORT)
- **PostgreSQL**: 5432
- **Redis**: 6379
- **pgAdmin**: 5050

## 🧪 Testing

```bash
make test              # Tests básicos
make test-coverage     # Tests con cobertura
```

## 📦 Build

```bash
make build             # Compila la aplicación
```

El binario se generará en `bin/goasync`.

## 🐳 Docker

### Construir imagen

```bash
make docker-build
```

### Ejecutar contenedor

```bash
make docker-run
```

### Con Docker Compose

```bash
make docker-compose-up
```

### Servicios incluidos

- **goasync-app**: API principal
- **postgres**: Base de datos PostgreSQL
- **redis**: Cache y sesiones
- **pgadmin**: Interfaz web para PostgreSQL

## 📝 Logs

La aplicación usa Logrus para logging estructurado. Los logs se muestran en consola en modo desarrollo.

## 🔍 Monitoreo

### pgAdmin

Accede a pgAdmin en http://localhost:5050:

- **Email**: admin@goasync.com
- **Contraseña**: admin123

### Health Checks

Los contenedores incluyen health checks automáticos para PostgreSQL y Redis.

## 🚨 Solución de Problemas

### Problemas Comunes

1. **Puerto 5432 ocupado**: Usa `make db-up` para levantar PostgreSQL en Docker
2. **Error de conexión**: Verifica que la base de datos esté ejecutándose
3. **Error de permisos**: Usa `make db-reset` para limpiar volúmenes

### Verificación de Estado

```bash
# Ver estado de contenedores
docker-compose ps

# Ver logs
make db-logs

# Conectar a base de datos
make db-connect
```

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 🆘 Soporte

Si tienes alguna pregunta o problema, por favor abre un issue en el repositorio.

## 📚 Documentación Adicional

- [Endpoints de la API](docs/API_ENDPOINTS.md)
- [Base de Datos](docs/DATABASE.md)
