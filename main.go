package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Task struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	UserFullName string `json:"user_full_name"`
	IsCompleted  bool   `json:"is_completed"`
	IsDeleted    bool   `json:"is_deleted"`
}

type DefaultResponse struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/tasks", TasksHandler)
	http.HandleFunc("/tasks/", TaskHandler)
	fmt.Println("Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

var tasks []Task

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetAllTasks(w, r)
	case "POST":
		AddTask(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetTask(w, r)
	case "PUT":
		UpdateTask(w, r)
	case "DELETE":
		DeleteTask(w, r)
	case "PATCH":
		PatchTask(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {

	t1 := Task{
		ID:           1,
		Title:        "Task 1",
		Description:  "Description 1",
		UserFullName: "Name Full Name 1",
		IsCompleted:  false,
		IsDeleted:    false,
	}

	tasks = append(tasks, t1)

	t2 := Task{
		ID:           2,
		Title:        "Task 2",
		Description:  "Description 2",
		UserFullName: "Name Full Name 2",
		IsCompleted:  false,
		IsDeleted:    false,
	}

	tasks = append(tasks, t2)

	jsonBody, err := json.Marshal(&tasks)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBody)
	if err != nil {
		fmt.Println(err)
	}

}

func GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, task := range tasks {
		if task.ID == id {
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	var t Task

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks = append(tasks, t)
	var response DefaultResponse
	response.Message = "Task successfully added"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println(err)
	}

}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var t Task

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	t.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		fmt.Println("Error writing response", http.StatusInternalServerError)
		return
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Println(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(DefaultResponse{Message: "Task successfully deleted"})
	if err != nil {
		fmt.Println(err)
	}
}

func PatchTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			if title, ok := updates["title"].(string); ok {
				task.Title = title
			}
			if description, ok := updates["description"].(string); ok {
				task.Description = description
			}
			if userFullName, ok := updates["user_full_name"].(string); ok {
				task.UserFullName = userFullName
			}
			if isCompleted, ok := updates["is_completed"].(bool); ok {
				task.IsCompleted = isCompleted
			}
			if isDeleted, ok := updates["is_deleted"].(bool); ok {
				task.IsDeleted = isDeleted
			}

			tasks[i] = task

			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}
