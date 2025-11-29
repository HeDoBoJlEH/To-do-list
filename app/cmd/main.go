package main

import (
	"fmt"
	"log"
	"net/http"
	"todo_list_api/app/internal/handlers"
	"todo_list_api/app/internal/repo"
)

func main() {
	tasksRepo, err := repo.TasksRepoInit()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./app/static/"))))
	http.HandleFunc("/", handlers.MainPage)
	http.HandleFunc("/create", handlers.CreateTask(tasksRepo))
	http.HandleFunc("/form/create", handlers.CreateTaskForm)
	http.HandleFunc("/form/update/status", handlers.UpdateStatusForm)
	http.HandleFunc("/form/delete", handlers.DeleteTaskForm)
	http.HandleFunc("/read", handlers.ReadTasks(tasksRepo))
	http.HandleFunc("/update/{id}", handlers.UpdateTask(tasksRepo))
	http.HandleFunc("/delete/{id}", handlers.DeleteTask(tasksRepo))

	fmt.Println("Started server on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
