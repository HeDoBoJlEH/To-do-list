package models

type Task struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	IsCompleted bool   `json:"is_completed"`
}
