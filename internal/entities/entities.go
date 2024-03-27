package entities

import (
	"encoding/json"
	"github.com/go-passwd/validator"
	"time"
)

// User структура пользователя
type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Advert структура объявления
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

// AdvList структура которая хранит массив объявлений
type AdvList struct {
	List []*Advert `json:"feed"`
}

// Filter структура для филтрации и сортировки ленты объявлений
type Filter struct {
	MinPrice           float64 `json:"min_price"`
	MaxPrice           float64 `json:"max_price"`
	ByPrice            bool    `json:"by_price"`
	AscendingDirection bool    `json:"ascending_direction"`
}

// ValidateAdvertData валидатор для данных в объявлении
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

// MarshalBinary сериализует структуру Advert для Redis
func (a *Advert) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

// UnmarshalBinary десериализует массив byte в структуру Advert
func (a *Advert) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &a)
}

// NewUser конструктор для User
func NewUser(id int, login string, password string, token string) *User {
	return &User{
		Id:       id,
		Login:    login,
		Password: password,
	}
}
