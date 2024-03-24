package db

import (
	"database/sql"
	"errors"
	"github.com/go-passwd/validator"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func NewUser(id int, login string, password string, token string) *User {
	return &User{
		Id:       id,
		Login:    login,
		Password: password,
	}
}

/*
func NewAdvert() *Advert{
	return &Advert
}
*/

func CheckIdExist(db *sql.DB, login string) bool {
	rows, err := db.Exec(getUserByLogin, login)
	if err != nil {
		return true
	}
	num, err := rows.RowsAffected()
	if err != nil {
		return true
	}
	if num == 0 {
		return false
	}
	return true
}

func (u *User) CreateUser(db *sql.DB) error {
	if CheckIdExist(db, u.Login) {
		return errors.New("already exist")
	}
	_, err := db.Exec(insertUser, u.Login, u.Password)
	if err != nil {
		return err
	}
	err = db.QueryRow(getUserByLogin, u.Login).Scan(&u.Id, &u.Login, &u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetUser(db *sql.DB) error {
	err := db.QueryRow(getUserByLoginAndPass, u.Login, u.Password).Scan(&u.Id, &u.Login, &u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (a *Advert) CreateAdvert(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO adverts (user_id, header, text, address, image_url, price, datetime) VALUES ($1, $2, $3, $4, $5, $6, $7)", a.UserId, a.Header, a.Text, a.Address, a.ImageURL, a.Price, time.Now())
	log.Println(err)
	if err != nil {
		return err
	}
	return nil
}

func (a *Advert) ValidateAdvertData() bool {
	textValidator := validator.New(validator.MinLength(8, nil), validator.MaxLength(512, nil), validator.ContainsOnly("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!_ /", nil))
	otherValidator := validator.New(validator.MinLength(8, nil), validator.MaxLength(64, nil), validator.ContainsOnly("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!_ /", nil))
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
	return true
}

func (f *Filter) GetRows(db *sql.DB) (*sql.Rows, error) {
	if !f.FromNewest {
		rows, err := db.Query("SELECT id, user_id, header, text, address, image_url, price, datetime FROM adverts WHERE price >= $1 AND price <= $2 ORDER BY datetime", f.MinPrice, f.MaxPrice)
		if err != nil {
			return nil, err
		}
		return rows, nil
	} else {
		rows, err := db.Query("SELECT id, user_id, header, text, address, image_url, price, datetime FROM adverts WHERE price >= $1 AND price <= $2 ORDER BY datetime DESC", f.MinPrice, f.MaxPrice)
		if err != nil {
			return nil, err
		}
		return rows, nil
	}
}

func (al *AdvList) GetAdvList(page int, db *sql.DB, filter *Filter) error {
	rows, err := filter.GetRows(db)
	if err != nil {
		return err
	}
	i := 0
	stop := (page + 1) * 10
	for rows.Next() {
		var adv Advert
		err = rows.Scan(&adv.Id, &adv.UserId, &adv.Header, &adv.Text, &adv.Address, &adv.ImageURL, &adv.Price, &adv.Datetime)
		if err != nil {
			break
		}
		if i >= stop-10 && i < stop {
			al.List = append(al.List, &adv)
		}
	}
	return nil
}

func NewDBConnection() *sql.DB {
	//pgPass := os.Getenv("POSTGRES_PASSWORD")
	//pgUser := os.Getenv("POSTGRES_USER")
	//pgDb := os.Getenv("POSTGRES_DB")
	//connStr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", pgUser, pgPass, pgDb)
	connStr := "postgres://server:server@localhost:5432/api_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
