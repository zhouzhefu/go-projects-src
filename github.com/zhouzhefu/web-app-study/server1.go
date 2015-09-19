package main 

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	//"strings"
)

func sayHallo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Web")
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("client Method: ", r.Method)
	if r.Method == "GET" { //GET means user just reach login panel
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else { //POST means user try to login
		r.ParseForm() //by default form will not be parsed until call out, 
		fmt.Println("username: ", r.Form["username"]) //only after ParseForm() was called, 
		fmt.Println("password: ", r.Form["password"]) //these fields can read value
	}
}

func main() {
	http.HandleFunc("/", sayHallo)
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}