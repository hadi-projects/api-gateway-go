package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/products/")
		log.Printf("[PRODUCT_SERVICE] Menerima request: %s %s", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Service-Name", "Product-Service")

		if r.Method == "GET" && (path == "" || strings.HasPrefix(path, "item/")) {
			if path == "" {
				json.NewEncoder(w).Encode([]map[string]string{
					{"productID": "prod001", "name": "Laptop", "price": "1200.00"},
					{"productID": "prod002", "name": "Mouse", "price": "25.00"},
				})
			} else if path == "item/prod001" {
				json.NewEncoder(w).Encode(map[string]string{"productID": "prod001", "name": "Laptop", "price": "1200.00", "description": "High-end gaming laptop"})
			} else {
				http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
			}
		} else if r.Method == "POST" && path == "" {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "Product created successfully (dummy)", "productID": "prod003"})
		} else {
			http.Error(w, "Endpoint tidak ditemukan di Product Service", http.StatusNotFound)
		}
	})
	port := "8082"
	fmt.Printf("Product Service berjalan di port :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
