package session

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SeanZhang-QED/easy-games-go/models"
)

const MAX_AGE = 1 * 60 * 60
// var Users = map[string]models.User{}       // user ID, user
var Sessions = map[string]models.Session{} // session ID, session
var LastCleaned time.Time = time.Now()

func AlreadyLoggedIn(w http.ResponseWriter, req *http.Request) (bool, string) {
	ck, err := req.Cookie("sessionId")
	if err != nil {
		return false, ""
	}
	s, ok := Sessions[ck.Value]
	if !ok {
		return false, ""		
	}
	// refresh session
	s.LastActivity = time.Now()
	Sessions[ck.Value] = s
	ck.MaxAge = MAX_AGE
	http.SetCookie(w, ck)
	return true, s.Email
}

// for demonstration purposes
func Show() {
	fmt.Println("***sessionID***")
	for k, v := range Sessions {
		fmt.Println(k, v.Email)
	}
	fmt.Println("***sessionID***")
}
