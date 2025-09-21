package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"lab1/model"
	"net/http"
	"time"

	w "github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	DB        *sqlx.DB
	JWTSecret []byte
	TokenTTL  time.Duration
}

func NewHandler(db *sqlx.DB, secret []byte, ttl time.Duration) *Handler {
	return &Handler{
		DB:        db,
		JWTSecret: secret,
		TokenTTL:  ttl,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var body struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var user model.User
	query := `SELECT id, login, password FROM users WHERE login=$1`
	if err := h.DB.Get(&user, query, body.Login); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	claims := w.MapClaims{
		"user_id": user.ID,
		"login":   user.Login,
		"exp":     time.Now().Add(h.TokenTTL).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := w.NewWithClaims(w.SigningMethodHS256, claims)
	ss, err := token.SignedString(h.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": ss, "expires_in": int(h.TokenTTL.Seconds())})
}

func (h *Handler) GetData(c *gin.Context) {
	var posts []model.Post
	query := `SELECT id, author, title, content, created_at FROM posts ORDER BY created_at DESC`
	if err := h.DB.Select(&posts, query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (h *Handler) CreatePost(c *gin.Context) {
	var body struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	loginVal, ok := c.Get("login")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}
	login := loginVal.(string)

	query := `INSERT INTO posts (author, title, content, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id, author, title, content, created_at`
	var post model.Post
	if err := h.DB.Get(&post, query, login, body.Title, body.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"post": post})
}

const temp = `
<!doctype html>
<html>
<head><meta charset="utf-8"><title>Posts</title></head>
<body>
<h1>Posts</h1>
{{range .}}
  <div><b>{{.Title}}</b> — {{.Author}}<br/>{{.Content}}</div>
{{else}}
  <p>No posts</p>
{{end}}
</body>
</html>
`

func (h *Handler) ShowPosts(c *gin.Context) {
	c.Header("Content-Security-Policy", "default-src 'self'; script-src 'none'; object-src 'none'; frame-ancestors 'none';")

	var posts []model.Post
	query := `SELECT id, author, title, content, created_at FROM posts ORDER BY created_at DESC`
	if err := h.DB.Select(&posts, query); err != nil {
		c.String(http.StatusInternalServerError, "db error")
		return
	}

	tmpl := template.Must(template.New("posts").Parse(temp))
	_ = tmpl.Execute(c.Writer, posts)
}
