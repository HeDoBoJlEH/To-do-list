package repo

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"todo_list_api/app/internal/models"
)

type TasksRepo struct {
	Tasks []*models.Task `json:"tasks"`
	id    uint64
}

func TasksRepoInit() (*TasksRepo, error) {
	var repo TasksRepo

	path := filepath.Join("app", "internal", "repo", "tasks.json")

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &repo.Tasks)
	if err != nil {
		return nil, err
	}

	repo.setId()

	return &repo, nil
}

func (t *TasksRepo) SaveFile() error {
	data, err := json.Marshal(t.Tasks)
	if err != nil {
		return err
	}

	path := filepath.Join("app", "internal", "repo", "tasks.json")

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	return nil
}

func (t *TasksRepo) Create(data []byte) error {
	task := models.Task{Id: t.id, IsCompleted: false}

	err := json.Unmarshal(data, &task)
	if err != nil {
		return errors.New("cannot unmarshal json")
	}

	t.id++

	t.Tasks = append(t.Tasks, &task)

	if err := t.SaveFile(); err != nil {
		return errors.New("cannot save file")
	}
	
	return nil
}

func (t *TasksRepo) Update(id uint64, data []byte) error {
	var taskRef *models.Task
	 
	for _, task := range t.Tasks {
		if task.Id == id {
			taskRef = task
		}
	}

	if taskRef == nil {
		return errors.New("task not found")
	}

	// Update fields
	if err := json.Unmarshal(data, taskRef); err != nil {
		return errors.New("cannot unmarshal data")
	}

	if err := t.SaveFile(); err != nil {
		return errors.New("cannot save file")
	}

	return nil
}

func (t *TasksRepo) Delete(id uint64) error {
	for i, task := range t.Tasks {
		if task.Id == id {
			t.Tasks = slices.Delete(t.Tasks, i, i+1)

			if err := t.SaveFile(); err != nil {
				return errors.New("cannot save file")
			}

			return nil
		}
	}

	return errors.New("task not found")
}

func (t *TasksRepo) GetTasks() []*models.Task {
	return t.Tasks
}

func (t *TasksRepo) setId() {
	var maxId uint64 = 0

	for _, task := range t.Tasks {
		maxId = max(maxId, task.Id)
	}

	t.id = maxId + 1
}