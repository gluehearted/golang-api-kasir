package main

import (
	"category-api-ss2/database"
	"category-api-ss2/handlers"
	"category-api-ss2/repositories"
	"category-api-ss2/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Category represents a product category in the cashier system
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// In-memory storage for categories (sementara, nanti ganti database)
var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Kategori untuk makanan"},
	{ID: 2, Name: "Minuman", Description: "Kategori untuk minuman"},
}

// GET /api/categories/{id}
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID dari URL path
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// Cari category dengan ID tersebut
	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	// Kalau tidak found
	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

// PUT /api/categories/{id}
func updateCategory(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// get data dari request
	var updateCat Category
	err = json.NewDecoder(r.Body).Decode(&updateCat)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// loop categories, cari id, ganti sesuai data dari request
	for i := range categories {
		if categories[i].ID == id {
			updateCat.ID = id
			categories[i] = updateCat

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCat)
			return
		}
	}

	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

// DELETE /api/categories/{id}
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// ganti id int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// loop categories cari ID, dapet index yang mau dihapus
	for i, c := range categories {
		if c.ID == id {
			// bikin slice baru dengan data sebelum dan sesudah index
			categories = append(categories[:i], categories[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete category",
			})
			return
		}
	}

	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Dependency Injection
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Setup routes for Product
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	// Category endpoints
	// GET localhost:8080/api/categories/{id}
	// PUT localhost:8080/api/categories/{id}
	// DELETE localhost:8080/api/categories/{id}
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}
	})

	// GET localhost:8080/api/categories
	// POST localhost:8080/api/categories
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
		} else if r.Method == "POST" {
			// baca data dari request
			var catBaru Category
			err := json.NewDecoder(r.Body).Decode(&catBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			// masukkin data ke dalam variable categories
			catBaru.ID = len(categories) + 1
			categories = append(categories, catBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(catBaru)
		}
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	addr := "localhost:" + config.Port
	fmt.Println("Server running di", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("gagal running server", err)
	}
}
