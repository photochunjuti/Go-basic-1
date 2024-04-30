package main

import (
	"context"
	"fmt"
	"gobasic/product"
	"log"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connnectMongoDB() (*mongo.Client, error) {
	// set client options
	client_options := options.Client().ApplyURI("mongodb://localhost:27017")

	// connect to mongoDB
	client, err := mongo.Connect(context.Background(), client_options)
	if err != nil {
		return nil, err
	}
	// Ping the MongoDB server to verify that the client can connect
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB!")

	return client, nil
}

func add_data(collection *mongo.Collection, id *int, con *string) {
	var name, category, key string
	var price int

	fmt.Print("Input product name : ")
	fmt.Scan(&name)
	fmt.Print("Input product prices : ")
	fmt.Scan(&price)
	fmt.Print("Input category : ")
	fmt.Scan(&category)
	// Create a test product object
	product := product.Product{
		Id:       *id,
		Name:     name,
		Price:    price,
		Category: category,
	}

	// Insert the Person object into the collection
	insertResult, err := collection.InsertOne(context.Background(), product)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted document ID:", insertResult.InsertedID)
	*id++

	fmt.Print("Continue adding data (y or n) : ")
	fmt.Scan(&key)
	fmt.Println("")
	*con = key
}
func delete_data(collection *mongo.Collection, id int, deleted_con *string) {
	var key string

	// filter := bson.M{"name": "Example"}
	target_data := bson.M{"id": id}
	deleteResult, err := collection.DeleteOne(context.Background(), target_data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %d document(s)\n", deleteResult.DeletedCount)

	fmt.Print("Continue adding data (y or n) : ")
	fmt.Scan(&key)
	fmt.Println("")
	*deleted_con = key
}

func get_product(w http.ResponseWriter, r *http.Request) {
	// Connect to MongoDB
	client, err := connnectMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Access database
	database := client.Database("companyDB")
	collection := database.Collection("products")

	// Query mongoDB for product data
	// fetch all products without any filters
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}
	defer cursor.Close(context.Background())

	// convert mongoDB cursor to a slice of product structs
	var products []product.Product
	if err := cursor.All(context.Background(), &products); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}

	// Marshal products slice to JSON
    jsonData, err := json.Marshal(products)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	// Set Content-Type header and write JSON response
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}

func main() {
	// Set up HTTP server
    http.HandleFunc("/api/products", get_product)
    log.Fatal(http.ListenAndServe(":8080", nil))
	
}



