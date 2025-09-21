## Эндпоинты API

### Аутентификация

#### POST `/auth/login`
метод для аутентификации пользователя (принимает
логин и пароль).

**Тело запроса:**
```json
{
  "login": "testuser",
  "password": "securepassword123"
}
```

**Ответ:**
```json
{
  "token": "User created successfully",
  "expires_at": "timestamp"
}
```

### Защищенные эндпоинты (требуют JWT токен)

#### GET `/api/data`
Получение всех постов.

**Заголовок:**
```
Authorization: Bearer <your_jwt_token>
```

**Ответ:**
```json
{
  "posts": [
    {
      "ID": 5,
      "Author": "alice",
      "Title": "My New Post",
      "Content": "This is the post body.",
      "CreatedAt": "2025-09-16T13:20:07.385666Z"
    },
  ]
}
```

#### POST `/api/posts`
Создание нового поста.

**Заголовок:**
```
Authorization: Bearer <your_jwt_token>
```

**Тело запроса:**
```json
{
  "title": "My New Post",
  "content": "This is the post body."
}
```

**Ответ:**
```json
{
  "post": {
    "ID": 6,
    "Author": "alice",
    "Title": "My New Post",
    "Content": "This is the post body.",
    "CreatedAt": "2025-09-21T07:11:11.386164Z"
  }
}
```

#### GET `/api/posts`
Получение постов.

**Заголовок:**
```
Authorization: Bearer <your_jwt_token>
```

**Ответ:**
```json
{
  "post": {
    "ID": 6,
    "Author": "alice",
    "Title": "My New Post",
    "Content": "This is the post body.",
    "CreatedAt": "2025-09-21T07:11:11.386164Z"
  }
}
```

## Быстрый старт

### Установка и запуск

```bash
# Клонирование репозитория
git clone https://github.com/KASYM-TECH/is_lab1.git

docker compose up

go run .
```

### Примеры использования с http

```bash
### 1) Login — получить JWT
POST http://localhost:8080/auth/login
Content-Type: application/json

{
  "login": "alice",
  "password": "password123"
}


### 2) Get protected data — /api/data
GET http://localhost:8080/api/data
Authorization: Bearer <your_jwt_token>


### 3) Get posts HTML — /posts
GET http://localhost:8080/posts
Authorization: Bearer <your_jwt_token>


### 4) Create new post
POST http://localhost:8080/api/posts
Content-Type: application/json
Authorization: Bearer <your_jwt_token>

{
  "title": "My New Post",
  "content": "This is the post body."
}
```

## Меры защиты

### Защита от SQL-инъекций

**Реализация:** Использование параметризованных запросов через github.com/jmoiron/sqlx.

```go
	query := `INSERT INTO posts (author, title, content, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id, author, title, content, created_at`
	var post model.Post
	if err := h.DB.Get(&post, query, login, body.Title, body.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}
```

### Защита от XSS атак

**Реализация:** Санитизация всех постов пользователей.

```go
	var posts []model.Post
	query := `SELECT id, author, title, content, created_at FROM posts ORDER BY created_at DESC`
	if err := h.DB.Select(&posts, query); err != nil {
		c.String(http.StatusInternalServerError, "db error")
		return
	}

	tmpl := template.Must(template.New("posts").Parse(temp))
	_ = tmpl.Execute(c.Writer, posts)
```

### Защита аутентификации

**Хеширование паролей:** Использование "golang.org/x/crypto/bcrypt" 
```go
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
```

**JWT токены:** 
- Срок действия: 1 час
- Middleware проверки на всех защищенных эндпоинтах

```go
// jwtMiddleware verifies Bearer token and sets user info in context
func jwtMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}
		parts := strings.Fields(auth)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
			return
		}
		tokenString := parts[1]

		claims := &Claims{}
		token, err := w.ParseWithClaims(tokenString, claims, func(token *w.Token) (interface{}, error) {
			if _, ok := token.Method.(*w.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("login", claims.Login)
		c.Next()
	}
}
```

## CI/CD Pipeline

### Запускаемые проверки

1. **gosec** - Проверка уязвимостей в зависимостях
3. **OWASP Dependency Check** - Анализ зависимостей
4. **golangci-lint** - Статический анализ кода (linter)
5. **Go tests* - Unit тесты

## Отчеты безопасности

### npm audit report

<img width="449" height="132" alt="image" src="https://github.com/user-attachments/assets/a5fde753-cbf6-467c-9f07-c6146cda4ac0" />

*Отчет показывает 0 уязвимостей высокого уровня критичности*

### OWASP Dependency Check

<img width="1070" height="1205" alt="image" src="https://github.com/user-attachments/assets/4a926390-e106-41c7-abe5-5abbc2250e4d" />

**OWASP Dependency Check обнаружил уязвимость**

### ESLint Security Analysis

<img width="586" height="985" alt="image" src="https://github.com/user-attachments/assets/073c890c-ee30-422f-b5a1-0cccfb73a89d" />

*ESLint с security plugin не обнаружил проблем безопасности*

### Test Results

<img width="832" height="332" alt="image" src="https://github.com/user-attachments/assets/fc554eac-7ee3-4e2d-bab6-2b5e92edf011" />

*Все тесты проходят успешно, включая security-тесты*
