package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv" // strconv package
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

// Item represents a simple item in our API
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB

func initDB() {
	// Connect to MySQL database
	connectionString := "root:root@tcp(127.0.0.1:3306)/testwork"
	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	// Check the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the MySQL database")
}

func main() {
	// Initialize the mux router
	router := mux.NewRouter()

	// Initialize the database connection
	initDB()

	// Define routes
	router.HandleFunc("/items", getItems).Methods("GET")     // Get all items
	router.HandleFunc("/items/{id}", getItem).Methods("GET") // Get an item by ID
	router.HandleFunc("/items", createItem).Methods("POST")  // Create a new item
        router.HandleFunc("/items", updateItem).Methods("PUT")    // Update an item
	router.HandleFunc("/items/{id}", deleteItem).Methods("DELETE") // Delete an item

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Handler functions for different routes

// getItems returns all items
func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Fetch items from the database
	rows, err := db.Query("SELECT id, name FROM items")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(items)
}

// getItem returns a specific item by ID
func getItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// Fetch item from the database by ID
	var item Item
	err := db.QueryRow("SELECT id, name FROM items WHERE id = ?", params["id"]).Scan(&item.ID, &item.Name)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(item)
}

// createItem creates a new item
func createItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the request body into an Item struct
	var newItem Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		log.Fatal(err)
	}

	// Insert the new item into the database
	result, err := db.Exec("INSERT INTO items (name) VALUES (?)", newItem.Name)
	if err != nil {
		log.Fatal(err)
	}

	// Get the ID of the last inserted item
	newItemID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	// Set the ID in the newItem struct
	newItem.ID = strconv.FormatInt(newItemID, 10)

	// Encode and return the new item with the generated ID
	json.NewEncoder(w).Encode(newItem)
}

// Update an existing item
func updateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the request body into an Item struct
	var updatedItem Item
	err := json.NewDecoder(r.Body).Decode(&updatedItem)
	if err != nil {
		log.Fatal(err)
	}

	// Update the item in the database by ID
	_, err = db.Exec("UPDATE items SET name = ? WHERE id = ?", updatedItem.Name, updatedItem.ID)
	if err != nil {
		log.Fatal(err)
	}

	

	json.NewEncoder(w).Encode(updatedItem)
}

// Delete an existing item
func deleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// Delete the item from the database by ID
	_, err := db.Exec("DELETE FROM items WHERE id = ?", params["id"])
	if err != nil {
		log.Fatal(err)
	}

	// Return a success message
	json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted successfully"})
}
