package main

import (
	"encoding/json"
	"fmt"
	contract "gohttp3/examples/clientServerContract"
	"gohttp3/server"
	"net/http"
	"strconv"
)

var tasks = map[int]contract.Task{
	1: {ID: "one", Description: "First task", Completed: true},
	2: {ID: "two", Description: "Second task", Completed: true},
	3: {ID: "three", Description: "Third task", Completed: false},
	4: {ID: "four", Description: "Fourth task", Completed: true},
}

func main() {
	mux := createMux()
	server.RunQuicHTTP3Server("../../tlsKeysAndCertificates/my-server-cert.crt", "../../tlsKeysAndCertificates/my-server-key.pem", mux)
}

func createMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", getTasks)
	mux.HandleFunc("GET /task/{id}", getTask)
	mux.HandleFunc("POST /task/create", postTask)
	return mux
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Error encountered", http.StatusInternalServerError)
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	fmt.Println("task id:", idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Error encountered: "+err.Error(), http.StatusInternalServerError)
	}
	task, ok := tasks[id]

	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Error encountered", http.StatusInternalServerError)
	}
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var task contract.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid task data", http.StatusBadRequest)
		return
	}
	id := len(tasks)
	id = id + 1
	tasks[id] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Error encountered", http.StatusInternalServerError)
	}
}
