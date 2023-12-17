package attendance

import (
	"errors"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"regexp"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// The tpl variable will store all templates in the templates folder
var tpl *template.Template

// The userList variable will store all users in the data.json file
var userList UserList

// The mapUserSessions variable will store the latest session in the server
var mapUserSessions = map[string]string{}

// The init function will run all templates in the templates folder
// It will also load the data from the data.json file and save it to the userList variable
func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	LoadData()
	// fmt.Println(userList)
}

// The User Struct contains the information for each user
type User struct {
	Username  string `json:"username"`
	Password  []byte `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CheckedIn bool   `json:"checked_in"`
	TimeIn    string `json:"time_in"`
}

// The UserList Struct contains a list of users
type UserList struct {
	Users      map[string]User   `json:"users"`
	CurSession map[string]string `json:"cur_session"`
}

// The Output Struct contains the information required to output the data into
type Output struct {
	Username  string
	FirstName string
	LastName  string
	CheckedIn bool
	TimeIn    string
}

// The AdminData Struct contains the information required to display the admin page
type AdminData struct {
	CurUser  User
	UserList UserList
}

// this function will read the chosen radio button and generate the attendance list in the chosen format
// the file will be returned to the client via serveFiles()
func Generate(res http.ResponseWriter, req *http.Request) {
	var outputList []Output

	for user, _ := range userList.Users {
		if user == "admin" {
			continue
		}
		//fmt.Println(user, userList.Users[user].CheckedIn, userList.Users[user].TimeIn)
		o := Output{
			Username:  userList.Users[user].Username,
			FirstName: userList.Users[user].FirstName,
			LastName:  userList.Users[user].LastName,
			CheckedIn: userList.Users[user].CheckedIn,
			TimeIn:    userList.Users[user].TimeIn,
		}
		outputList = append(outputList, o)

		// SaveData(outputList, "output.json")
	}

	format := req.FormValue("format")
	fmt.Println("Selected Format to output in: ", format)

	var filename string
	switch format {
	case "JSON":
		// Generate JSON
		fmt.Println("Generating Attendance List in JSON")
		SaveDataInJSON(outputList, "output.json")
		filename = "output.json"
		// Generate XML
	case "XML":
		fmt.Println("Generating Attendance List in XML")
		saveDataInXML(outputList, "output.xml")
		filename = "output.xml"

	default:
		fmt.Println("error in generateeee")
	}

	// Set the Content-Disposition header so the browser knows it's a downloadable file
	res.Header().Set("Content-Disposition", "attachment; filename="+filename)

	// Send the file to the client
	http.ServeFile(res, req, "internal/output/"+filename)
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

// The logoutPage function will handle all requests to the logout page
// If there is existing Cookie in the browser, the logout page will delete the cookie and redirect to the index page
// If there is no existing Cookie in the browser, the logout page will redirect to the index page
func LogoutPage(res http.ResponseWriter, req *http.Request) {
	myCookie, err := req.Cookie("userCookie")
	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	delete(mapUserSessions, myCookie.Value)
	delete(userList.CurSession, myCookie.Value)
	SaveData(userList, "data.json")
	myCookie.MaxAge = -1
	myCookie.Value = ""

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
			http.Error(res, "Invalid account", http.StatusForbidden)
			return
		}

		err := bcrypt.CompareHashAndPassword(userList.Users[username].Password, []byte(pwd))
		if err != nil {
			http.Error(res, "Invalid account", http.StatusForbidden)
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
		//return
	}
	tpl.ExecuteTemplate(res, "login.gohtml", nil)
}

// The admin function will handle all requests to the admin page
// The admin page will display all users and their information
func Admin(res http.ResponseWriter, req *http.Request) {
	curUser := getCurUser(res, req)
	adminData := AdminData{
		CurUser:  curUser,
		UserList: userList,
	}
	err := tpl.ExecuteTemplate(res, "admin.gohtml", adminData)
	if err != nil {
		fmt.Println("Error in Admin Page")
		fmt.Println(err)
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
				fmt.Println("Error in Admin Page")
				fmt.Println(err)
			}
		} else {
			err := tpl.ExecuteTemplate(res, "index.gohtml", curUser)
			if err != nil {
				fmt.Println("Error in Index Page")
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println("No User Logged In")
		err := tpl.ExecuteTemplate(res, "index.gohtml", nil)
		if err != nil {
			fmt.Println("Error in Login Page")
			fmt.Println(err)
		}
	}
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

// The signupPage will only be called by Admin to create new students
func SignupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		err := createNewUser(res, req)
		if err != nil {
			fmt.Println(err)
			return
		}
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(res, "signup.gohtml", nil)
}

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

// this function creates a brand new session and stored it in the client browser
func createSessionCookie(res http.ResponseWriter, username string) {
	id := uuid.NewV4()
	myCookie := &http.Cookie{
		Name:  "userCookie",
		Value: id.String(),
	}
	http.SetCookie(res, myCookie)
	mapUserSessions[myCookie.Value] = username
	userList.CurSession = mapUserSessions
	fmt.Println("Cookie Created")
}
