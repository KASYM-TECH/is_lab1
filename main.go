package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (using system environment)")
	}

	db, err := sqlx.Connect("postgres", dsn())
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	schema := `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    author TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
`
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	if err := seedDemoSQLX(db); err != nil {
		return fmt.Errorf("failed to seed demo data: %w", err)
	}

	h := NewHandler(db, []byte(os.Getenv("JWT_SECRET")), 6*time.Hour)

	r := gin.Default()

	r.POST("/auth/login", h.Login)
	r.GET("/posts", h.ShowPosts)

	api := r.Group("/api", jwtMiddleware())
	{
		api.GET("/data", h.GetData)
		api.POST("/posts", h.CreatePost)
	}

	return r.Run(":8080")
}
