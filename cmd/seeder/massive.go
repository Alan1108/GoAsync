package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

// runMassiveSeeder ejecuta un seeder masivo con miles de registros
func runMassiveSeeder(db *sql.DB) error {
	// Limpiar datos existentes
	if err := cleanDatabase(db); err != nil {
		return fmt.Errorf("error limpiando base de datos: %w", err)
	}

	// Insertar categorías masivas
	if err := seedMassiveCategories(db); err != nil {
		return fmt.Errorf("error insertando categorías masivas: %w", err)
	}

	// Insertar usuarios masivos
	if err := seedMassiveUsers(db); err != nil {
		return fmt.Errorf("error insertando usuarios masivos: %w", err)
	}

	// Insertar tags masivos
	if err := seedMassiveTags(db); err != nil {
		return fmt.Errorf("error insertando tags masivos: %w", err)
	}

	// Insertar posts masivos
	if err := seedMassivePosts(db); err != nil {
		return fmt.Errorf("error insertando posts masivos: %w", err)
	}

	// Insertar comentarios masivos
	if err := seedMassiveComments(db); err != nil {
		return fmt.Errorf("error insertando comentarios masivos: %w", err)
	}

	// Insertar relaciones post-tag masivas
	if err := seedMassivePostTags(db); err != nil {
		return fmt.Errorf("error insertando relaciones post-tag masivas: %w", err)
	}

	// Insertar logs de actividad masivos
	if err := seedMassiveActivityLogs(db); err != nil {
		return fmt.Errorf("error insertando logs de actividad masivos: %w", err)
	}

	return nil
}

func seedMassiveCategories(db *sql.DB) error {
	log.Println("📂 Insertando categorías masivas...")

	categories := []Category{
		{Name: "Tecnología", Description: "Artículos sobre tecnología, programación y desarrollo", Slug: "tecnologia"},
		{Name: "Ciencia", Description: "Artículos sobre ciencia, investigación y descubrimientos", Slug: "ciencia"},
		{Name: "Salud", Description: "Artículos sobre salud, bienestar y medicina", Slug: "salud"},
		{Name: "Educación", Description: "Artículos sobre educación, aprendizaje y desarrollo personal", Slug: "educacion"},
		{Name: "Entretenimiento", Description: "Artículos sobre entretenimiento, cultura y ocio", Slug: "entretenimiento"},
		{Name: "Deportes", Description: "Artículos sobre deportes, fitness y actividades físicas", Slug: "deportes"},
		{Name: "Negocios", Description: "Artículos sobre negocios, emprendimiento y economía", Slug: "negocios"},
		{Name: "Viajes", Description: "Artículos sobre viajes, turismo y aventuras", Slug: "viajes"},
		{Name: "Cocina", Description: "Recetas, técnicas culinarias y gastronomía", Slug: "cocina"},
		{Name: "Arte", Description: "Arte, diseño, fotografía y creatividad", Slug: "arte"},
		{Name: "Música", Description: "Música, instrumentos y teoría musical", Slug: "musica"},
		{Name: "Literatura", Description: "Libros, escritura y análisis literario", Slug: "literatura"},
		{Name: "Historia", Description: "Historia, arqueología y eventos históricos", Slug: "historia"},
		{Name: "Filosofía", Description: "Filosofía, ética y pensamiento crítico", Slug: "filosofia"},
		{Name: "Psicología", Description: "Psicología, comportamiento humano y bienestar mental", Slug: "psicologia"},
	}

	for _, cat := range categories {
		query := `
			INSERT INTO categories (name, description, slug, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`

		var id string
		err := db.QueryRow(query, cat.Name, cat.Description, cat.Slug, true, time.Now(), time.Now()).Scan(&id)
		if err != nil {
			return fmt.Errorf("error insertando categoría %s: %w", cat.Name, err)
		}

		cat.ID = id
		log.Printf("  - Categoría '%s' insertada con ID: %s", cat.Name, id)
	}

	return nil
}

