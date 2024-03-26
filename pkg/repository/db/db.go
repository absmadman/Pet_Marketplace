package db

import (
	"VK_Internship_Marketplace/internal/entities"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
	"time"
)

type Database struct {
	connection *sql.DB
}

func (db *Database) UpdateAdvert(a *entities.Advert) error {
	_, err := db.connection.Exec(updateAdvert, a.Header, a.Text, a.Address, a.ImageURL, a.Price, time.Now(), a.Id)
	if err != nil {
		return err
	}
	a.ByThisUser = true
	return nil
}

func (db *Database) RemoveAdvert(a *entities.Advert) error {
	if _, err := db.connection.Exec("DELETE FROM adverts WHERE id = $1", a.Id); err != nil {
		return err
	}
	return nil
}

func (db *Database) CheckUserIdExist(login string) bool {
	rows, err := db.connection.Exec(getUserByLogin, login)
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

func (db *Database) CreateUser(u *entities.User) error {
	if db.CheckUserIdExist(u.Login) {
		return errors.New("already exist")
	}
	_, err := db.connection.Exec(insertUser, u.Login, u.Password)
	if err != nil {
		return err
	}
	err = db.connection.QueryRow(getUserByLogin, u.Login).Scan(&u.Id, &u.Login, &u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetUser(u *entities.User) error {
	err := db.connection.QueryRow(getUserByLoginAndPass, u.Login, u.Password).Scan(&u.Id, &u.Login, &u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) CreateAdvert(a *entities.Advert) error {
	splitTime := strings.Split(time.Now().String(), " ")
	formattedTime := splitTime[0] + " " + splitTime[1]
	err := db.connection.QueryRow(fmt.Sprintf(insertAdvertWithIdReturn,
		a.UserId, a.Header, a.Text, a.Address, a.ImageURL, a.Price, formattedTime)).Scan(&a.Id)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetRows(f *entities.Filter) (*sql.Rows, error) {
	if !f.FromNewest {
		rows, err := db.connection.Query(selectFromAdvertsAscending, f.MinPrice, f.MaxPrice)
		if err != nil {
			return nil, err
		}
		return rows, nil
	} else {
		rows, err := db.connection.Query(selectFromAdvertsDescending, f.MinPrice, f.MaxPrice)
		if err != nil {
			return nil, err
		}
		return rows, nil
	}
}

func (db *Database) GetAdvList(page int, al *entities.AdvList, filter *entities.Filter, authUserId int) error {
	rows, err := db.GetRows(filter)
	if err != nil {
		return err
	}
	i := 0
	stop := (page + 1) * 10
	for rows.Next() {
		var adv entities.Advert
		err = rows.Scan(&adv.Id, &adv.UserId, &adv.Header, &adv.Text, &adv.Address, &adv.ImageURL, &adv.Price, &adv.Datetime)
		adv.ByThisUser = false
		if authUserId == adv.UserId {
			adv.ByThisUser = true
		}
		if err != nil {
			break
		}
		if i >= stop-10 && i < stop {
			al.List = append(al.List, &adv)
		}
		i++
	}
	return nil
}

func (db *Database) GetAdv(a *entities.Advert, advId int) error {
	err := db.connection.QueryRow(selectAdvertByAdvertId, advId).
		Scan(
			&a.Id,
			&a.UserId,
			&a.Header,
			&a.Text,
			&a.Address,
			&a.ImageURL,
			&a.Price,
			&a.Datetime)
	if err != nil {
		return err
	}
	return nil
}

func NewDatabase() *Database {
	return &Database{
		connection: NewDBConnection(),
	}
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
