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
	// Step 2: Encrypt the (user + password => unique) as password before save to the database
	bsEmail := []byte(user.Email)
	bsPassword := []byte(user.Password)
	bs, err := bcrypt.GenerateFromPassword(append(bsEmail, bsPassword...), bcrypt.MinCost) //... is required, because append() is a variadic function that accepts an unlimited number of arguments.
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
		http.Error(w, "User already exists", http.StatusConflict)
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
	w.Header().Set("Content-Type", "application/json")
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
	firstName, err := uh.verifyUser(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, "Error occured in verifying User's password", http.StatusUnauthorized)
		fmt.Printf("Error occured in verifying User's password\n")
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
	// return user's first name(prepared for frontend)
	var loginResp models.LoginResponse
	loginResp.Email = credentials.Email
	loginResp.Name = firstName
	loginRespJson, err := json.Marshal(loginResp) 
	if err != nil {
		http.Error(w, "Fail to create response for login success.", http.StatusInternalServerError)
		return
	}
	w.Write(loginRespJson)
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

	ck, err := r.Cookie("sessionId")
	if err != nil {
		http.Error(w, "Haven't logged in", http.StatusUnauthorized)
		fmt.Println("Fail to read cookie from http request")
		return 
	}
	
	loggedSession, err := session.AlreadyLoggedIn(w, r, uh.session)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	} 
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

// MongoDB Connections -------------------

func (uh UserHandler) addUser(u *models.User) (bool, error) {
	// check user existence
	var users []models.User
	if err := uh.session.DB("easy-games-db").C("users").FindId(u.Email).All(&users); err != nil {
		fmt.Printf("Failed to check user existence from MongoDB %v\n", err)
		return false, err
	}
	if len(users) != 0 {
		return false, nil
	}

	// store the user in mongodb
	err := uh.session.DB("easy-games-db").C("users").Insert(u)
	if err != nil {
		fmt.Printf("Failed to insert user to MongoDB %v\n", err)
		return false, err
	}
	return true, nil
}

func (uh UserHandler) verifyUser(email string, password string) (string, error) {
	// get the user 
	u, err := getUserByEmail(email, uh.session)
	if ; err != nil {
		return "", err
	}

	// does the entered password match the stored password?
	bsEmail := []byte(email)
	bsPassword := []byte(password)
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), append(bsEmail, bsPassword...))
	if err != nil {
		fmt.Println("Wrong Password.")
		return "", err
	}
	return u.FirstName, nil
}

func getUserByEmail(email string, session *mgo.Session) (models.User, error) {
	var u models.User
	err := session.DB("easy-games-db").C("users").FindId(email).One(&u)
	if err != nil {
		fmt.Println("Fail to fetch the user info from users document")
		return models.User{}, err
	}
	return u, nil
}