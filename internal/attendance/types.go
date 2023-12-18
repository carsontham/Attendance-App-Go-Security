package attendance

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
