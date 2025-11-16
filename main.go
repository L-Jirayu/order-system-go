package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Go Gateway started on port 8080")

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		resultChan := runTestClientAsync() // async
		result := <-resultChan             // await

		fmt.Fprintln(w, "Running gRPC test client...")
		fmt.Fprintln(w)
		fmt.Fprintln(w, result)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
