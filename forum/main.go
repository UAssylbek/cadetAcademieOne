package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	ID            string
	Username      string
	Authenticated bool
}

type App struct {
	DB          *sql.DB
	SessionsMap map[string]Session
	Mutex       sync.Mutex
}

func main() {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	createTables(db)

	app := &App{
		DB:          db,
		SessionsMap: make(map[string]Session),
	}
	fs := http.FileServer(http.Dir("uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", fs))
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	http.HandleFunc("/", app.homeHandler)
	http.HandleFunc("/dashboard", app.dashboardHandler)
	http.HandleFunc("/register", app.registerHandler)
	http.HandleFunc("/login", app.loginHandler)
	http.HandleFunc("/logout", app.logoutHandler)
	http.HandleFunc("/post", app.postHandler)
	http.HandleFunc("/category", app.categoryHandler)
	http.HandleFunc("/create-post", app.createPostHandler)
	http.HandleFunc("/create-category", app.createCategoryHandler)
	log.Println("Server started listening on port 8080")
	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (app *App) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "dashboard.html", nil)
}

func (app *App) postHandler(w http.ResponseWriter, r *http.Request) {
	// Запрос всех категорий из базы данных
	categoryes, err := app.DB.Query("SELECT Name FROM Categories")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer categoryes.Close()

	// Создание структуры для хранения всех постов
	type Category struct {
		Name string `json:"name"`
	}
	var allCategory []Category

	// Сканирование результатов запроса в структуру
	for categoryes.Next() {
		var category Category
		if err := categoryes.Scan(&category.Name); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		allCategory = append(allCategory, category)
	}
	if err := categoryes.Err(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "post.html", allCategory)
}

func (app *App) categoryHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "category.html", nil)
}

func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Запрос всех категорий из базы данных
	postcategoryes, err := app.DB.Query("SELECT CategoryID FROM PostCategories")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer postcategoryes.Close()

	// Создание структуры для хранения всех постов
	type PostCategory struct {
		CategoryID string `json:"category"`
	}
	var allPostCategory []PostCategory

	// Сканирование результатов запроса в структуру
	for postcategoryes.Next() {
		var postcategory PostCategory
		if err := postcategoryes.Scan(&postcategory.CategoryID); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		allPostCategory = append(allPostCategory, postcategory)
	}
	if err := postcategoryes.Err(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Запрос всех категорий из базы данных
	categoryes, err := app.DB.Query("SELECT Name FROM Categories")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer categoryes.Close()

	// Создание структуры для хранения всех постов
	type Category struct {
		Name string `json:"name"`
	}
	var allCategory []Category

	// Сканирование результатов запроса в структуру
	for categoryes.Next() {
		var category Category
		if err := categoryes.Scan(&category.Name); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		allCategory = append(allCategory, category)
	}
	if err := categoryes.Err(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Запрос всех постов из базы данных
	posts, err := app.DB.Query("SELECT Title, Content, ImagePath, CreatedAt FROM Posts")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer posts.Close()

	// Создание структуры для хранения всех постов
	type Post struct {
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		ImagePath string    `json:"image_path,omitempty"` // omitempty to omit the field if it's empty
		CreatedAt time.Time `json:"created_at"`
	}
	var allPosts []Post

	// Сканирование результатов запроса в структуру
	for posts.Next() {
		var post Post
		if err := posts.Scan(&post.Title, &post.Content, &post.ImagePath, &post.CreatedAt); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		allPosts = append(allPosts, post)
	}
	if err := posts.Err(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	type HomePageData struct {
		Session      Session
		Posts        []Post
		Category     []Category
		PostCategory []PostCategory
	}
	var homeData HomePageData
	cookie, err := r.Cookie("sessionID")
	if err != nil || cookie.Value == "" {
		homeData = HomePageData{
			Posts:        allPosts,
			Category:     allCategory,
			PostCategory: allPostCategory,
		}
		renderTemplate(w, "home.html", homeData)
		return
	}

	app.Mutex.Lock()
	session, ok := app.SessionsMap[cookie.Value]
	app.Mutex.Unlock()
	if !ok || !session.Authenticated {
		homeData = HomePageData{
			Posts:        allPosts,
			Category:     allCategory,
			PostCategory: allPostCategory,
		}
		renderTemplate(w, "home.html", homeData)
		return
	}

	homeData = HomePageData{
		Session:      session,
		Posts:        allPosts,
		Category:     allCategory,
		PostCategory: allPostCategory,
	}

	// Передача всех постов в шаблон для отображения
	renderTemplate(w, "home.html", homeData)
}

func (app *App) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var userData User
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if userData.Email == "" || userData.Username == "" || userData.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Все поля должны быть заполнены"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO Users (Email, Username, PasswordHash) VALUES (?, ?, ?)"
	stmt, err := app.DB.Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userData.Email, userData.Username, string(hashedPassword))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			var errorMessage, field string
			if strings.Contains(err.Error(), "Users.Username") {
				errorMessage = "Этот логин уже используется"
				field = "username"
			} else if strings.Contains(err.Error(), "Users.Email") {
				errorMessage = "Этот email уже используется"
				field = "email"
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": errorMessage, "field": field})
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно зарегистрирован"})
}

func (app *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginData User
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	row := app.DB.QueryRow("SELECT UserID, PasswordHash, Username FROM Users WHERE Email = ?", loginData.Email)

	var userID int
	var passwordHash, username string

	err = row.Scan(&userID, &passwordHash, &username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(loginData.Password)) != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный логин или пароль", "field": "wrong"})
		return
	}

	sessionID, err := uuid.NewV4()
	if err != nil {
		log.Printf("Error generating UUID: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "sessionID",
		Value: sessionID.String(),
		Path:  "/",
	})

	app.Mutex.Lock()
	app.SessionsMap[sessionID.String()] = Session{ID: sessionID.String(), Username: username, Authenticated: true}
	app.Mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"success": "Login successful", "sessionID": sessionID.String()})
}

