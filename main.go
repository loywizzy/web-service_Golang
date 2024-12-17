package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// โครงสร้างข้อมูล User
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ตัวอย่างข้อมูล (จำลองฐานข้อมูล)
var users = []User{
	{ID: "1", Name: "Alice", Email: "alice@example.com"},
	{ID: "2", Name: "Bob", Email: "bob@example.com"},
}

// ฟังก์ชัน handler: แสดง User ทั้งหมด
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// ฟังก์ชัน handler: เพิ่ม User ใหม่
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	users = append(users, newUser)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

// ฟังก์ชัน handler: แสดง User ตาม ID
func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	for _, user := range users {
		if user.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var updatedUser User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for index, user := range users {
		if user.ID == id {
			// อัปเดตข้อมูล User ใน Slice
			users[index] = updatedUser
			users[index].ID = id // คงค่า ID เดิมเอาไว้
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(users[index])
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	for index, user := range users {
		if user.ID == id {
			// ลบ User ออกจาก Slice โดยใช้ append
			users = append(users[:index], users[index+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func main() {
	// สร้าง Router ด้วย Mux
	r := mux.NewRouter()

	// กำหนด Route
	r.HandleFunc("/users", getUsersHandler).Methods("GET")           // ดึง Users ทั้งหมด
	r.HandleFunc("/users", createUserHandler).Methods("POST")        // เพิ่ม User ใหม่
	r.HandleFunc("/users/{id}", getUserByIDHandler).Methods("GET")   // ดึง User ตาม ID
	r.HandleFunc("/users/{id}", updateUserHandler).Methods("PUT")    // อัปเดตข้อมูล User ตาม ID
	r.HandleFunc("/users/{id}", deleteUserHandler).Methods("DELETE") // ลบ User ตาม ID

	// Start Server
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
