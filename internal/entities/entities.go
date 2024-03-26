package entities

import (
	"encoding/json"
	"github.com/go-passwd/validator"
	"time"
)

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Advert struct {
	Id         int       `json:"id"`
	UserId     int       `json:"user_id"`
	Header     string    `json:"header"`
	Text       string    `json:"text"`
	ImageURL   string    `json:"image_url"`
	Address    string    `json:"address"`
	Price      float64   `json:"price"`
	Datetime   time.Time `json:"datetime"`
	ByThisUser bool      `json:"by_this_user"`
}

type AdvList struct {
	List []*Advert `json:"feed"`
}

type Filter struct {
	MinPrice   float64 `json:"min_price"`
	MaxPrice   float64 `json:"max_price"`
	FromNewest bool    `json:"from_newest"`
}

func (a *Advert) ValidateAdvertData() bool {
	textValidator := validator.New(
		validator.MinLength(8, nil),
		validator.MaxLength(512, nil),
		validator.ContainsOnly(AllowedSymbols, nil))
	otherValidator := validator.New(
		validator.MinLength(8, nil),
		validator.MaxLength(64, nil),
		validator.ContainsOnly(AllowedSymbols, nil))
	err := textValidator.Validate(a.Text)
	if err != nil {
		return false
	}
	err = otherValidator.Validate(a.ImageURL)
	if err != nil {
		return false
	}
	err = otherValidator.Validate(a.Address)
	if err != nil {
		return false
	}
	err = otherValidator.Validate(a.Header)
	if err != nil {
		return false
	}
	if a.Price < 1 || a.Price > 1000000 {
		return false
	}
	return true
}

func (a *Advert) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Advert) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &a)
}

func NewUser(id int, login string, password string, token string) *User {
	return &User{
		Id:       id,
		Login:    login,
		Password: password,
	}
}
