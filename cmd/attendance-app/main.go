package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"attendanceapp/internal/attendance" // Add import statement for the attendance package

	"github.com/gorilla/mux"
)

var tpl *template.Template
var userList attendance.UserList
var mapUserSessions = map[string]string{}

func main() {
	router := mux.NewRouter()
	router.Use(attendance.LoggingMiddleware)
	// The main function will listen to port 5332 and handle all requests to the server
	router.Handle("/", http.HandlerFunc(attendance.IndexPage))
	router.Handle("/signup", http.HandlerFunc(attendance.SignupPage))
	router.Handle("/login", http.HandlerFunc(attendance.LoginPage))
	router.Handle("/logout", http.HandlerFunc(attendance.LogoutPage))
	router.Handle("/checkin", http.HandlerFunc(attendance.Checkin))
	router.Handle("/admin", http.HandlerFunc(attendance.Admin))
	router.Handle("/generate", http.HandlerFunc(attendance.Generate))

	router.Handle("/favicon.ico", http.NotFoundHandler())

	certFile := "cert/cert.pem"
	keyFile := "cert/key.pem"

	log.Println("Listening on port 5332..")

	// http.ListenAndServe(":5332", nil) - HTTP

	err := http.ListenAndServeTLS(":5332", certFile, keyFile, router) // HTTPS
	if err != nil {
		fmt.Println("Error when listening to port:", err)
	}
}
