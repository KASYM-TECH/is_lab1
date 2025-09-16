package model

import "time"

type Post struct {
	ID        uint      `db:"id"`
	Author    string    `db:"author"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}
