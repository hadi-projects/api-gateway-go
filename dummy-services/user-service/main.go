package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/users/")
		log.Printf("[USER_SERVICE] Menerima request: %s %s", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Service-Name", "User-Service")

		if r.Method == "GET" && path == "profile" {
			json.NewEncoder(w).Encode(map[string]string{"userID": "user123", "name": "John Doe", "email": "john.doe@example.com"})
		} else if r.Method == "GET" && path == "" {
			json.NewEncoder(w).Encode([]map[string]string{
				{"userID": "user123", "name": "John Doe"},
				{"userID": "user456", "name": "Jane Smith"},
			})
		} else if r.Method == "POST" && path == "" {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully (dummy)", "userID": "user789"})
		} else {
			http.Error(w, "Endpoint tidak ditemukan di User Service", http.StatusNotFound)
		}
	})
	port := "8081"
	fmt.Printf("User Service berjalan di port :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
