package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"wbl0-server/storage/broker"
	"wbl0-server/storage/db"
)

func RunServer() {
	// Define the directory where your static files (HTML, CSS, JS) are located
	fs := http.FileServer(http.Dir("static"))

	// Create a new HTTP server
	http.Handle("/", fs)
	http.HandleFunc("/getOrder", Handler)

	// Start the server on port 8081
	log.Println("Server started on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// Handle HTTP requests
func Handler(w http.ResponseWriter, req *http.Request) {
	queryParameters := req.URL.Query()
	var output db.Output
	// Get the value of a specific parameter, for example, "name"
	data := queryParameters.Get("data")
	if intdata, err := strconv.Atoi(data); err == nil {
		fmt.Printf("String: %s\nInteger: %d\n", data, intdata)
		output = broker.SendOrderRequest(intdata)
	} else {
		fmt.Println("Error:", err)
	}

	// Check if the request is an AJAX request
	if req.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// Marshal the response data to JSON
		responseData, err := json.Marshal(output)
		if err != nil {
			log.Println("Error marshaling response data:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Write the response data
		w.Write(responseData)
		return
	}
}
