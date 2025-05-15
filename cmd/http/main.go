package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	storage "project/internal/storage"
	postgres "project/internal/storage/postgres"
	model "project/internal/storage/postgres/models"
	transport "project/internal/transport"
	"strconv"
	"strings"
)

func main() {
	dsn := "host=localhost port=5432 user=postgres password=qwe123 dbname=feedbox sslmode=disable"
	db, err := storage.NewStoragePG(dsn)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.InitTables(); err != nil {
		fmt.Printf("Failed to initialize tables: %v", err)
	}

	// Запуск сервера
	transport.NewServer()

	// POST и GET /projects
	http.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// var req ProjectRequest
			var req model.Project

			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, " JSON", http.StatusBadRequest)
				return
			}

			err = postgres.CreateProject(db.DB, req.Title, req.Description)
			if err != nil {
				http.Error(w, " Ошибка сохранения", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"Проект создан!"}`))

		case http.MethodGet:
			projects, err := postgres.GetProjects(db.DB)
			if err != nil {
				http.Error(w, " Ошибка получения проектов", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(projects)

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	// DELETE /projects/{id}
	http.HandleFunc("/projects/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		// Парсинг ID из URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 {
			http.Error(w, "Неверный путь", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(parts[2])
		if err != nil {
			http.Error(w, "ID", http.StatusBadRequest)
			return
		}

		err = postgres.DeleteProject(db.DB, id)
		if err != nil {
			http.Error(w, "Ошибка при удалении", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"Удалено успешно"}`))
	})

	// fmt.Println("📡 Сервер запущен на http://localhost:8080 (чистый Go)")
	// http.ListenAndServe(":8080", nil)
}
