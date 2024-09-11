package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/olivere/elastic/v7"
	"log"
	"net/http"
	"os"
)

// Book represents a document in Elasticsearch
type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

var client *elastic.Client
var ctx = context.Background()

func main() {

	dotEnv := godotenv.Load(".env")
	if dotEnv != nil {
		log.Fatalf("Error loading .env file: %s", dotEnv)
	}

	// Initialize Elasticsearch client
	var err error
	client, err = elastic.NewClient(elastic.SetURL(os.Getenv("elastic_url")), elastic.SetSniff(false))
	if err != nil {
		log.Fatalf("Error creating the client: %v", err)
	}

	// Create a new router
	router := mux.NewRouter()

	// Define CRUD routes
	router.HandleFunc("/books/{id}", createBookHandler).Methods("POST")
	router.HandleFunc("/books/{id}", getBookHandler).Methods("GET")
	router.HandleFunc("/books/{id}", updateBookHandler).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBookHandler).Methods("DELETE")

	// Start HTTP server
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Create a book
func createBookHandler(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	_, err = client.Index().
		Index("books").
		Id(id).
		BodyJson(book).
		Refresh("wait_for").
		Do(ctx)
	if err != nil {
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Book with ID %s created\n", id)
}

// Get a book
func getBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := client.Get().
		Index("books").
		Id(id).
		Do(ctx)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	if res.Found {
		var book Book
		json.Unmarshal(res.Source, &book)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	} else {
		http.Error(w, "Book not found", http.StatusNotFound)
	}
}

// Update a book
func updateBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updates map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err = client.Update().
		Index("books").
		Id(id).
		Doc(updates).
		Do(ctx)
	if err != nil {
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Book with ID %s updated\n", id)
}

// Delete a book
func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := client.Delete().
		Index("books").
		Id(id).
		Refresh("wait_for").
		Do(ctx)
	if err != nil {
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Book with ID %s deleted\n", id)
}
