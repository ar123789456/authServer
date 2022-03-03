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

	// var us models.UserInfo
	// var user models.User
	// user.Name = "qaz"
	// user.Password = "123"
	// fmt.Println(user.Create())
	// us.Name = "qaz"
	// us.Password = "123"
	// a, _ := json.Marshal(us)
	// fmt.Println(string(a))
}
