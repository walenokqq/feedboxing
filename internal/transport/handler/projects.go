	package handler
	import {
	"net/http"
	}
	

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
