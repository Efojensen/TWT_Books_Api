package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity uint8  `json:"quantity"`
}

var allBooks = []book{
	{ID: "1", Title: "Alex Rider  Scorpia", Author: "Anthony Horowitz", Quantity: 2},
	{ID: "2", Title: "The Red Pyramid", Author: "Rick Riordan", Quantity: 5},
	{ID: "3", Title: "Harry Potter: The Prisoner of Azkaban", Author: "J.K Rowling", Quantity: 8},
}

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		for _, book := range allBooks {
			json.NewEncoder(w).Encode(book)
		}
	}
}

func getBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	w.Header().Set("Content-Type", "application/json")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "missing path parameter",
		})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string {
			"error": "book not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func getBookById(id string) (*book, error) {
	for index, book := range allBooks {
		if book.ID == id {
			return &allBooks[index], nil
		}
	}

	return nil, errors.New("book not found")
}

func checkoutBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		json.NewEncoder(w).Encode(map[string]string {
			"error" : "wrong request type",
		})
		return
	}

	queryParams := r.URL.Query()
	id := queryParams.Get("id")

	book, err := getBookById(id)
	w.Header().Set("Content-type", "application/json")

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string {
			"error" : "book not found",
		})
		return
	}

	book.Quantity -= 1

	json.NewEncoder(w).Encode(book)
}

func main() {
	http.HandleFunc("/books", getAllBooks)
	http.HandleFunc("/books/{id}", getBook)
	http.HandleFunc("/checkout", checkoutBook)
	fmt.Println("Server running on localhost:8080")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})

	http.ListenAndServe(":8080", nil)
}
