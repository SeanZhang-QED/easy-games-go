package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "time"

	"github.com/SeanZhang-QED/easy-games-go/models"
	"github.com/SeanZhang-QED/easy-games-go/session"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// added session to our userController
type UserHandler struct {
	session *mgo.Session
}

// added session to our userController
func NewUserHandler(s *mgo.Session) *UserHandler {
	return &UserHandler{s}
}

func (uh UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signup request")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/plain")

	if r.Method == http.MethodOptions {
		return
	}

	// Step 1: Decode username and possword from http request Json body
	decoder := json.NewDecoder(r.Body)
	var user models.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode user data from client %v\n", err)
		return
	}
	// Step 2: password encrypt
	bs, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
		fmt.Printf("Failed to encrypt password\n")
		return
	}
	user.Password = string(bs)

	// Step 3: Call addUser() to insert the user into the database
	success, err := uh.addUser(&user) // note: pass by pointer
	if err != nil {
		http.Error(w, "Failed to save user to MongoDB", http.StatusInternalServerError)
		fmt.Printf("Failed to save user to MongoDB %v\n", err)
		return
	}

	if !success {
		http.Error(w, "User already exists", http.StatusBadRequest)
		fmt.Println("User already exists")
		return
	}
	w.WriteHeader(http.StatusCreated) // 201 Created
	fmt.Printf("User added successfully: %s.\n", user.Email)
}

func (uh UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one Login request")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	w.Header().Set("Content-Type", "text/plain")
	if r.Method == http.MethodOptions {
		return
	}

	// Step 1: check user name and possword
	decoder := json.NewDecoder(r.Body)
	var credentials models.Credentials
	if err := decoder.Decode(&credentials); err != nil {
		http.Error(w, "Cannot decode Login credentials from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode Login credentials from client %v\n", err)
		return
	}
	exists, err := uh.verifyUser(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, "Failed to read user from MongoDB", http.StatusInternalServerError)
		fmt.Printf("Failed to read user from MongoDB %v\n", err)
		return
	}
	if !exists {
		http.Error(w, "User doesn't exists or wrong password", http.StatusUnauthorized)
		fmt.Printf("User doesn't exists or wrong password\n")
		return
	}

	// Step 2: handle seesion cookie
	sID := uuid.NewV4()
	ck := &http.Cookie{
		Name:   "sessionId",
		Value:  sID.String(),
		MaxAge: session.MAX_AGE,
	}
	http.SetCookie(w, ck)
	// step 3: manage session collection on mongoDB
	ss, err := session.SearchSessionByEmail(uh.session, credentials.Email)
	if err != nil {
		http.Error(w, "Failed to fetch session from MongoDB", http.StatusInternalServerError)
		fmt.Println("Fail to fetch session from mongoDB")
		return
	}
	if len(ss) == 0 {
		// insert
		err := session.InsertSessionBySId(uh.session, credentials.Email, sID.String())
		if err != nil {
			http.Error(w, "Fail to insert session from mongoDB", http.StatusInternalServerError)
			return
		}
	} else {
		// update
		err := session.UpdateSessionById(uh.session, ss[0].Id, sID.String())
		if err != nil {
			http.Error(w, "Fail to insert session from mongoDB", http.StatusInternalServerError)
			return
		}
	}
	// Step 4: return
	w.WriteHeader(http.StatusOK) // 200
	fmt.Printf("User Login successfully: %s.\n", credentials.Email)
}

func (uh UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one Log out request")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	w.Header().Set("Content-Type", "text/plain")
	if r.Method == http.MethodOptions {
		return
	}

	loggedSession, err := session.AlreadyLoggedIn(w, r, uh.session)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	} else if loggedSession == (models.Session{}) {
		http.Error(w, "Haven't logged in", http.StatusBadRequest)
		return
	}

	ck, _ := r.Cookie("sessionId")
	// delete the session
	session.DeleteSessionBySId(uh.session, ck.Value)
	// remove the cookie
	ck = &http.Cookie{
		Name:   "sessionId",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, ck)
	w.WriteHeader(http.StatusOK) // 200
	fmt.Printf("User Logout successfully: %s.\n", loggedSession.Email)
}

func (uh UserHandler) addUser(u *models.User) (bool, error) {
	// create bson ID
	u.Id = bson.NewObjectId()

	// store the user in mongodb
	err := uh.session.DB("easy-games-db").C("users").Insert(u)
	if err != nil {
		fmt.Printf("Failed to insert user to MongoDB %v\n", err)
		return false, err
	}
	return true, nil
}

func (uh UserHandler) verifyUser(email string, password string) (bool, error) {
	// composite literal
	var u models.User

	// Fetch user
	if err := uh.session.DB("easy-games-db").C("users").Find(bson.M{"email": email}).One(&u); err != nil {
		fmt.Printf("Failed to fetch user from MongoDB %v\n", err)
		return false, err
	}

	// does the entered password match the stored password?
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		fmt.Printf("Wrong Password.")
		return false, err
	}
	return true, nil
}
