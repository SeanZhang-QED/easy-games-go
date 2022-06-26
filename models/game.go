package models

type Game struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	BoxArtUrl string `json:"box_art_url"`
}

type Pagination struct {
	Cursor string `json:"cursor"`
}

type TwitchGameResponse struct {
	Data       []Game     `json:"data"`
	Pagination Pagination `json:"pagination"`
}
