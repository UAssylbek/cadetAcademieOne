package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"
)

// User представляет собой структуру пользователя
type User struct {
	Authenticated bool
	Email         string `json:"email"`
	Username      string `json:"username"`
	Password      string `json:"password"`
}

// Post представляет собой структуру поста
type Post struct {
	PostID    int
	UserID    int
	Title     string
	Content   string
	CreatedAt string
}

// Comment представляет собой структуру комментария
type Comment struct {
	CommentID int
	PostID    int
	UserID    int
	Content   string
	CreatedAt string
}

var user User
var db *sql.DB
var sessions map[string]time.Time

func main() {
	// Устанавливаем соединение с базой данных SQLite
	var err error
	db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверяем соединение с базой данных
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	// Создаем таблицы, если они еще не существуют
	createTables(db)

	// Инициализируем карту sessions
	sessions = make(map[string]time.Time)

	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/createpost", createPostHandler)
	http.HandleFunc("/createcomment", createCommentHandler)
	http.HandleFunc("/post", postHandler) // Добавляем обработчик для страницы с постом
	log.Println("serve start listening on port 8080")
	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "dashboard.html", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//// Парсим параметр запроса "id" для получения идентификатора поста
	//postID := r.URL.Query().Get("id")
	//if postID == "" {
	//	http.Error(w, "Bad Request", http.StatusBadRequest)
	//	return
	//}
	//
	//// Здесь нужно добавить логику для получения отдельного поста из базы данных
	//post, err := getPost(db, postID)
	//if err != nil {
	//	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	//	return
	//}

	// Отображаем HTML-шаблон отдельного поста
	renderTemplate(w, "home.html", user)
}

// getPost получает отдельный пост из базы данных по его идентификатору
func getPost(db *sql.DB, postID string) (Post, error) {
	var post Post
	err := db.QueryRow("SELECT * FROM Posts WHERE PostID = ?", postID).Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем данные из тела запроса
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Извлекаем данные из формы
	email := r.Form.Get("email")
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	// Проверяем, что все поля заполнены
	if email == "" || username == "" || password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Проверяем, что email не используется
	// Здесь нужно добавить логику проверки уникальности email в базе данных

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Server Error1", http.StatusInternalServerError)
		return
	}

	// Создаем нового пользователя
	user = User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}
	// Здесь нужно добавить логику добавления пользователя в базу данных
	query := "INSERT INTO Users (Email, Username, PasswordHash) VALUES (?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error2", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Email, user.Username, user.Password)
	if err != nil {
		http.Error(w, "Internal Server Error3", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ об успешной регистрации
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем данные из тела запроса
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Извлекаем данные из формы
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Запрос на выборку пользователя по email из базы данных
	row := db.QueryRow("SELECT UserID, PasswordHash, Username FROM Users WHERE Email = ?", email)

	// Переменные для хранения результата запроса
	var userID int
	var passwordHash, username string

	// Сканируем данные из результата запроса в переменные
	err = row.Scan(&userID, &passwordHash, &username)
	if err != nil {
		// Пользователь с таким email не найден
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
	}

	// Проверяем соответствие пароля хэшу из базы данных
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		// Пароль не совпадает
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	sessionID := generateSessionID()
	expiration := time.Now().Add(24 * time.Hour)
	// Если все проверки пройдены успешно, создаем сессию и отправляем ответ об успешном входе
	// В реальном приложении используйте пакет для управления сессиями, такой как sessions или gorilla/sessions
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	user.Authenticated = true
	user.Username = username
	sessions[sessionID] = expiration
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := filepath.Join("web", tmpl) // Изменяем путь к файлу на корректный
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем данные из тела запроса
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Извлекаем данные из формы
	title := r.Form.Get("title")
	content := r.Form.Get("content")

	// Проверяем, что оба поля заполнены
	if title == "" || content == "" {
		http.Error(w, "Title and Content are required", http.StatusBadRequest)
		return
	}

	// Создаем новый пост
	post := Post{
		UserID:  1, // Здесь нужно использовать идентификатор текущего пользователя
		Title:   title,
		Content: content,
	}
	err = createPost(db, post)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ об успешном создании поста
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Post created successfully")
}

func createCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем данные из тела запроса
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Извлекаем данные из формы
	postID := r.Form.Get("post_id")
	content := r.Form.Get("content")

	// Проверяем, что оба поля заполнены
	if postID == "" || content == "" {
		http.Error(w, "Post ID and Content are required", http.StatusBadRequest)
		return
	}

	// Создаем новый комментарий
	comment := Comment{
		PostID:  1, // Здесь нужно использовать идентификатор текущего пользователя
		Content: content,
	}
	err = createComment(db, comment)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ об успешном создании комментария
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Comment created successfully")
}

