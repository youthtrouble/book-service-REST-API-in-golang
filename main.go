package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//Book mod
type book struct {
	ID 			string `json:"ID"`	
	Name          string `json:"name"`
	Author       string `json:"author"`
	PublishedAt string `json:"published_at"`
}

type allBooks []book

var books = allBooks{
	{
		Name:          "A man of the people",
		Author:       "Chinua Achebe",
		PublishedAt: time.Now().Local().String(),
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var newBook book
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Enter book details in the correct order to update")
	}

	json.Unmarshal(reqBody, &newBook)

	// Add the newly created event to the array of events
	books = append(books, newBook)

	// Return the 201 created status code
	w.WriteHeader(http.StatusCreated)
	// Return the newly created event
	json.NewEncoder(w).Encode(newBook)
}

func getOneBook(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	bookID := mux.Vars(r)["id"]

	// Get the details from an existing event
	// Use the blank identifier to avoid creating a value that will not be used
	for _, singleBook := range books {
		if singleBook.ID == bookID {
			json.NewEncoder(w).Encode(singleBook)
		}
	}	
}

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	bookID := mux.Vars(r)["id"]
	var updatedBook book
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "kindle enter book data for update")
	}

	json.Unmarshal(reqBody, &updatedBook)

	for i, singleBook := range books {
		if singleBook.ID == bookID {
			singleBook.Author = updatedBook.Author
			singleBook.Name = updatedBook.Name
			books[i] = singleBook
			json.NewEncoder(w).Encode(singleBook)
		}
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	bookID := mux.Vars(r)["id"]

	// Get the details from an existing event
	// Use the blank identifier to avoid creating a value that will not be used
	for i, singleBook := range books {
		if singleBook.ID == bookID {
			books = append(books[:i], books[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", bookID)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/book", createBook).Methods("POST")
	router.HandleFunc("/books", getAllBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getOneBook).Methods("GET")
	router.HandleFunc("/books/{id}", updateBook).Methods("PATCH")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}