package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/SeanZhang-QED/easy-games-go/models"
	"gopkg.in/mgo.v2"
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
    // Step 2:Data format validation - Also can be done in front-end
    if user.Email == "" || user.Password == "" || regexp.MustCompile(`^[a-zA-Z0-9]$`).MatchString(user.Email) {
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        fmt.Printf("Invalid username or password\n")
        return
    }
	// Step 3: Call addUser() to insert the user into the database
    success, err := addUser(&user) // note: pass by pointer
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
	exists, err := checkUser(credentials.Email, credentials.Password)
    if err != nil {
        http.Error(w, "Failed to read user from MongoDB", http.StatusInternalServerError)
        fmt.Printf("Failed to read user from MongoDB %v\n", err)
        return
    }
    // fail to sign in
    if !exists {
        http.Error(w, "User doesn't exists or wrong password", http.StatusUnauthorized)
        fmt.Printf("User doesn't exists or wrong password\n")
        return
    }
    // Step 2: handle seesion cookie
    // Step 3: return
	w.WriteHeader(http.StatusOK) // 200
	fmt.Printf("User Login successfully: %s.\n", credentials.Email) 
}

func addUser(user *models.User) (bool, error){
	return true, nil
}

func checkUser(email string, password string) (bool, error) {
	return true, nil
}
