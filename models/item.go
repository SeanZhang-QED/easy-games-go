package models

var ItemType = [...]string{"STREAM", "CLIP", "VIDEO"}

type Item struct {
	Id              string `json:"id" bson:"_id"`
	Title           string `json:"title" bson:"title"`
	Url             string `json:"url,omitempty" bson:"url"`
	ThumbnailUrl    string `json:"thumbnail_url" bson:"thumbnail_url"`
	BroadcasterName string `json:"broadcaster_name,omitempty" bson:"broadcaster_name"`
	UserName        string `json:"user_name,omitempty" bson:"user_name"`
	GameId          string `json:"game_id" bson:"game_id"`
	ItemType        string `json:"item_type" bson:"item_type"`
}

type TwitchItemResponse struct {
	Data       []Item     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Favorite struct {
	Item Item `json:"favorite"`
}
