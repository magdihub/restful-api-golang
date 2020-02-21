package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func main() {

	// init connection with mysql
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/go_with_mysql")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// fmt.Println("Connected to mysql db!")

	// init routers
	r := mux.NewRouter()

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/book/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/book/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/book/{id}", deleteBook).Methods("DELETE")

	// fmt.Printf("server running on Port %s\n", port)
	log.Fatal(http.ListenAndServe(":8000", r))

}

// struct Models

type Book struct {
	ID     string  `json:"_id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

// Handlers
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("inside getBooks func!")

	results, err := db.Query("SELECT `books`.`_id`, `title`, `isbn`, `firstname`, `lastname` FROM `books` LEFT JOIN `authors` ON `books`.`author_id` = `authors`.`_id` WHERE 1;")

	if err != nil {
		panic(err.Error())
	}

	defer results.Close()

	var books []Book

	for results.Next() {
		var book Book
		var author Author

		err := results.Scan(
			&book.ID,
			&book.Isbn,
			&book.Title,
			&author.FirstName,
			&author.LastName)
		if err != nil {
			panic(err.Error())
		}

		book.Author = &author

		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("inside getBook func!")

	params := mux.Vars(r)

	result, err := db.Query("SELECT `books`.`_id`, `title`, `isbn`, `firstname`, `lastname` FROM `books` LEFT JOIN `authors` ON `books`.`author_id` = `authors`.`_id` WHERE `books`.`_id` = ?;", params["id"])

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var book Book
	var author Author

	for result.Next() {

		err := result.Scan(&book.ID, &book.Isbn, &book.Title, &author.FirstName, &author.LastName)

		if err != nil {
			panic(err.Error())
		}

		book.Author = &author

	}

	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("inside createBook func!")
	fmt.Printf("%T\n", r.Body)

	type NewBook struct {
		ID       string `json:"_id"`
		Isbn     string `json:"isbn"`
		Title    string `json:"title"`
		AuthorID int    `json:"author_id"`
	}

	type Resp struct {
		status int
		msg    string
	}

	var book NewBook

	json.NewDecoder(r.Body).Decode(&book)

	insert, err := db.Query("INSERT INTO `books` (`_id`, `title`, `isbn`, `author_id`, `updatedAt`, `createdAt`) VALUES (null, ?, ?, ?, NOW(), NOW());", book.Title, book.Isbn, book.AuthorID)
	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()

	payload := Resp{}
	payload.status = 1
	payload.msg = "Book Inserted!"

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(payloadJson)
	// json.NewDecoder("{ \"status\": 1, \"msg\": \"Book Inserted!\" }").Decode(&payload)
	w.WriteHeader(http.StatusOK)
	w.Write(payloadJson)

}

// Update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Whoa, Go is neat!")
}

// respondwithJSON write json response format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	fmt.Println(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