func seedMassiveUsers(db *sql.DB) error {
	log.Println("👥 Insertando usuarios masivos...")

	// Nombres y apellidos para generar usuarios realistas
	firstNames := []string{
		"Alex", "Jordan", "Taylor", "Casey", "Morgan", "Riley", "Quinn", "Avery", "Blake", "Cameron",
		"Jamie", "Drew", "Emery", "Finley", "Harper", "Kendall", "Logan", "Parker", "Reese", "Sage",
		"Skyler", "Tatum", "Wren", "Zion", "Adrian", "Bailey", "Charlie", "Dakota", "Eden", "Frankie",
		"Gray", "Hayden", "Indigo", "Jules", "Kai", "Lane", "Milan", "Nova", "Ocean", "Phoenix",
		"River", "Sage", "Teagan", "Winter", "Zen", "Aria", "Bella", "Chloe", "Diana", "Eva",
		"Fiona", "Grace", "Hannah", "Iris", "Jade", "Kate", "Luna", "Maya", "Nina", "Opal",
		"Paige", "Quinn", "Ruby", "Sofia", "Tara", "Uma", "Vera", "Willow", "Xena", "Yara", "Zara",
		"Adam", "Ben", "Carl", "Dan", "Eli", "Finn", "Gabe", "Hank", "Ian", "Jack",
		"Kyle", "Leo", "Max", "Nick", "Owen", "Paul", "Ryan", "Sam", "Tom", "Vince",
		"Wade", "Xander", "Yves", "Zane", "Aaron", "Brian", "Chris", "David", "Eric", "Frank",
		"George", "Henry", "Ivan", "James", "Kevin", "Liam", "Mark", "Noah", "Oscar", "Peter",
	}

	lastNames := []string{
		"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez",
		"Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin",
		"Lee", "Perez", "Thompson", "White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson",
		"Walker", "Young", "Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores",
		"Green", "Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell", "Mitchell", "Carter", "Roberts",
		"Gomez", "Kim", "Chen", "Wong", "Chang", "Singh", "Shah", "Patel", "Kumar", "Singh",
		"Kaur", "Sharma", "Verma", "Gupta", "Malhotra", "Kapoor", "Chopra", "Reddy", "Rao", "Iyer",
		"Menon", "Nair", "Krishnan", "Pillai", "Nambiar", "Kurup", "Kurian", "Mathew", "Philip", "Thomas",
		"Fernandez", "Silva", "Santos", "Costa", "Oliveira", "Pereira", "Carvalho", "Almeida", "Nascimento", "Lima",
		"Ribeiro", "Ferreira", "Alves", "Pinto", "Cunha", "Mendes", "Dias", "Castro", "Monteiro", "Moreira",
	}

	// Dominios de email populares
	emailDomains := []string{
		"gmail.com", "yahoo.com", "hotmail.com", "outlook.com", "icloud.com",
		"protonmail.com", "fastmail.com", "zoho.com", "mail.com", "aol.com",
		"live.com", "me.com", "mac.com", "msn.com", "rocketmail.com",
	}

	// Generar 1000 usuarios
	userCount := 1000
	log.Printf("  - Generando %d usuarios...", userCount)

	// Insertar usuarios en lotes para mejor rendimiento
	batchSize := 100
	for i := 0; i < userCount; i += batchSize {
		end := i + batchSize
		if end > userCount {
			end = userCount
		}

		// Preparar batch
		values := []interface{}{}
		placeholders := []string{}
		placeholderIndex := 1

		for j := i; j < end; j++ {
			firstName := firstNames[rand.Intn(len(firstNames))]
			lastName := lastNames[rand.Intn(len(lastNames))]
			username := fmt.Sprintf("%s%s%d", strings.ToLower(firstName), strings.ToLower(lastName), j)
			email := fmt.Sprintf("%s.%s%d@%s", strings.ToLower(firstName), strings.ToLower(lastName), j, emailDomains[rand.Intn(len(emailDomains))])

			values = append(values, username, email, "$2a$10$hashedpassword", firstName, lastName, true, time.Now(), time.Now())

			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				placeholderIndex, placeholderIndex+1, placeholderIndex+2, placeholderIndex+3,
				placeholderIndex+4, placeholderIndex+5, placeholderIndex+6, placeholderIndex+7))
			placeholderIndex += 8
		}

		query := fmt.Sprintf(`
			INSERT INTO users (username, email, password_hash, first_name, last_name, is_active, created_at, updated_at)
			VALUES %s
			RETURNING id`, strings.Join(placeholders, ", "))

		rows, err := db.Query(query, values...)
		if err != nil {
			return fmt.Errorf("error insertando batch de usuarios %d-%d: %w", i, end-1, err)
		}

		// Cerrar rows después de procesar
		rows.Close()

		if (i+batchSize)%500 == 0 || end == userCount {
			log.Printf("  - %d usuarios insertados", end)
		}
	}

	log.Printf("  - ✅ %d usuarios insertados exitosamente", userCount)
	return nil
}

