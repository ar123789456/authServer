package main

import (
	handler "auth/handlers"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/login", handler.Access)
	http.HandleFunc("/refrash", handler.Refresh)

	log.Println("ListenAndServe: localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
