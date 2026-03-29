package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cipher0ne/shopping-list/backend/models"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Handles /products (GET, POST)
func AllProductsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	switch r.Method {
	case http.MethodPost:
		AddProduct(w, r)
	case http.MethodGet:
		GetProducts(w, r)
	default:
		log.Printf("[WARN] Method Not Allowed: %s", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Handles /products/{id} (PATCH, DELETE)
func ProductHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	switch r.Method {
	case http.MethodPatch:
		UpdateProduct(w, r)
	case http.MethodDelete:
		DeleteProduct(w, r)
	default:
		log.Printf("[WARN] Method Not Allowed: %s", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetUserIDFromToken(r *http.Request) (primitive.ObjectID, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return primitive.NilObjectID, http.ErrNoCookie
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return primitive.NilObjectID, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["userID"] == nil {
		return primitive.NilObjectID, jwt.ErrSignatureInvalid
	}

	userIDHex, ok := claims["userID"].(string)
	if !ok {
		return primitive.NilObjectID, jwt.ErrSignatureInvalid
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return primitive.NilObjectID, jwt.ErrSignatureInvalid
	}

	return userID, nil
}

func GetUsershoppingList(userID primitive.ObjectID) (models.ShoppingList, error) {
	var shoppingList models.ShoppingList
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("[DEBUG] Looking for shopping list with userId: %v\n", userID.Hex())
	err := ShoppingListsCollection.FindOne(ctx, bson.M{"userId": userID}).Decode(&shoppingList)
	if err != nil {
		log.Printf("[DEBUG] Shopping list not found for userId: %v, error: %v\n", userID.Hex(), err)
	} else {
		log.Printf("[DEBUG] Found shopping list: %+v\n", shoppingList)
	}
	return shoppingList, err
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromToken(r)
	log.Printf("[DEBUG] AddProduct: userID from token: %v\n", userID.Hex())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// find the shopping list of this user
	shoppingList, err := GetUsershoppingList(userID)
	if err != nil {
		log.Printf("[DEBUG] AddProduct: Shopping list not found for userID: %v\n", userID.Hex())
		http.Error(w, "Shopping list not found", http.StatusNotFound)
		return
	}

	// decode product from request body and give it IDs
	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	product.ID = primitive.NewObjectID()
	product.ShoppingListID = shoppingList.ID

	log.Printf("[DEBUG] AddProduct: Inserting product: %+v\n", product)

	// insert product
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := ProductsCollection.InsertOne(ctx, product)
	if err != nil {
		log.Printf("[DEBUG] AddProduct: Database error: %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"insertedID": res.InsertedID,
	})
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] GetProducts called by %s", r.RemoteAddr)
	userID, err := GetUserIDFromToken(r)
	if err != nil {
		log.Printf("[WARN] Unauthorized access in GetProducts: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// find the shopping list of this user
	shoppingList, err := GetUsershoppingList(userID)
	if err != nil {
		log.Printf("[WARN] Shopping list not found in GetProducts for userID: %v", userID.Hex())
		http.Error(w, "Shopping list not found", http.StatusNotFound)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// find products that belong to this user's list
	cursor, err := ProductsCollection.Find(ctx, bson.M{"shoppingListID": shoppingList.ID})
	if err != nil {
		log.Printf("[ERROR] Database error in GetProducts: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		log.Printf("[ERROR] Database error in GetProducts (cursor): %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Returning %d products for userID: %v", len(products), userID.Hex())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] UpdateProduct called by %s", r.RemoteAddr)
	userID, err := GetUserIDFromToken(r)
	if err != nil {
		log.Printf("[WARN] Unauthorized access in UpdateProduct: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// find the shopping list of this user
	shoppingList, err := GetUsershoppingList(userID)
	if err != nil {
		log.Printf("[WARN] Shopping list not found in UpdateProduct for userID: %v", userID.Hex())
		http.Error(w, "Shopping list not found", http.StatusNotFound)
		return
	}

	// get product ID from URL
	id := strings.TrimPrefix(r.URL.Path, "/products/")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("[WARN] Invalid product ID in UpdateProduct: %v", id)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// decode updated product data from request body
	var updates map[string]any
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Printf("[WARN] Invalid JSON in UpdateProduct: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	update := bson.M{"$set": updates}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// only update products that belong to the user
	filter := bson.M{"_id": objID, "shoppingListID": shoppingList.ID}

	res, err := ProductsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("[ERROR] Database error in UpdateProduct: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Updated product %v for userID: %v, matched: %d, modified: %d", id, userID.Hex(), res.MatchedCount, res.ModifiedCount)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"matchedCount":  res.MatchedCount,
		"modifiedCount": res.ModifiedCount,
	})
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] DeleteProduct called by %s", r.RemoteAddr)
	userID, err := GetUserIDFromToken(r)
	if err != nil {
		log.Printf("[WARN] Unauthorized access in DeleteProduct: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// find the shopping list of this user
	shoppingList, err := GetUsershoppingList(userID)
	if err != nil {
		log.Printf("[WARN] Shopping list not found in DeleteProduct for userID: %v", userID.Hex())
		http.Error(w, "Shopping list not found", http.StatusNotFound)
		return
	}

	// get product ID from URL
	id := strings.TrimPrefix(r.URL.Path, "/products/")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("[WARN] Invalid product ID in DeleteProduct: %v", id)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// filter out the product by id passed in the URL and only if it belong to the user's list
	filter := bson.M{"_id": objID, "shoppingListID": shoppingList.ID}

	res, err := ProductsCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("[ERROR] Database error in DeleteProduct: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Deleted product %v for userID: %v, deletedCount: %d", id, userID.Hex(), res.DeletedCount)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"deletedCount": res.DeletedCount,
	})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] RegisterHandler called by %s", r.RemoteAddr)
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("[WARN] Invalid JSON in RegisterHandler: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[ERROR] Password hashing failed in RegisterHandler: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check if user already exists
	var existingUser models.User
	result := UsersCollection.FindOne(ctx, bson.M{"email": req.Email})
	// if user was found Decode will return nil error
	err = result.Decode(&existingUser)
	if err == nil {
		log.Printf("[WARN] RegisterHandler: User already exists: %s", req.Email)
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// create user
	user := models.User{
		ID:       primitive.NewObjectID(),
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	_, err = UsersCollection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("[ERROR] Failed to insert user in RegisterHandler: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// create shopping list for user
	shoppingList := models.ShoppingList{
		ID:     primitive.NewObjectID(),
		UserID: user.ID,
	}
	_, err = ShoppingListsCollection.InsertOne(ctx, shoppingList)
	if err != nil {
		log.Printf("[ERROR] Failed to insert shopping list in RegisterHandler: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Registered new user: %s, userID: %v", req.Email, user.ID.Hex())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"userID": user.ID,
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] LoginHandler called by %s", r.RemoteAddr)
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("[WARN] Invalid JSON in LoginHandler: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check if email exists in the database
	var user models.User
	err = UsersCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		log.Printf("[WARN] Invalid email in LoginHandler: %s", req.Email)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// compare password with hashed password in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Printf("[WARN] Invalid password in LoginHandler for email: %s", req.Email)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID.Hex(),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("[ERROR] Failed to sign JWT in LoginHandler: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] User logged in: %s, userID: %v", req.Email, user.ID.Hex())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"token": tokenString,
	})
}
