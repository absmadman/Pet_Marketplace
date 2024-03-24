package db

import "time"

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Advert struct {
	Id       int       `json:"id"`
	UserId   int       `json:"user_id"`
	AdvertId int       `json:"advert_id"`
	Header   string    `json:"header"`
	Text     string    `json:"text"`
	ImageURL string    `json:"image_url"`
	Address  string    `json:"address"`
	Price    float64   `json:"price"`
	Datetime time.Time `json:"datetime"`
}

type AdvList struct {
	List []*Advert `json:"feed"`
}

type Filter struct {
	MinPrice   float64 `json:"min_price"`
	MaxPrice   float64 `json:"max_price"`
	FromNewest bool    `json:"from_newest"`
}
