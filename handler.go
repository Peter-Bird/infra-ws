// File: handler.go
package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Handler function to use reflection to choose the appropriate service
func Handler(w http.ResponseWriter, r *http.Request, wsName string) {
	startTime := time.Now()
	log.Printf("[%s] Received request", wsName)

	// Extract the last part of the URL path as the service name
	pathParts := strings.Split(r.URL.Path, "/")
	serviceType := pathParts[len(pathParts)-1] // Get the endpoint name
	log.Printf("[%s] Service type requested: %s", wsName, serviceType)

	var input map[string]interface{}

	// Support both GET and POST methods
	if r.Method == http.MethodGet {
		input = make(map[string]interface{})
		for key, values := range r.URL.Query() {
			if len(values) > 0 {
				input[key] = values[0]
			}
		}
		log.Printf("[%s] Request payload from query parameters: %v", wsName, input)
	} else if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			log.Printf("[%s] Error decoding JSON: %v", wsName, err)
			return
		}
		log.Printf("[%s] Request payload from body: %v", wsName, input)
	} else {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		log.Printf("[%s] Unsupported method: %s", wsName, r.Method)
		return
	}

	service, exists := GetService(serviceType)
	if !exists {
		http.Error(w, "Unknown service type", http.StatusBadRequest)
		log.Printf("[%s] Unknown service type: %s", wsName, serviceType)
		return
	}

	result, err := service.Process(input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing request: %v", err), http.StatusInternalServerError)
		log.Printf("[%s] Error processing request with service %s: %v", wsName, serviceType, err)
		return
	}
	log.Printf("[%s] Service %s processed request successfully", wsName, serviceType)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("[%s] Error encoding response: %v", wsName, err)
	}

	endTime := time.Now()
	log.Printf("[%s] Request completed, duration: %v", wsName, endTime.Sub(startTime))
}
