package attendance

import (
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// this function creates a brand new session and stored it in the client browser
func createSessionCookie(res http.ResponseWriter, username string) {
	id := uuid.NewV4()

	claims := jwt.MapClaims{
		"username":  username,
		"sessionId": id.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println("Error signing token:", err)
		return
	}

	myCookie := &http.Cookie{
		Name:  "userCookie",
		Value: signedToken, //id.String(),
	}
	http.SetCookie(res, myCookie)
	mapUserSessions[myCookie.Value] = username
	userList.CurSession = mapUserSessions
	log.Println("Cookie Created")
}
