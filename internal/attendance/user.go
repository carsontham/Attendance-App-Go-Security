package attendance

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// the submit button in singupPage is linked to this function, which creates a new user
// the data will be stored in the json file which acts as a database
func createNewUser(res http.ResponseWriter, req *http.Request) error {
	username := req.FormValue("username")
	password := req.FormValue("password")
	firstname := req.FormValue("firstname")
	lastname := req.FormValue("lastname")

	if username != "" {
		// check if username exist/ taken
		if _, ok := userList.Users[username]; ok {
			http.Error(res, "Username already taken", http.StatusForbidden)
			return errors.New("username already taken")
		}

		bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return errors.New("internal server error")
		}

		curUser := User{
			Username:  username,
			Password:  bPassword,
			FirstName: firstname,
			LastName:  lastname,
		}
		userList.Users[username] = curUser
		fmt.Printf("New User %v Created", username)
		//createSessionCookie(res, username)
		SaveData(userList, "data.json")

	}
	return nil
}

// The checkin function will handle all requests to the checkin page
// Students can check in to the class by clicking the check in button
// The check in button will set the CheckedIn field to true and the TimeIn field to the current time
func Checkin(res http.ResponseWriter, req *http.Request) {
	curUser := getCurUser(res, req)
	curUser.CheckedIn = true
	curUser.TimeIn = time.Now().Format("2006-01-02 15:04:05")

	userList.Users[curUser.Username] = curUser
	SaveData(userList, "data.json")
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

// this function will try to retrieve existing cookies if user is logged in
// if there is existing cookie, the function will return the current user
func getCurUser(res http.ResponseWriter, req *http.Request) User {
	myCookie, err := req.Cookie("userCookie")
	if err != nil {
		fmt.Println("No Cookie here")
		return User{}
	}
	username := userList.CurSession[myCookie.Value]
	curUser := userList.Users[username]
	return curUser
}