func (app *App) logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		renderTemplate(w, "home.html", nil)
		return
	}
	sessionID := cookie.Value

	app.Mutex.Lock()
	delete(app.SessionsMap, sessionID)
	app.Mutex.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := filepath.Join("web", tmpl)
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

func createTables(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS Users (
    UserID INTEGER PRIMARY KEY AUTOINCREMENT,
    Username TEXT NOT NULL UNIQUE,
    Email TEXT NOT NULL UNIQUE,
    PasswordHash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Categories (
    CategoryID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Posts (
    PostID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID INTEGER,
    Title TEXT NOT NULL,
    Content TEXT NOT NULL,
    ImagePath TEXT,  -- новое поле для хранения пути к изображению
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (UserID) REFERENCES Users(UserID)
);


CREATE TABLE IF NOT EXISTS PostCategories (
    PostID INTEGER,
    CategoryID INTEGER,
    FOREIGN KEY (PostID) REFERENCES Posts(PostID),
    FOREIGN KEY (CategoryID) REFERENCES Categories(CategoryID)
);

CREATE TABLE IF NOT EXISTS Comments (
    CommentID INTEGER PRIMARY KEY AUTOINCREMENT,
    PostID INTEGER,
    UserID INTEGER,
    Content TEXT NOT NULL,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (PostID) REFERENCES Posts(PostID),
    FOREIGN KEY (UserID) REFERENCES Users(UserID)
);

CREATE TABLE IF NOT EXISTS Likes (
    LikeID INTEGER PRIMARY KEY AUTOINCREMENT,
    PostID INTEGER,
    CommentID INTEGER,
    UserID INTEGER,
    IsLike BOOLEAN,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (PostID) REFERENCES Posts(PostID),
    FOREIGN KEY (CommentID) REFERENCES Comments(CommentID),
    FOREIGN KEY (UserID) REFERENCES Users(UserID)
);

`)
	if err != nil {
		log.Fatal(err)
	}
}

func (app *App) createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("sessionID")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	app.Mutex.Lock()
	session, ok := app.SessionsMap[cookie.Value]
	app.Mutex.Unlock()
	if !ok || !session.Authenticated {
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	// Parse the multipart form, with a max upload size of 10MB
	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Retrieve form values
	title := r.FormValue("title")
	content := r.FormValue("content")
	categories := r.Form["categories"]
	fmt.Println(categories)

	// Handle image upload
	file, handler, err := r.FormFile("image")
	var imagePath string
	if err == nil {
		defer file.Close()
		// Ensure the "uploads" directory exists
		if _, err := os.Stat("uploads"); os.IsNotExist(err) {
			os.Mkdir("uploads", 0755)
		}

		imagePath = filepath.Join("uploads", fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename))
		fixedPath := strings.ReplaceAll(imagePath, "\\", "/")
		imagePath = fixedPath
		dst, err := os.Create(imagePath)
		if err != nil {
			fmt.Println("1")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		// Copy the uploaded file to the destination
		if _, err := file.Seek(0, 0); err != nil {
			fmt.Println("2")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if _, err := dst.ReadFrom(file); err != nil {
			fmt.Println("3")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	tx, err := app.DB.Begin()
	if err != nil {
		fmt.Println("4")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	result, err := tx.Exec("INSERT INTO Posts (UserID, Title, Content, ImagePath) VALUES (?, ?, ?, ?)", session.ID, title, content, imagePath)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		fmt.Println("6")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, category := range categories {
		var categoryID int64
		err = tx.QueryRow("SELECT CategoryID FROM Categories WHERE Name = ?", category).Scan(&categoryID)
		if err == sql.ErrNoRows {
			result, err = tx.Exec("INSERT INTO Categories (Name) VALUES (?)", category)
			if err != nil {
				tx.Rollback()
				fmt.Println("7")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			categoryID, err = result.LastInsertId()
			if err != nil {
				tx.Rollback()
				fmt.Println("8")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else if err != nil {
			tx.Rollback()
			fmt.Println("9")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("INSERT INTO PostCategories (PostID, CategoryID) VALUES (?, ?)", postID, categoryID)
		if err != nil {
			tx.Rollback()
			fmt.Println("10")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully"})
}

func (app *App) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var categoryData struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&categoryData)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if categoryData.Name == "" {
		fmt.Println("2")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	_, err = app.DB.Exec("INSERT INTO Categories (Name) VALUES (?)", categoryData.Name)
	if err != nil {
		fmt.Println("3")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Category created successfully"})
}
