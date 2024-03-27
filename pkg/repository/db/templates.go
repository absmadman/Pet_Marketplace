package db

const (
	insertUser                         = `INSERT INTO users (login, password) VALUES ($1, $2)`
	getUserByLogin                     = `SELECT id, login, password FROM users WHERE login = $1`
	getUserByLoginAndPass              = `SELECT id, login, password FROM users WHERE login = $1 AND password = $2`
	insertAdvertWithIdReturn           = "INSERT INTO adverts (user_id, header, text, address, image_url, price, datetime) VALUES (%d, '%s', '%s', '%s', '%s', %f, '%v') RETURNING id"
	selectFromAdvertsByDateDescending  = `SELECT id, user_id, header, text, address, image_url, price, datetime FROM adverts WHERE price >= $1 AND price <= $2 ORDER BY datetime DESC`
	selectFromAdvertsByDateAscending   = `SELECT id, user_id, header, text, address, image_url, price, datetime FROM adverts WHERE price >= $1 AND price <= $2 ORDER BY datetime`
	selectFromAdvertsByPriceDescending = `SELECT id, user_id, header, text, address, image_url, price, datetime FROM adverts WHERE price >= $1 AND price <= $2 ORDER BY price DESC`
	selectFromAdvertsByPriceAscending  = `SELECT id, user_id, header, text, address, image_url, price, datetime FROM adverts WHERE price >= $1 AND price <= $2 ORDER BY price`
	selectAdvertByAdvertId             = `SELECT id, user_id, header, text, address, image_url, price, datetime FROM adverts WHERE id = $1`
	updateAdvert                       = `UPDATE adverts SET header = $1, text = $2, address = $3, image_url = $4, price = $5, datetime = $6 WHERE id = $7`
)
