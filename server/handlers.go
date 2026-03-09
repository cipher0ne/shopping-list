package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/cipher0ne/shopping-list/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

func AddProduct(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var product models.Product
	err := json.NewDecoder(request.Body).Decode(&product)
	if err != nil {
		http.Error(response, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := ProductsCollection.InsertOne(ctx, product)
	if err != nil {
		http.Error(response, "Database error", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(map[string]any{
		"insertedID": res.InsertedID,
	})
}

func GetProducts(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := ProductsCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(response, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var products []models.Product
	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			http.Error(response, "Error decoding product", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		http.Error(response, "Cursor error", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(products)
}

func UpdateProduct(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPatch {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(request.URL.Path, "/products/")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(response, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(request.Body).Decode(&updates)
	if err != nil {
		http.Error(response, "Invalid JSON", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": updates,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := ProductsCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		http.Error(response, "Database error", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(map[string]any{
		"matchedCount": res.MatchedCount,
		"modifiedCount": res.ModifiedCount,
	})
}

func DeleteProduct(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(request.URL.Path, "/products/")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(response, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := ProductsCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		http.Error(response, "Database error", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(map[string]any{
		"deletedCount": res.DeletedCount,
	})
}
