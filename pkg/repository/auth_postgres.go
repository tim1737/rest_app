package repository

import (
	"fmt"

	todo "github.com/Tim-Masuda/rest_todo"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}



// sing up
func (r *AuthPostgres) CreateUser(user todo.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable) // добавление в таблицу, формирование запроса

	row := r.db.QueryRow(query, user.Name, user.Username, user.Password) // добавление данных в query

	if err := row.Scan(&id); err != nil { // чек, запись id 
		return 0, err
	}	

	return id, nil
}

// sing in
func (r *AuthPostgres) GetUser(username, password string) (todo.User, error) {
	var user todo.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable) // формирование запроса
	err := r.db.Get(&user, query, username, password) // записываем в user подставоляем query в запрос

	return user, err
} 