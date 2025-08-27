# Base de Datos GoAsync

Esta documentación describe la configuración y uso de la base de datos PostgreSQL en el proyecto GoAsync.

## 🗄️ Configuración de la Base de Datos

### Especificaciones Técnicas

- **Motor**: PostgreSQL 15
- **Puerto**: 5432
- **Base de datos**: `goasync`
- **Usuario**: `postgres`
- **Contraseña**: `password`
- **Host**: `localhost` (local) / `postgres` (Docker)

### Extensiones Instaladas

- **uuid-ossp**: Para generación de UUIDs
- **pgcrypto**: Para funciones de criptografía (hash de contraseñas)

## 🚀 Inicio Rápido

### Opción 1: Con Docker (Recomendado)

```bash
# Levantar solo la base de datos
make db-up

# Levantar todos los servicios
make docker-compose-up

# Ver logs de la base de datos
make db-logs
```

### Opción 2: Instalación Local

1. Instalar PostgreSQL 15
2. Crear base de datos:

```sql
CREATE DATABASE goasync;
CREATE USER postgres WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE goasync TO postgres;
```

## 📊 Estructura de la Base de Datos

### Tablas Principales

#### `users`

- Almacena información de usuarios del sistema
- Contraseñas hasheadas con bcrypt
- Soporte para perfiles activos/inactivos

#### `user_profiles`

- Información adicional de usuarios
- Relación 1:1 con `users`

#### `categories`

- Categorías para organizar posts
- Sistema de slugs únicos
- Soporte para categorías activas/inactivas

#### `posts`

- Artículos y contenido del sistema
- Sistema de estados (draft, published, archived)
- Relaciones con usuarios y categorías

#### `tags`

- Etiquetas para categorizar posts
- Sistema de slugs únicos

#### `post_tags`

- Tabla de relación many-to-many entre posts y tags

#### `comments`

- Sistema de comentarios en posts
- Soporte para comentarios anidados (replies)
- Sistema de aprobación de comentarios

#### `user_sessions`

- Gestión de sesiones de usuario
- Tokens hasheados para seguridad

#### `activity_logs`

- Logs de actividad del sistema
- Almacenamiento en formato JSONB para flexibilidad

### Índices y Optimización

La base de datos incluye índices optimizados para:

- Búsquedas por email y username
- Filtros por estado de posts
- Ordenamiento por fechas
- Búsquedas por slugs
- Consultas de comentarios

### Triggers y Funciones

#### `update_updated_at_column()`

- Actualiza automáticamente el campo `updated_at`
- Se ejecuta en todas las tablas relevantes

#### `generate_unique_slug()`

- Genera slugs únicos para posts y categorías
- Evita conflictos de nombres duplicados

## 🌱 Seeder y Datos de Prueba

### Seeder Automático

El proyecto incluye un seeder completo que puede poblar la base de datos con diferentes volúmenes de datos:

#### 🔥 Seeder Masivo (60,000+ registros)

```bash
# Ejecutar seeder masivo
make seed-massive

# O usar el script directamente
./scripts/seed-db.sh --massive
```

**Genera:**

- **1,000 usuarios** con nombres realistas (Alex, Jordan, Taylor, etc.)
- **15 categorías** (Tecnología, Ciencia, Salud, Educación, etc.)
- **100+ tags** tecnológicos y generales
- **5,000 posts** con contenido variado
- **15,000 comentarios** de diferentes usuarios
- **25,000 relaciones** post-tag
- **10,000 logs** de actividad del sistema

#### 📝 Seeder Básico (30 registros)

```bash
# Ejecutar seeder básico
make seed-small

# O usar el script directamente
./scripts/seed-db.sh --small
```

**Genera:**

- **5 usuarios** básicos (admin, johndoe, janesmith, etc.)
- **8 categorías** principales
- **10 tags** esenciales
- **4 posts** de ejemplo
- **Comentarios** básicos

#### ⚡ Seeder por Defecto

```bash
# Ejecutar seeder por defecto
make seed
```

### Características del Seeder

