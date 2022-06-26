package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id     	bson.ObjectId `json:"id" bson:"_id"` // automatically assigned to a document upon insert
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
