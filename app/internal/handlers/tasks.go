package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"todo_list_api/app/internal/models"
)

type Repo interface {
	Create([]byte) error
	GetTasks() []*models.Task
	Update(uint64, []byte) error
	Delete(uint64) error
}

func CreateTask(repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check method is POST
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Create task
		if err := repo.Create(body); err != nil {
			http.Error(w, "cannot create task", http.StatusBadRequest)
			return
		}

		// Response
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]string{"message": "Task created successfully"})
		w.Write(response)
	}
}

func ReadTasks(repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check method is GET
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	
		// Send tasks
		resp, err := json.Marshal(repo.GetTasks())
		if err != nil {
			http.Error(w, "cannot marshal tasks", http.StatusInternalServerError)
			return
		}
	
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func UpdateTask(repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check method is PATCH
		if r.Method != http.MethodPatch {
			w.Header().Set("Allow", http.MethodPatch)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "cannot get id", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
	
		// Update fields	
		if err := repo.Update(id, body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	
		// Response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]string{"message": "Task updated successfully"})
		w.Write(response)
	}
}

func DeleteTask(repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check method is DELETE
		if r.Method != http.MethodDelete {
			w.Header().Set("Allows", http.MethodDelete)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	
		// Find and delete task with given id
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "cannot get id", http.StatusBadRequest)
			return
		}

		if err := repo.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]string{"message": "Task deleted successfully"})
			w.Write(response)
		}
	}
}