- **Inserción en lotes**: Optimizado para insertar miles de registros eficientemente
- **Datos realistas**: Nombres, emails y contenido que simulan un entorno real
- **Relaciones coherentes**: Mantiene integridad referencial entre entidades
- **Logging detallado**: Progreso en tiempo real con emojis y estadísticas
- **Manejo de errores**: Recuperación robusta y mensajes informativos

### Comandos Disponibles

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

### Datos Generados

#### Usuarios

- Nombres y apellidos realistas de diferentes orígenes
- Emails únicos con dominios populares (gmail.com, yahoo.com, etc.)
- Contraseñas hasheadas (todas usan "password123" por defecto)
- Perfiles activos con timestamps realistas

#### Categorías

- 15 categorías principales que cubren diferentes temas
- Descripciones detalladas y slugs únicos
- Estado activo por defecto

#### Tags

- **Tags tecnológicos**: Go, Python, JavaScript, Docker, Kubernetes, AWS, etc.
- **Tags generales**: Salud, Fitness, Educación, Negocios, Arte, etc.
- Slugs automáticos y descripciones contextuales

#### Posts

- Títulos variados sobre temas tecnológicos
- Contenido de ejemplo extenso
- Autores asignados aleatoriamente
- Categorías distribuidas uniformemente
- Fechas de publicación del último año

#### Comentarios

- 15 plantillas de comentarios realistas
- Distribución aleatoria entre posts y usuarios
- 90% de comentarios aprobados por defecto

#### Relaciones Post-Tag

- Asignación aleatoria pero coherente
- Evita duplicados con `ON CONFLICT DO NOTHING`
- Distribución equilibrada entre entidades

#### Logs de Actividad

- 20 tipos de acciones diferentes
- 8 tipos de recursos
- IPs y user agents simulados
- Timestamps distribuidos

### Rendimiento

El seeder está optimizado para:

- **Inserción en lotes**: Reduce el número de consultas a la base de datos
- **Transacciones eficientes**: Mejor rendimiento en grandes volúmenes
- **Logging inteligente**: Progreso visible sin saturar la consola
- **Manejo de memoria**: Procesa datos en chunks manejables

### Monitoreo

Durante la ejecución del seeder masivo, verás:

- Progreso en tiempo real con emojis
- Contadores de registros insertados
- Logs de cada entidad procesada
- Estadísticas finales completas

### Troubleshooting

#### Error de Memoria

Si el seeder masivo consume mucha memoria:

```bash
# Reducir tamaño de lotes en el código
batchSize := 50  # En lugar de 100 o 500
```

#### Error de Timeout

Si hay timeouts en la base de datos:

```bash
# Aumentar timeouts en docker-compose.yml
command: ["postgres", "-c", "statement_timeout=300000"]
```

#### Limpieza Manual

Si necesitas limpiar manualmente:

```sql
-- Limpiar todas las tablas
TRUNCATE post_tags, comments, posts, tags, user_profiles, users, categories, activity_logs RESTART IDENTITY CASCADE;
```

### Usuarios de Prueba

| Usuario    | Email                   | Contraseña  | Rol            |
| ---------- | ----------------------- | ----------- | -------------- |
| admin      | admin@goasync.com       | password123 | Administrador  |
| johndoe    | john.doe@example.com    | password123 | Desarrollador  |
| janesmith  | jane.smith@example.com  | password123 | Arquitecta     |
| bobwilson  | bob.wilson@example.com  | password123 | DevOps         |
| alicebrown | alice.brown@example.com | password123 | Data Scientist |

## 🛠️ Comandos Útiles

### Makefile

```bash
# Gestión de base de datos
make db-up          # Levantar servicios de BD
make db-down        # Detener servicios
make db-reset       # Reiniciar BD (eliminar volúmenes)
make db-logs        # Ver logs de PostgreSQL
make db-connect     # Conectar a PostgreSQL

# Seeder
make seed           # Ejecutar seeder local
make seed-docker    # Ejecutar seeder en Docker
make seed-clean     # Limpiar y ejecutar seeder

# Desarrollo completo
make dev-full       # Levantar BD + seeder + aplicación
make dev-docker     # Todo en Docker
```

### Scripts

