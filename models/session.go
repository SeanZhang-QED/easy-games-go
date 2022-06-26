package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Session struct {
	Id        bson.ObjectId `bson:"_id"`
	SessionId string        `bson:"session_id"`
	Email     string        `bson:"email"`
}
