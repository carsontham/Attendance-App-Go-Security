package attendance

import (
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

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
