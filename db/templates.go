package db

const (
	intertUser = `INSERT INTO users (login, password) VALUES ($1, $2)`
)
