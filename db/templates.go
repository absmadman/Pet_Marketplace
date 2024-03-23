package db

const (
	insertUser            = `INSERT INTO users (login, password) VALUES ($1, $2)`
	getUserByLogin        = `SELECT id, login, password FROM users WHERE login = $1`
	getUserByLoginAndPass = `SELECT id, login, password FROM users WHERE login = $1 AND password = $2`
)
