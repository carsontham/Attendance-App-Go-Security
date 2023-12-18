package attendance

import (
	"html/template"
	"log"
	"net/http"
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
	log.Println("Server start up...")
	LoadData()
}

// this function will read the chosen radio button and generate the attendance list in the chosen format
// the file will be returned to the client via serveFiles()
func Generate(res http.ResponseWriter, req *http.Request) {
	var outputList []Output

	for user, _ := range userList.Users {
		if user == "admin" {
			continue
		}
		o := Output{
			Username:  userList.Users[user].Username,
			FirstName: userList.Users[user].FirstName,
			LastName:  userList.Users[user].LastName,
			CheckedIn: userList.Users[user].CheckedIn,
			TimeIn:    userList.Users[user].TimeIn,
		}
		outputList = append(outputList, o)
	}

	format := req.FormValue("format")
	log.Println("Selected Format to output in: ", format)

	var filename string
	switch format {
	case "JSON":
		// Generate JSON
		log.Println("Generating Attendance List in JSON")
		SaveDataInJSON(outputList, "output.json")
		filename = "output.json"
		// Generate XML
	case "XML":
		log.Println("Generating Attendance List in XML")
		saveDataInXML(outputList, "output.xml")
		filename = "output.xml"

	default:
		log.Println("error when generating ouput file")
	}

	// Set the Content-Disposition header so the browser knows it's a downloadable file
	res.Header().Set("Content-Disposition", "attachment; filename="+filename)

	// Send the file to the client
	http.ServeFile(res, req, "internal/output/"+filename)
}