func createComment(db *sql.DB, comment Comment) error {
	// Подготавливаем запрос на вставку данных
	query := "INSERT INTO Comments (PostID, Content) VALUES (?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполняем запрос на вставку данных
	_, err = stmt.Exec(comment.PostID, comment.Content)
	if err != nil {
		return err
	}

	return nil
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID поста из запроса
	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	// Запрос к базе данных для получения информации о посте и его комментариях
	post, comments, err := getPostAndComments(postID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Выводим страницу с информацией о посте и его комментариях
	renderTemplate(w, "post.html", struct {
		Post     Post
		Comments []Comment
	}{
		Post:     post,
		Comments: comments,
	})
}

// getPostAndComments извлекает информацию о посте и его комментариях из базы данных
func getPostAndComments(postID string) (Post, []Comment, error) {
	// Запрос к базе данных для получения информации о посте по его ID
	row := db.QueryRow("SELECT * FROM Posts WHERE PostID = ?", postID)
	var post Post
	err := row.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		return Post{}, nil, err
	}

	// Запрос к базе данных для получения комментариев к посту
	rows, err := db.Query("SELECT * FROM Comments WHERE PostID = ?", postID)
	if err != nil {
		return Post{}, nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return Post{}, nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return Post{}, nil, err
	}

	return post, comments, nil
}

// logout обрабатывает запросы на выход из системы
func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		delete(sessions, cookie.Value) // Удаляем сессию из карты sessions
		// Очистка куки сессии
		cookie := http.Cookie{
			Name:    "session_id",
			Value:   "",
			Expires: time.Unix(0, 0),
		}
		http.SetCookie(w, &cookie)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

//func logout(w http.ResponseWriter, r *http.Request) {
//	// Очистка куки сессии
//	cookie := http.Cookie{
//		Name:    "session_id",
//		Value:   "",
//		Expires: time.Unix(0, 0),
//	}
//	http.SetCookie(w, &cookie)
//	user.Username = ""
//	user.Authenticated = false
//	http.Redirect(w, r, "/", http.StatusSeeOther)
//}

func generateSessionID() string {
	rand.Seed(time.Now().UnixNano()) // Инициализация генератора случайных чисел
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(rand.Intn(26) + 65) // Генерация случайной буквы
	}
	return string(b)
}

// createTables создает таблицы в базе данных, если они еще не существуют
func createTables(db *sql.DB) {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS Users (
            UserID INTEGER PRIMARY KEY AUTOINCREMENT,
            Username TEXT NOT NULL UNIQUE,
            Email TEXT NOT NULL UNIQUE,
            PasswordHash TEXT NOT NULL
        );
        CREATE TABLE IF NOT EXISTS Posts (
            PostID INTEGER PRIMARY KEY AUTOINCREMENT,
            UserID INTEGER NOT NULL,
            Title TEXT NOT NULL,
            Content TEXT NOT NULL,
            CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (UserID) REFERENCES Users(UserID)
        );
        CREATE TABLE IF NOT EXISTS Comments (
            CommentID INTEGER PRIMARY KEY AUTOINCREMENT,
            PostID INTEGER NOT NULL,
            UserID INTEGER NOT NULL,
            Content TEXT NOT NULL,
            CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (PostID) REFERENCES Posts(PostID),
            FOREIGN KEY (UserID) REFERENCES Users(UserID)
        );
    `)
	if err != nil {
		log.Fatal(err)
	}
}

// createPost добавляет новый пост в базу данных
func createPost(db *sql.DB, post Post) error {
	// Подготавливаем запрос на вставку данных
	query := "INSERT INTO Posts (UserID, Title, Content) VALUES (?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполняем запрос на вставку данных
	_, err = stmt.Exec(post.UserID, post.Title, post.Content)
	if err != nil {
		return err
	}

	return nil
}
