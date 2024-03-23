package db

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
)

func NewUser(id int, login string, password string, token string) *User {
	return &User{
		Id:       id,
		Login:    login,
		Password: password,
	}
}

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
