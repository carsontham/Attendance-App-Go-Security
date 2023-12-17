package main

import (
	"fmt"
	"html/template"
	"net/http"

	"attendanceapp/internal/attendance" // Add import statement for the attendance package
)

var tpl *template.Template
var userList attendance.UserList
var mapUserSessions = map[string]string{}

func main() {
	// The main function will listen to port 5332 and handle all requests to the server
	http.Handle("/", http.HandlerFunc(attendance.IndexPage))
	http.Handle("/signup", http.HandlerFunc(attendance.SignupPage)) // Assuming SignupPage is defined in the attendance package
	http.Handle("/login", http.HandlerFunc(attendance.LoginPage))   // Assuming LoginPage is defined in the attendance package
	http.Handle("/logout", http.HandlerFunc(attendance.LogoutPage)) // Assuming LogoutPage is defined in the attendance package
	http.Handle("/checkin", http.HandlerFunc(attendance.Checkin))   // Assuming Checkin is defined in the attendance package
	http.Handle("/admin", http.HandlerFunc(attendance.Admin))       // Assuming Admin is defined in the attendance package
	http.Handle("/generate", http.HandlerFunc(attendance.Generate)) // Assuming Generate is defined in the attendance package

	http.Handle("/favicon.ico", http.NotFoundHandler())

	certFile := "cert/cert.pem"
	keyFile := "cert/key.pem"

	fmt.Println("Listening on port 5332..")
	// http.ListenAndServe(":5332", nil)

	err := http.ListenAndServeTLS(":5332", certFile, keyFile, nil)
	if err != nil {
		fmt.Println("Error when listening to port:", err)
	}
}
