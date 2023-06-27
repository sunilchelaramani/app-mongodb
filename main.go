// Sample go code that demonstrate CRUD operations using mongodb

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Contact struct represents a contact entity
type Contact struct {
	Name  string
	Email string
	Phone string
}

func main() {
	// MongoDB connection settings
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB server to check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Failed to ping the MongoDB server:", err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get a handle to the "contacts" collection
	collection := client.Database("testdb").Collection("contacts")

	// Insert a contact
	insertContact(collection, Contact{
		Name:  "John Doe",
		Email: "johndoe@example.com",
		Phone: "1234567890",
	})

	// Find all contacts
	contacts, err := findAllContacts(collection)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("All contacts:")
	for _, contact := range contacts {
		fmt.Println(contact)
	}

	// Update a contact
	updateContact(collection, "johndoe@example.com", Contact{
		Name:  "John Doe",
		Email: "johndoe@example.com",
		Phone: "9876543210",
	})

	// Find a contact by email
	email := "johndoe@example.com"
	contact, err := findContactByEmail(collection, email)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Contact found by email:", contact)

	// Delete a contact by email
	deleteContact(collection, email)

	// Close the MongoDB connection
	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Disconnected from MongoDB!")
}

// Insert a contact into the collection
func insertContact(collection *mongo.Collection, contact Contact) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, contact)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Contact inserted successfully!")
}

// Find all contacts in the collection
func findAllContacts(collection *mongo.Collection) ([]Contact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var contacts []Contact
	for cur.Next(ctx) {
		var contact Contact
		err := cur.Decode(&contact)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

// Update a contact in the collection by email
func updateContact(collection *mongo.Collection, email string, updatedContact Contact) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": email}
	update := bson.M{"$set": updatedContact}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Contact updated successfully!")
}

// Find a contact in the collection by email
func findContactByEmail(collection *mongo.Collection, email string) (Contact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var contact Contact
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&contact)
	if err != nil {
		return contact, err
	}

	return contact, nil
}

// Delete a contact in the collection by email
func deleteContact(collection *mongo.Collection, email string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"email": email})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Contact deleted successfully!")
}