```bash
# Seeder con opciones
./scripts/seed-db.sh --help
./scripts/seed-db.sh --clean --verbose
./scripts/seed-db.sh --docker
```

## 🔍 Consultas Útiles

### Estadísticas de la Base de Datos

```sql
SELECT * FROM get_database_stats();
```

### Posts Activos con Información

```sql
SELECT * FROM active_posts;
```

### Estadísticas de Usuarios

```sql
SELECT * FROM user_stats;
```

### Posts por Categoría

```sql
SELECT
    c.name as category,
    COUNT(p.id) as post_count
FROM categories c
LEFT JOIN posts p ON c.id = p.category_id
WHERE c.is_active = true
GROUP BY c.id, c.name
ORDER BY post_count DESC;
```

### Usuarios Más Activos

```sql
SELECT
    u.username,
    COUNT(p.id) as posts_written,
    COUNT(c.id) as comments_made
FROM users u
LEFT JOIN posts p ON u.id = p.author_id
LEFT JOIN comments c ON u.id = c.author_id
GROUP BY u.id, u.username
ORDER BY (posts_written + comments_made) DESC;
```

## 🐳 Docker

### Servicios Incluidos

- **PostgreSQL**: Base de datos principal
- **Redis**: Cache y sesiones (para futuras implementaciones)
- **pgAdmin**: Interfaz web para administrar PostgreSQL

### Acceso a pgAdmin

- **URL**: http://localhost:5050
- **Email**: admin@goasync.com
- **Contraseña**: admin123

### Configuración de Red

- **Red**: `goasync_goasync-network`
- **Subnet**: `172.20.0.0/16`
- **Puertos expuestos**: 5432, 6379, 5050

## 🔧 Variables de Entorno

```bash
# Base de datos
DB_HOST=localhost          # o 'postgres' para Docker
DB_PORT=5432
DB_NAME=goasync
DB_USER=postgres
DB_PASSWORD=password
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost       # o 'redis' para Docker
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

## 📈 Monitoreo y Mantenimiento

### Health Checks

Los contenedores incluyen health checks automáticos:

- **PostgreSQL**: `pg_isready` cada 10 segundos
- **Redis**: `redis-cli ping` cada 10 segundos

### Logs

```bash
# Ver logs de PostgreSQL
make db-logs

# Ver logs de Redis
docker-compose logs -f redis

# Ver logs de pgAdmin
docker-compose logs -f pgadmin
```

### Backup y Restore

```bash
# Backup
docker-compose exec postgres pg_dump -U postgres goasync > backup.sql

# Restore
docker-compose exec -T postgres psql -U postgres goasync < backup.sql
```

## 🚨 Solución de Problemas

### Problemas Comunes

#### 1. Puerto 5432 ya en uso

```bash
# Ver qué está usando el puerto
lsof -i :5432

# Detener servicio local de PostgreSQL
sudo service postgresql stop
```

#### 2. Error de conexión en Docker

```bash
# Verificar que los contenedores estén ejecutándose
docker-compose ps

# Reiniciar servicios
docker-compose restart postgres
```

#### 3. Error de permisos en volúmenes

```bash
# Limpiar volúmenes y reiniciar
make db-reset
```

### Verificación de Estado

```bash
# Verificar estado de contenedores
docker-compose ps

# Verificar logs
docker-compose logs postgres

# Verificar conectividad
make db-connect
```

## 🔮 Futuras Mejoras

- [ ] Migraciones automáticas con herramientas como `golang-migrate`
- [ ] Replicación maestro-esclavo para alta disponibilidad
- [ ] Backup automático programado
- [ ] Monitoreo con Prometheus y Grafana
- [ ] Particionamiento de tablas grandes
- [ ] Compresión de datos históricos
- [ ] Auditoría completa de cambios

## 📚 Recursos Adicionales

- [Documentación oficial de PostgreSQL](https://www.postgresql.org/docs/)
- [Docker Hub PostgreSQL](https://hub.docker.com/_/postgres)
- [pgAdmin Documentation](https://www.pgadmin.org/docs/)
- [PostgreSQL Performance Tuning](https://www.postgresql.org/docs/current/performance.html)