func seedMassiveTags(db *sql.DB) error {
	log.Println("🏷️ Insertando tags masivos...")

	// Tags relacionados con tecnología y desarrollo
	techTags := []string{
		"Go", "Python", "JavaScript", "TypeScript", "Java", "C++", "C#", "Rust", "Kotlin", "Swift",
		"React", "Vue", "Angular", "Node.js", "Express", "Django", "Flask", "Spring", "Laravel", "Symfony",
		"PostgreSQL", "MySQL", "MongoDB", "Redis", "Elasticsearch", "Cassandra", "Neo4j", "InfluxDB", "TimescaleDB",
		"Docker", "Kubernetes", "Terraform", "Ansible", "Jenkins", "GitLab", "GitHub", "Bitbucket", "Jira", "Confluence",
		"AWS", "Azure", "GCP", "DigitalOcean", "Heroku", "Vercel", "Netlify", "Cloudflare", "Fastly", "Akamai",
		"Machine Learning", "AI", "Data Science", "Big Data", "Analytics", "Business Intelligence", "ETL", "Data Warehousing",
		"Microservices", "API", "REST", "GraphQL", "gRPC", "WebSocket", "Serverless", "Event-Driven", "CQRS", "Event Sourcing",
		"Security", "Authentication", "Authorization", "OAuth", "JWT", "OIDC", "Encryption", "HTTPS", "SSL", "TLS",
		"Testing", "Unit Testing", "Integration Testing", "E2E Testing", "TDD", "BDD", "Performance Testing", "Load Testing",
		"DevOps", "CI/CD", "Infrastructure as Code", "Monitoring", "Logging", "Tracing", "Alerting", "Metrics", "Observability",
	}

	// Tags de otras categorías
	otherTags := []string{
		"Salud", "Fitness", "Nutrición", "Medicina", "Psicología", "Bienestar", "Yoga", "Meditación", "Ejercicio",
		"Educación", "Aprendizaje", "Online", "Cursos", "Certificaciones", "Habilidades", "Desarrollo Personal", "Productividad",
		"Negocios", "Emprendimiento", "Marketing", "Ventas", "Finanzas", "Estrategia", "Innovación", "Liderazgo", "Management",
		"Arte", "Diseño", "Ilustración", "Arquitectura", "Moda", "Interiorismo", "Creatividad", "Inspiración",
		"Viajes", "Turismo", "Aventura", "Cultura", "Gastronomía", "Historia", "Arqueología", "Naturaleza",
		"Deportes", "Fútbol", "Baloncesto", "Tenis", "Golf", "Running", "Ciclismo", "Natación", "CrossFit",
		"Música", "Instrumentos", "Teoría", "Composición", "Producción", "DJ", "Bandas", "Conciertos", "Festivales",
		"Literatura", "Libros", "Escritura", "Poesía", "Novelas", "Cuentos", "Crítica", "Autores", "Editoriales",
	}

	allTags := append(techTags, otherTags...)
	tagCount := len(allTags)

	log.Printf("  - Generando %d tags...", tagCount)

	// Insertar tags en lotes
	batchSize := 50
	for i := 0; i < tagCount; i += batchSize {
		end := i + batchSize
		if end > tagCount {
			end = tagCount
		}

		values := []interface{}{}
		placeholders := []string{}
		placeholderIndex := 1

		for j := i; j < end; j++ {
			tagName := allTags[j]
			slug := strings.ToLower(strings.ReplaceAll(tagName, " ", "-"))
			description := fmt.Sprintf("Artículos y contenido relacionado con %s", tagName)

			values = append(values, tagName, slug, description, time.Now())

			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)",
				placeholderIndex, placeholderIndex+1, placeholderIndex+2, placeholderIndex+3))
			placeholderIndex += 4
		}

		query := fmt.Sprintf(`
			INSERT INTO tags (name, slug, description, created_at)
			VALUES %s
			RETURNING id`, strings.Join(placeholders, ", "))

		_, err := db.Exec(query, values...)
		if err != nil {
			return fmt.Errorf("error insertando batch de tags %d-%d: %w", i, end-1, err)
		}

		if (i+batchSize)%100 == 0 || end == tagCount {
			log.Printf("  - %d tags insertados", end)
		}
	}

	log.Printf("  - ✅ %d tags insertados exitosamente", tagCount)
	return nil
}

