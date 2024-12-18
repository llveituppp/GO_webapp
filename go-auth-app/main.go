package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func InitDB() {
	var err error
	connStr := "postgres://postgres:123@localhost/go_auth_app?sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошбика подключения к базе данных:", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("Ошибка при проверке подключения к базе данных:", err)
	}
	log.Println("Подключение к базе данных выполнено успешно!")
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Отсутствует заголовок авторизации", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Недопустимый формат авторизации", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Недействительный или просроченный токен", http.StatusUnauthorized)
			return
		}

		// Передаём email в заголовках
		r.Header.Set("UserEmail", claims.Email)
		next.ServeHTTP(w, r)
	})
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("UserEmail")
	response := map[string]string{
		"message": "Welcome to the protected route!",
		"user":    email,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, "Ошибка хэширования пароля", http.StatusInternalServerError)
		return
	}

	_, err = DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, hash)
	if err != nil {
		http.Error(w, "Ошибка сохранения пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно зарегистрирован"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неправильный ввод", http.StatusBadRequest)
		return
	}

	var storedPassword string
	err := DB.QueryRow("SELECT password FROM users WHERE email=$1", user.Email).Scan(&storedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Неправильная почта или пароль", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Ошибка при запросе пользователя", http.StatusInternalServerError)
		return
	}

	if !checkPasswordHash(user.Password, storedPassword) {
		http.Error(w, "Неправильная почта или пароль", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, name, price, quantity FROM products")
	if err != nil {
		http.Error(w, "Ошибка при выборе товаров", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity); err != nil {
			http.Error(w, "Ошибка при проверке товаров", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func PurchaseProductHandler(w http.ResponseWriter, r *http.Request) {
	var purchase struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&purchase); err != nil {
		http.Error(w, "Неправильный ввод", http.StatusBadRequest)
		return
	}

	var currentQuantity int
	err := DB.QueryRow("SELECT quantity FROM products WHERE id=$1", purchase.ProductID).Scan(&currentQuantity)
	if err != nil {
		http.Error(w, "Товар не найден", http.StatusNotFound)
		return
	}

	if currentQuantity < purchase.Quantity {
		http.Error(w, "Недостаточно товаров", http.StatusBadRequest)
		return
	}

	_, err = DB.Exec("UPDATE products SET quantity = quantity - $1 WHERE id = $2", purchase.Quantity, purchase.ProductID)
	if err != nil {
		http.Error(w, "Ошибка обновления товаров", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Product purchased successfully"})
}

func main() {
	InitDB()
	defer DB.Close()

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")

	protected := r.PathPrefix("/protected").Subrouter()
	protected.Use(JWTMiddleware)
	protected.HandleFunc("/example", ProtectedHandler).Methods("GET")

	r.HandleFunc("/products", GetProductsHandler).Methods("GET")
	r.HandleFunc("/purchase", PurchaseProductHandler).Methods("POST")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	fmt.Println("Server running on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
