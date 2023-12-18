package attendance

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// The logoutPage function will handle all requests to the logout page
// If there is existing Cookie in the browser, the logout page will delete the cookie and redirect to the index page
// If there is no existing Cookie in the browser, the logout page will redirect to the index page
func LogoutPage(res http.ResponseWriter, req *http.Request) {
	myCookie, err := req.Cookie("userCookie")
	if err != nil {
		log.Println("Error when getting cookie:", err)
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	delete(mapUserSessions, myCookie.Value)
	delete(userList.CurSession, myCookie.Value)
	SaveData(userList, "data.json")
	myCookie.MaxAge = -1
	myCookie.Value = ""
	log.Println("User successfully logged out")
	http.SetCookie(res, myCookie)
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

// The loginPage function will handle all requests to the login page
// A new cookie will be created to allow user to store their session
func LoginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {

		// using escape string for sanitization
		username := html.EscapeString(req.FormValue("username"))
		fmt.Println(username)
		pwd := req.FormValue("password")

		// Input Validation using regex to check if username is alphanumeric
		// If username contains special characters, the login page will return an error
		isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
		if !isAlphanumeric(username) {
			res.WriteHeader(http.StatusBadRequest)
			tpl.ExecuteTemplate(res, "login.gohtml", "Invalid Account")
			return
		}

		if userList.Users[username].Username == "" {
			res.WriteHeader(http.StatusBadRequest)
			tpl.ExecuteTemplate(res, "login.gohtml", "Invalid Account")
			return
		}

		err := bcrypt.CompareHashAndPassword(userList.Users[username].Password, []byte(pwd))
		if err != nil {
			log.Println("Error in Login")
			res.WriteHeader(http.StatusBadRequest)
			tpl.ExecuteTemplate(res, "login.gohtml", "Invalid Account")
			return
		}

		createSessionCookie(res, username)
		SaveData(userList, "data.json")

		if username == "admin" {
			http.Redirect(res, req, "/admin", http.StatusSeeOther)
			tpl.ExecuteTemplate(res, "admin.gohtml", nil)
			return
		} else {
			http.Redirect(res, req, "/", http.StatusSeeOther)
			return
		}
	}
	tpl.ExecuteTemplate(res, "login.gohtml", nil)
}

// The admin function will handle all requests to the admin page
// The admin page will display all users and their information
func Admin(res http.ResponseWriter, req *http.Request) {
	log.Println("Loading Admin Page")
	curUser := getCurUser(res, req)
	adminData := AdminData{
		CurUser:  curUser,
		UserList: userList,
	}
	err := tpl.ExecuteTemplate(res, "admin.gohtml", adminData)
	if err != nil {
		log.Println("Error in Admin Page", err)
	}
}

// The indexPage function will handle all requests to the index page
// If there is existing Cookie in the browser, the index page will show the user's information and logged in time
// If there is no existing Cookie in the browser, the index page allow user to log in or sign up
func IndexPage(res http.ResponseWriter, req *http.Request) {
	curUser := getCurUser(res, req)
	if curUser.Username != "" {
		if curUser.Username == "admin" {
			adminData := AdminData{
				CurUser:  curUser,
				UserList: userList,
			}
			err := tpl.ExecuteTemplate(res, "admin.gohtml", adminData)
			if err != nil {
				log.Println("Error in Admin Page", err)

			}
		} else {
			err := tpl.ExecuteTemplate(res, "index.gohtml", curUser)
			if err != nil {
				log.Println("Error in Index Page", err)
			}
		}
	} else {
		log.Println("No User Logged In")
		err := tpl.ExecuteTemplate(res, "index.gohtml", nil)
		if err != nil {
			log.Println("Error in Login Page", err)
		}
	}
}

// The signupPage will only be called by Admin to create new students
func SignupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		err := createNewUser(res, req)
		if err != nil {
			log.Println(err)
			return
		}
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(res, "signup.gohtml", nil)
}
