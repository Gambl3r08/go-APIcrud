package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int
	Name    string
	Content string
}

type allTask []task

var tasks = allTask{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Some content",
	},
}

//enviar las tareas
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

//Crear tareas
func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid Task")
	}
	json.Unmarshal(reqBody, &newTask)
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)

}

//Delete Task

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "invalid ID")
		return
	}

	for index, taskf := range tasks {
		if taskf.ID == taskID {
			tasks = append(tasks[:index], tasks[index+1:]...)
			fmt.Fprintf(w, "the task with ID %v has been remove", taskID)
		}
	}

}

//Update by ID

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updatedTask task
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "Please enter valid data ")
		return
	}

	json.Unmarshal(reqBody, &updatedTask)

	for i, taskf := range tasks {
		if taskf.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			updatedTask.ID = taskID
			tasks = append(tasks, updatedTask)

			fmt.Fprintf(w, "the Task with ID %v has been updated", taskID)
		}
	}

}

//Taks by ID
func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for _, taskf := range tasks {
		if taskf.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(taskf)
		}
	}

}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Wellcome to my API")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))

}