func seedMassivePosts(db *sql.DB) error {
	log.Println("📝 Insertando posts masivos...")

	// Obtener IDs de usuarios y categorías
	var userIDs []string
	var categoryIDs []string

	rows, err := db.Query("SELECT id FROM users LIMIT 100")
	if err != nil {
		return fmt.Errorf("error obteniendo IDs de usuarios: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("error escaneando ID de usuario: %w", err)
		}
		userIDs = append(userIDs, id)
	}

	rows, err = db.Query("SELECT id FROM categories")
	if err != nil {
		return fmt.Errorf("error obteniendo IDs de categorías: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("error escaneando ID de categoría: %w", err)
		}
		categoryIDs = append(categoryIDs, id)
	}

	// Títulos de posts de ejemplo
	postTitles := []string{
		"Introducción a Go: El lenguaje del futuro",
		"Construyendo APIs RESTful con Go y Gin",
		"Docker para desarrolladores: Una guía práctica",
		"PostgreSQL vs MySQL: ¿Cuál elegir para tu proyecto?",
		"Microservicios: Arquitectura y mejores prácticas",
		"Machine Learning para principiantes",
		"DevOps: Automatización y CI/CD",
		"Seguridad en aplicaciones web modernas",
		"Testing en Go: Estrategias y herramientas",
		"Optimización de bases de datos",
		"Cloud Computing: AWS, Azure y GCP",
		"GraphQL vs REST: Comparación completa",
		"Kubernetes: Orquestación de contenedores",
		"Event-Driven Architecture",
		"Domain-Driven Design en Go",
		"Clean Architecture: Principios y implementación",
		"Performance y escalabilidad",
		"Monitoring y observabilidad",
		"API Design: Mejores prácticas",
		"Database Sharding y particionamiento",
	}

	// Generar 5000 posts
	postCount := 5000
	log.Printf("  - Generando %d posts...", postCount)

	// Insertar posts en lotes
	batchSize := 100
	for i := 0; i < postCount; i += batchSize {
		end := i + batchSize
		if end > postCount {
			end = postCount
		}

		values := []interface{}{}
		placeholders := []string{}
		placeholderIndex := 1

		for j := i; j < end; j++ {
			title := postTitles[j%len(postTitles)]
			slug := fmt.Sprintf("%s-%d", strings.ToLower(strings.ReplaceAll(title, " ", "-")), j)
			content := fmt.Sprintf("Contenido del artículo %d: %s. Este es un artículo de ejemplo con contenido extenso sobre el tema.", j+1, title)
			excerpt := fmt.Sprintf("Resumen del artículo sobre %s", title)
			authorID := userIDs[rand.Intn(len(userIDs))]
			categoryID := categoryIDs[rand.Intn(len(categoryIDs))]
			status := "published"
			publishedAt := time.Now().AddDate(0, 0, -rand.Intn(365)) // Posts del último año

			values = append(values, title, slug, content, excerpt, authorID, categoryID, status, publishedAt, time.Now(), time.Now())

			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				placeholderIndex, placeholderIndex+1, placeholderIndex+2, placeholderIndex+3,
				placeholderIndex+4, placeholderIndex+5, placeholderIndex+6, placeholderIndex+7,
				placeholderIndex+8, placeholderIndex+9))
			placeholderIndex += 10
		}

		query := fmt.Sprintf(`
			INSERT INTO posts (title, slug, content, excerpt, author_id, category_id, status, published_at, created_at, updated_at)
			VALUES %s
			RETURNING id`, strings.Join(placeholders, ", "))

		_, err := db.Exec(query, values...)
		if err != nil {
			return fmt.Errorf("error insertando batch de posts %d-%d: %w", i, end-1, err)
		}

		if (i+batchSize)%1000 == 0 || end == postCount {
			log.Printf("  - %d posts insertados", end)
		}
	}

	log.Printf("  - ✅ %d posts insertados exitosamente", postCount)
	return nil
}

