package main

import (
	"context"
	"fmt"
	"go-sandbox/api/internal/config"
	"go-sandbox/api/internal/database"
	"log"
	"math/rand"
	"strings"
	"time"
)

var firstNames = []string{
	"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda",
	"David", "Elizabeth", "William", "Barbara", "Richard", "Susan", "Joseph", "Jessica",
	"Thomas", "Sarah", "Charles", "Karen", "Daniel", "Lisa", "Matthew", "Nancy",
	"Anthony", "Betty", "Mark", "Margaret", "Steven", "Ashley", "Paul", "Emily",
	"Andrew", "Donna", "Joshua", "Michelle", "Kenneth", "Carol", "Kevin", "Amanda",
	"Brian", "Melissa", "George", "Deborah", "Timothy", "Stephanie", "Ronald", "Rebecca",
	"Jason", "Sharon", "Jeffrey", "Laura", "Ryan", "Cynthia", "Jacob", "Kathleen",
	"Gary", "Amy", "Nicholas", "Angela", "Eric", "Shirley", "Jonathan", "Anna",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
	"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson",
	"Thomas", "Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson",
	"White", "Harris", "Clark", "Lewis", "Robinson", "Walker", "Young", "Allen",
	"King", "Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores", "Green",
	"Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell", "Mitchell", "Carter",
	"Roberts", "Chen", "Kim", "Patel", "Singh", "Kumar", "Ali", "Sato",
	"Müller", "Schmidt", "Weber", "Fischer", "Meyer", "Wagner", "Becker", "Schulz",
}

var countries = []string{
	"US", "GB", "DE", "FR", "CA", "AU", "JP", "BR", "IN", "NL",
	"SE", "NO", "DK", "FI", "ES", "IT", "PT", "PL", "UA", "KR",
}

var roles = []string{"user", "user", "user", "user", "user", "user", "user", "admin", "moderator", "moderator"}
var statuses = []string{"active", "active", "active", "active", "active", "active", "active", "inactive", "inactive", "banned"}

var postTitles = []string{
	"Getting Started with Go", "Understanding Goroutines", "Building REST APIs",
	"PostgreSQL Tips and Tricks", "Docker Best Practices", "Clean Code Principles",
	"Microservices Architecture", "Database Optimization", "Error Handling in Go",
	"Testing Strategies", "CI/CD Pipeline Setup", "Kubernetes Basics",
	"React Performance Tips", "TypeScript Advanced Patterns", "GraphQL vs REST",
	"Caching Strategies", "Message Queues Explained", "API Rate Limiting",
	"Authentication Best Practices", "Monitoring and Observability",
}

var postBodies = []string{
	"In this post, we explore the fundamentals and best practices that every developer should know.",
	"After working on several production systems, I've collected these insights for other developers.",
	"One of the most common questions I get asked is how to approach this topic properly.",
	"This is a comprehensive guide covering everything from basics to advanced topics.",
	"Performance is crucial in modern applications. Here are practical techniques you can apply today.",
}

var postStatuses = []string{"published", "published", "published", "draft"}

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Clear existing data
	if _, err := db.Exec(ctx, "TRUNCATE users, posts CASCADE"); err != nil {
		log.Fatalf("Failed to truncate: %v", err)
	}

	var userValues []string
	var userArgs []any
	argIdx := 1
	startTime := time.Now().AddDate(-1, 0, 0) // Start 1 year ago

	for i := 0; i < 500; i++ {
		createdAt := startTime.Add(time.Duration(i) * 17 * time.Hour) // ~17 hours apart
		updatedAt := createdAt
		first := firstNames[rand.Intn(len(firstNames))]
		last := lastNames[rand.Intn(len(lastNames))]
		email := fmt.Sprintf("%s.%s%d@example.com", strings.ToLower(first), strings.ToLower(last), i)
		country := countries[rand.Intn(len(countries))]
		role := roles[rand.Intn(len(roles))]
		status := statuses[rand.Intn(len(statuses))]
		avatar := fmt.Sprintf("https://i.pravatar.cc/150?u=%d", i)

		userValues = append(userValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			argIdx, argIdx+1, argIdx+2, argIdx+3, argIdx+4, argIdx+5, argIdx+6, argIdx+7, argIdx+8))
		userArgs = append(userArgs, first, last, email, role, status, country, avatar, createdAt, updatedAt)
		argIdx += 9
	}

	userQuery := "INSERT INTO users (first_name, last_name, email, role, status, country, avatar_url, created_at, updated_at) VALUES " +
		strings.Join(userValues, ", ") + " RETURNING id"

	rows, err := db.Query(ctx, userQuery, userArgs...)
	if err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	var userIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Fatalf("Failed to scan user id: %v", err)
		}
		userIDs = append(userIDs, id)
	}
	rows.Close()
	log.Println("Seeding 500 users...")

	log.Printf("✅ Seeded %d users", len(userIDs))

	// Seed 2000 posts
	log.Println("Seeding 2000 posts...")
	var postValues []string
	var postArgs []any
	argIdx = 1
	startTime = time.Now().AddDate(-1, 0, 0)

	for i := 0; i < 2000; i++ {
		createdAt := startTime.Add(time.Duration(i) * 4 * time.Hour) // ~4 hours apart
		updatedAt := createdAt
		userID := userIDs[rand.Intn(len(userIDs))]
		title := postTitles[rand.Intn(len(postTitles))]
		body := postBodies[rand.Intn(len(postBodies))]
		status := postStatuses[rand.Intn(len(postStatuses))]

		postValues = append(postValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			argIdx, argIdx+1, argIdx+2, argIdx+3, argIdx+4, argIdx+5))
		postArgs = append(postArgs, userID, title, body, status, createdAt, updatedAt)
		argIdx += 6
	}

	postQuery := "INSERT INTO posts (user_id, title, content, status, created_at, updated_at) VALUES " +
		strings.Join(postValues, ", ") + " RETURNING id"
	_, err = db.Exec(ctx, postQuery, postArgs...)

	if err != nil {
		log.Fatalf("Failed to seed posts: %v", err)
	}
	log.Println("✅ Seeded 2000 posts")
	log.Println("🎉 Seeding complete!")

}
