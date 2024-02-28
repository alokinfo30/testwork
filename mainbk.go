package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// Item represents a simple item in our API
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var items = []Item{
	{ID: "1", Name: "Item 1"},
	{ID: "2", Name: "Item 2"},
	{ID: "3", Name: "Item 3"},
}

func main() {
	// Initialize the mux router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/items", getItems).Methods("GET")     // Get all items
	router.HandleFunc("/items/{id}", getItem).Methods("GET") // Get an item by ID
	router.HandleFunc("/items", createItem).Methods("POST")  // Create a new item
	router.HandleFunc("/test", testEndpoint).Methods("GET")  // Test endpoint

	// Start the server on a different port, e.g., :8081
	port := ":8081"
	log.Printf("Server is starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

// Handler functions for different routes

// getItems returns all items
func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("GET request to /items")
	json.NewEncoder(w).Encode(items)
}

// getItem returns a specific item by ID
func getItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	log.Printf("GET request to /items/%s", params["id"])
	for _, item := range items {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Item{})
}

// createItem creates a new item
func createItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)
	items = append(items, item)
	log.Printf("POST request to /items. New item created: %+v", item)
	json.NewEncoder(w).Encode(item)
}

// testEndpoint is a test endpoint for demonstration
func testEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("GET request to /test")
	json.NewEncoder(w).Encode(map[string]string{"message": "This is a test endpoint"})
}

// Example test function
func TestGetItems(t *testing.T) {
	req, err := http.NewRequest("GET", "/items", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getItems)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"id":"1","name":"Item 1"},{"id":"2","name":"Item 2"},{"id":"3","name":"Item 3"}]`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