func seedMassiveComments(db *sql.DB) error {
	log.Println("💬 Insertando comentarios masivos...")

	// Obtener IDs de posts y usuarios
	var postIDs []string
	var userIDs []string

	rows, err := db.Query("SELECT id FROM posts LIMIT 1000")
	if err != nil {
		return fmt.Errorf("error obteniendo IDs de posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("error escaneando ID de post: %w", err)
		}
		postIDs = append(postIDs, id)
	}

	rows, err = db.Query("SELECT id FROM users LIMIT 200")
	if err != nil {
		return fmt.Errorf("error obteniendo IDs de usuarios: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("error escaneando ID de usuario: %w", err)
		}
		userIDs = append(userIDs, id)
	}

	// Comentarios de ejemplo
	commentTemplates := []string{
		"Excelente artículo! Muy útil para principiantes.",
		"Gracias por compartir esta información.",
		"Interesante perspectiva sobre el tema.",
		"¿Podrías profundizar más en este aspecto?",
		"Me gustó mucho la explicación.",
		"Buena información, pero creo que falta mencionar...",
		"Este enfoque me parece muy práctico.",
		"¿Tienes algún ejemplo adicional?",
		"Perfecto para mi proyecto actual.",
		"Me ayudó mucho a entender el concepto.",
		"Excelente trabajo, muy bien explicado.",
		"¿Podrías hacer un artículo sobre...?",
		"Me gustaría ver más contenido similar.",
		"Buena introducción al tema.",
		"¿Hay alguna alternativa a esta solución?",
	}

	// Generar 15000 comentarios
	commentCount := 15000
	log.Printf("  - Generando %d comentarios...", commentCount)

	// Insertar comentarios en lotes
	batchSize := 200
	for i := 0; i < commentCount; i += batchSize {
		end := i + batchSize
		if end > commentCount {
			end = commentCount
		}

		values := []interface{}{}
		placeholders := []string{}
		placeholderIndex := 1

		for j := i; j < end; j++ {
			postID := postIDs[rand.Intn(len(postIDs))]
			authorID := userIDs[rand.Intn(len(userIDs))]
			content := commentTemplates[rand.Intn(len(commentTemplates))]
			isApproved := rand.Float32() > 0.1 // 90% aprobados

			values = append(values, postID, authorID, content, isApproved, time.Now(), time.Now())

			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
				placeholderIndex, placeholderIndex+1, placeholderIndex+2,
				placeholderIndex+3, placeholderIndex+4, placeholderIndex+5))
			placeholderIndex += 6
		}

		query := fmt.Sprintf(`
			INSERT INTO comments (post_id, author_id, content, is_approved, created_at, updated_at)
			VALUES %s
			RETURNING id`, strings.Join(placeholders, ", "))

		_, err := db.Exec(query, values...)
		if err != nil {
			return fmt.Errorf("error insertando batch de comentarios %d-%d: %w", i, end-1, err)
		}

		if (i+batchSize)%2000 == 0 || end == commentCount {
			log.Printf("  - %d comentarios insertados", end)
		}
	}

	log.Printf("  - ✅ %d comentarios insertados exitosamente", commentCount)
	return nil
}

