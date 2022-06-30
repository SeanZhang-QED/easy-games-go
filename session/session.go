package session

import (
	"fmt"
	"github.com/SeanZhang-QED/easy-games-go/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

const MAX_AGE = 15 * 60

func AlreadyLoggedIn(w http.ResponseWriter, req *http.Request, ms *mgo.Session) (models.Session, error) {
	ck, _ := req.Cookie("sessionId")
	var s models.Session
	err := ms.DB("easy-games-db").C("sessions").Find(bson.M{"session_id": ck.Value}).One(&s)

	if err != nil {
		fmt.Println("Fail to fetch session from mongoDB")
		return models.Session{}, err
	}

	ck.MaxAge = MAX_AGE
	http.SetCookie(w, ck)
	return s, nil
}

func SearchSessionByEmail(ms *mgo.Session, email string) ([]models.Session, error) {
	var ss []models.Session
	err := ms.DB("easy-games-db").C("sessions").Find(bson.M{"email": email}).All(&ss)

	if err != nil {
		fmt.Println("Fail to fetch session from mongoDB")
		return nil, err
	}
	return ss, nil
}

func InsertSessionBySId(ms *mgo.Session, email string, sId string) error {
	var s models.Session
	s.Id = bson.NewObjectId()
	s.Email = email
	s.SessionId = sId
	if err := ms.DB("easy-games-db").C("sessions").Insert(s); err != nil {
		fmt.Println("Fail to insert session from mongoDB")
		return err
	}
	return nil
}

func UpdateSessionById(ms *mgo.Session, oid bson.ObjectId, sId string) error {

	if err := ms.DB("easy-games-db").C("sessions").UpdateId(oid, bson.M{"$set": bson.M{"session_id": sId}}); err != nil {
		fmt.Println("Fail to update session from mongoDB")
		return err
	}

	return nil
}

func DeleteSessionBySId(ms *mgo.Session, sId string) error {
	if _, err := ms.DB("easy-games-db").C("sessions").RemoveAll(bson.M{"session_id": sId}); err != nil {
		fmt.Println("Fail to delete session from mongoDB")
		return err
	}
	return nil
}
