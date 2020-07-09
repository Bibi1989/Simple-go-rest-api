package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bibi1989/restfulapi/connectmongo"
	connectionmongo "github.com/bibi1989/restfulapi/connectmongo"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Connect (Func)
// func Connect() *mongo.Collection {
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

// 	// Connect to MongoDB
// 	client, err := mongo.Connect(context.TODO(), clientOptions)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Check the connection
// 	err = client.Ping(context.TODO(), nil)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Connected to MongoDB!")

// 	collection := client.Database("go_book").Collection("books")

// 	err = client.Disconnect(context.TODO())

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// fmt.Println("Connection to MongoDB closed.")

// 	return collection
// }

// NewBook Struct (Model)
type NewBook struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Isbn   string             `json:"isbn,omitempty" bson:"isbn,omitempty"`
	Title  string             `json:"title,omitempty" bson:"title,omitempty"`
	Author *Author            `json:"author,omitempty" bson:"author,omitempty"`
}

// NewAuthor Struct (Model)
type NewAuthor struct {
	Firstname string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

// Book Struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author Struct (Model)
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// init book slice
var books []Book

// get all books func
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection := connectmongo.Connect()
	// findOptions := options.Find()
	// findOptions.SetLimit(2)
	// ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	col, err := collection.Find(context.Background(), bson.D{{}})

	var results []primitive.M
	for col.Next(context.Background()) {
		var result bson.M
		e := col.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		// fmt.Println("cur..>", cur, "result", reflect.TypeOf(result), reflect.TypeOf(result["_id"]))
		results = append(results, result)

	}

	fmt.Println(col, err)
	json.NewEncoder(w).Encode(results)
}

// get one book func
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// create a book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// var book Book
	var newbook NewBook
	_ = json.NewDecoder(r.Body).Decode(&newbook)
	// _ = json.NewDecoder(r.Body).Decode(&book)
	// book.ID = strconv.Itoa(rand.Intn(10000000))
	// books = append(books, book)
	collection := connectionmongo.Connect()
	col, _ := collection.InsertOne(context.Background(), newbook)
	json.NewEncoder(w).Encode(col)
}

// update a book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

// delete a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	route := mux.NewRouter()

	// mock data - @todo - implement DB
	books = append(books, Book{ID: "1", Title: "Book One", Isbn: "123456", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", Title: "Book Two", Isbn: "654321", Author: &Author{Firstname: "Mary", Lastname: "Doe"}})
	books = append(books, Book{ID: "3", Title: "Book Three", Isbn: "567890", Author: &Author{Firstname: "Peter", Lastname: "Smith"}})

	fmt.Println("Server running at port 8000")

	connectionmongo.Connect()

	route.HandleFunc("/api/v1/books", getBooks).Methods("GET")
	route.HandleFunc("/api/v1/books/{id}", getBook).Methods("GET")
	route.HandleFunc("/api/v1/books", createBook).Methods("POST")
	route.HandleFunc("/api/v1/books/{id}", updateBook).Methods("PUT")
	route.HandleFunc("/api/v1/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", route))
}
