package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"todo_list_api/app/internal/models"
)

type MainPageData struct {
	Tasks []models.Task
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	// Check method is GET
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load html
	tmpl, err := template.ParseFiles("app/static/templates/index.html")
	if err != nil {
		http.Error(w, "cannot parse html", http.StatusInternalServerError)
		return
	}

	// Get data from server
	var mainPageData MainPageData
	
	// Tasks:
	var tasks []models.Task

	resp, err := http.Get("http://localhost:8080/read")
	if err != nil {
		log.Println(err)
		http.Error(w, "cannot get data", http.StatusInternalServerError)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "cannot read data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		http.Error(w, "cannot unmarshal data", http.StatusInternalServerError)
		return
	}

	mainPageData.Tasks = tasks

	tmpl.Execute(w, mainPageData)
}

func CreateTaskForm(w http.ResponseWriter, r *http.Request) {
		// Check method is POST
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.PostFormValue("name")
		desc := r.PostFormValue("desc")

		bodyMap := map[string]interface{}{"name":name,"desc":desc,"is_completed":false}

		bodyJson, err := json.Marshal(bodyMap)
		if err != nil {
			http.Error(w, "cannot marshal data", http.StatusBadRequest)
			return
		}

		body := bytes.NewBuffer(bodyJson)

		// Create task
		resp, err := http.Post("http://localhost:8080/create", "application/json", body)
		if err != nil {
			http.Error(w, err.Error(), resp.StatusCode)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			http.Error(w, "cannot update task", resp.StatusCode)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UpdateStatusForm(w http.ResponseWriter, r *http.Request) {
	// Check method is POST
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PostFormValue("id")

	_, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	compl, err := strconv.ParseBool(r.PostFormValue("is_completed"))
	if err != nil {
		http.Error(w, "cannot get status", http.StatusBadRequest)
		return
	}

	bodyMap := map[string]interface{}{"is_completed":!compl}

	bodyJson, err := json.Marshal(bodyMap)
	if err != nil {
		http.Error(w, "cannot marshal data", http.StatusBadRequest)
		return
	}

	body := bytes.NewBuffer(bodyJson)

	req, err := http.NewRequest("PATCH", "http://localhost:8080/update/" + id, body)
	if err != nil {
		http.Error(w, "cannot send request", http.StatusBadRequest)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "cannot send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		http.Error(w, "cannot update task", resp.StatusCode)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteTaskForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		w.Header().Set("Allow", http.MethodPost)
		return
	}

	id := r.PostFormValue("id")

	_, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("DELETE", "http://localhost:8080/delete/" + id, nil)
	if err != nil {
		http.Error(w, "cannot create request", http.StatusBadRequest)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "cannot send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		http.Error(w, "cannot delete task", resp.StatusCode)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}