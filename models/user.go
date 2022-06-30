package models

type User struct {
	Email     string        `json:"email" bson:"_id"`
	Password  string        `json:"password" bson:"password"`
	FirstName string        `json:"first_name" bson:"first_name"`
	LastName  string        `json:"last_name" bson:"last_name"`
	FavoriteRecords map[string]Void `bson:"favorite_records"` 
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Email    string `json:"email"`
	Name	 string `json:"name"`
}

type Void struct{} // map[]Void works as a set