package main

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"lab1/model"
)

func seedDemoSQLX(db *sqlx.DB) error {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE login=$1", "alice")
	if err != nil {
		return err
	}

	if count == 0 {
		pwHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if _, err := db.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", "alice", string(pwHash)); err != nil {
			return err
		}
	}

	// seed posts, если нет
	err = db.Get(&count, "SELECT COUNT(*) FROM posts")
	if err != nil {
		return err
	}

	if count == 0 {
		posts := []model.Post{
			{Author: "alice", Title: "Welcome", Content: "Hello, this is a safe post."},
			{Author: "bob", Title: "About XSS", Content: "Example content. <script>alert('xss')</script> — this will be escaped in HTML view."},
			{Author: "charlie", Title: "Markdown?", Content: "You can store markup, but templates escape it by default."},
		}
		for _, p := range posts {
			if _, err := db.Exec("INSERT INTO posts (author, title, content) VALUES ($1, $2, $3)", p.Author, p.Title, p.Content); err != nil {
				return err
			}
		}
	}

	return nil
}
