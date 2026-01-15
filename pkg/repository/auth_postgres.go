package repository

import (
	"fmt"
	"restapi-app"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db}
}

func (p *AuthPostgres) CreateUser(user restapi.User) (int, error) {
	tr, err := p.db.Begin()
	if err != nil {
		return -1, err
	}

	defer func() {
		if err != nil {
			tr.Rollback()
		}
	}()

	var id int
	CreateUserQuery := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)
	row := tr.QueryRow(CreateUserQuery, user.Name, user.Username, user.Password)

	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id, tr.Commit()
}

func (p *AuthPostgres) GetUser(username, password_hash string) (restapi.User, error) {
	var user restapi.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := p.db.Get(&user, query, username, password_hash)

	return user, err
}
