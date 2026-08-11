package main

import (
	"fmt"
	"log"
	"net/http"

	"go-app-layer/apps/handlers"
)

func main() {
	fmt.Println("Go Gateway started on port 8080")

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		resultChan := runTestClientAsync()
		result := <-resultChan

		fmt.Fprintln(w, "Running gRPC test client...")
		fmt.Fprintln(w)
		fmt.Fprintln(w, result)
	})

	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/profile", handlers.Profile)
	http.HandleFunc("/addresses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.Addresses(w, r)
		case http.MethodPost:
			handlers.AddAddress(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/restaurants", handlers.Restaurants)
	http.HandleFunc("/restaurants/", handlers.RestaurantMenu)

	http.HandleFunc("/payments", handlers.CreatePayment)
	http.HandleFunc("/payments/", handlers.GetPayment)

	http.HandleFunc("/orders", handlers.CreateOrder)
	http.HandleFunc("/orders/", handlers.GetOrder)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
