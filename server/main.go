package main

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

func ping(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "pong\n")
}

func main() {
	http.HandleFunc("/ping", ping)

	ConnectDB()
	http.HandleFunc("/products", AddProduct)
	http.HandleFunc("/products", UpdateProduct)
	http.HandleFunc("/products", GetProducts)
	http.HandleFunc("/products", DeleteProduct)

	c := cors.Default()
	handler := c.Handler(http.DefaultServeMux)
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