func seedMassivePostTags(db *sql.DB) error {
	log.Println("🔗 Insertando relaciones post-tag masivas...")

	// Obtener IDs de posts y tags
	var postIDs []string
	var tagIDs []string

	rows, err := db.Query("SELECT id FROM posts LIMIT 2000")
	if err != nil {
		return fmt.Errorf("error obteniendo IDs de posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("error escaneando ID de post: %w", err)
		}
		postIDs = append(postIDs, id)
	}

	rows, err = db.Query("SELECT id FROM tags")
	if err != nil {
		return fmt.Errorf("error obteniendo IDs de tags: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("error escaneando ID de tag: %w", err)
		}
		tagIDs = append(tagIDs, id)
	}

	// Generar 25000 relaciones post-tag
	relationCount := 25000
	log.Printf("  - Generando %d relaciones post-tag...", relationCount)

	// Insertar relaciones en lotes
	batchSize := 500
	for i := 0; i < relationCount; i += batchSize {
		end := i + batchSize
		if end > relationCount {
			end = relationCount
		}

		values := []interface{}{}
		placeholders := []string{}
		placeholderIndex := 1

		for j := i; j < end; j++ {
			postID := postIDs[rand.Intn(len(postIDs))]
			tagID := tagIDs[rand.Intn(len(tagIDs))]

			values = append(values, postID, tagID)

			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)",
				placeholderIndex, placeholderIndex+1))
			placeholderIndex += 2
		}

		query := fmt.Sprintf(`
			INSERT INTO post_tags (post_id, tag_id)
			VALUES %s
			ON CONFLICT DO NOTHING`, strings.Join(placeholders, ", "))

		_, err := db.Exec(query, values...)
		if err != nil {
			return fmt.Errorf("error insertando batch de relaciones %d-%d: %w", i, end-1, err)
		}

		if (i+batchSize)%5000 == 0 || end == relationCount {
			log.Printf("  - %d relaciones insertadas", end)
		}
	}

	log.Printf("  - ✅ %d relaciones post-tag insertadas exitosamente", relationCount)
	return nil
}

func seedMassiveActivityLogs(db *sql.DB) error {
	log.Println("📊 Insertando logs de actividad masivos...")

	// Obtener IDs de usuarios
	var userIDs []string
	rows, err := db.Query("SELECT id FROM users LIMIT 100")
	if err != nil {
		return fmt.Errorf("error obteniendo IDs de usuarios: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("error escaneando ID de usuario: %w", err)
		}
		userIDs = append(userIDs, id)
	}

	// Acciones de ejemplo
	actions := []string{
		"user_login", "user_logout", "post_created", "post_updated", "post_deleted",
		"comment_added", "comment_updated", "comment_deleted", "user_registered",
		"profile_updated", "password_changed", "email_verified", "account_locked",
		"search_performed", "file_uploaded", "file_downloaded", "api_call",
		"admin_action", "moderation_action", "backup_created", "system_maintenance",
	}

	// Tipos de recursos
	resourceTypes := []string{
		"user", "post", "comment", "category", "tag", "file", "session", "api",
	}

	// Generar 10000 logs de actividad
	logCount := 10000
	log.Printf("  - Generando %d logs de actividad...", logCount)

	// Insertar logs en lotes
	batchSize := 500
	for i := 0; i < logCount; i += batchSize {
		end := i + batchSize
		if end > logCount {
			end = logCount
		}

		values := []interface{}{}
		placeholders := []string{}
		placeholderIndex := 1

		for j := i; j < end; j++ {
			userID := userIDs[rand.Intn(len(userIDs))]
			action := actions[rand.Intn(len(actions))]
			resourceType := resourceTypes[rand.Intn(len(resourceTypes))]
			details := fmt.Sprintf(`{"ip": "192.168.1.%d", "user_agent": "Mozilla/5.0...", "timestamp": "%s"}`,
				rand.Intn(255), time.Now().Format(time.RFC3339))
			ipAddress := fmt.Sprintf("192.168.1.%d", rand.Intn(255))

			values = append(values, userID, action, resourceType, details, ipAddress, time.Now())

			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
				placeholderIndex, placeholderIndex+1, placeholderIndex+2,
				placeholderIndex+3, placeholderIndex+4, placeholderIndex+5))
			placeholderIndex += 6
		}

		query := fmt.Sprintf(`
			INSERT INTO activity_logs (user_id, action, resource_type, details, ip_address, created_at)
			VALUES %s
			RETURNING id`, strings.Join(placeholders, ", "))

		_, err := db.Exec(query, values...)
		if err != nil {
			return fmt.Errorf("error insertando batch de logs %d-%d: %w", i, end-1, err)
		}

		if (i+batchSize)%2000 == 0 || end == logCount {
			log.Printf("  - %d logs insertados", end)
		}
	}

	log.Printf("  - ✅ %d logs de actividad insertados exitosamente", logCount)
	return nil
}
