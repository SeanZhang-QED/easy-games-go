package models

import "gopkg.in/mgo.v2/bson"

type FavoriteRecords struct {
	Id 		bson.ObjectId `bson:"_id"`
	UserId string `bson:"email"`
	ItemId string `bson:"item_id"`
}