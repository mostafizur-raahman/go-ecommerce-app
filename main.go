package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = make(map[int]User)

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Error while creating user", http.StatusBadRequest)
		return
	}

	if _, exists := users[user.ID]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// set in map
	users[user.ID] = user

	fmt.Println("user", user)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func getAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	fmt.Println("All users retrieved:", users)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser User // Fixed typo here
	err = json.NewDecoder(r.Body).Decode(&updatedUser)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, exists := users[userId]

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return // Added return to stop further execution
	}

	if updatedUser.Name != "" {
		user.Name = updatedUser.Name
	}

	if updatedUser.Email != "" {
		user.Email = updatedUser.Email
	}

	users[userId] = user

	// Return the updated user as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, exists := users[userId]

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return // Added return to stop further execution
	}

	delete(users, userId)
	// Return the updated user as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/user/create", createUser).Methods("POST")
	r.HandleFunc("/users", getAllUser).Methods("GET")
	r.HandleFunc("/user/update/{id}", updateUser).Methods("PATCH")
	r.HandleFunc("/user/delete/{id}", deleteUser).Methods("DELETE")

	fmt.Println("Server listen on port 8080") // Fixed typo here
	http.ListenAndServe(":8080", r)           // Pass the router here
}
