package db

const (
	insertUser                  = `INSERT INTO users (login, password) VALUES ($1, $2)`
	getUserByLogin              = `SELECT id, login, password FROM users WHERE login = $1`
	getUserByLoginAndPass       = `SELECT id, login, password FROM users WHERE login = $1 AND password = $2`
	insertAdvertWithIdReturn    = "INSERT INTO adverts (user_id, header, text, address, image_url, price, datetime) VALUES (%d, '%s', '%s', '%s', '%s', %f, '%v') RETURNING id"
	selectFromAdvertsDescending = `SELECT id, user_id, header, text, address, image_url, price, datetime FROM adverts WHERE price >= $1 AND price <= $2 ORDER BY datetime DESC`
	selectFromAdvertsAscending  = `SELECT id, user_id, header, text, address, image_url, price, datetime FROM adverts WHERE price >= $1 AND price <= $2 ORDER BY datetime`
	selectAdvertByAdvertId      = `SELECT id, user_id, header, text, address, image_url, price, datetime FROM adverts WHERE id = $1`
	AllowedSymbols              = `1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!_ /`
)
