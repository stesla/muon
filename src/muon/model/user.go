package model

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"errors"
)

type User struct {
	Id int64
	Login string
	Email string
	pwhash []byte
}

var InvalidUserOrPass = errors.New("invalid username or password")

const pwCost = bcrypt.DefaultCost

func createUsersTable() error {
	sql :=`CREATE IF NOT EXISTS users (
		   id        INTEGER      PRIMARY KEY,
		   login     VARCHAR(32)  UNIQUE NOT NULL,
		   email     VARCHAR(255) UNIQUE NOT NULL,
		   pwhash    CHAR(60)     NOT NULL);`
	_, err := db.Exec(sql)
	return err
}

func CreateUser(login, pw, email string) (*User, error) {
	var err error

	user := &User{Login: login, Email: email}
	user.pwhash, err = bcrypt.GenerateFromPassword([]byte(pw), pwCost)

	r, err := db.Exec(
		`INSERT INTO users (login, email, pwhash)
		 VALUES (?, ?, ?)`,
		user.Login, user.Email, user.pwhash)
	if err != nil {
		return nil, err
	}

	user.Id, err = r.LastInsertId()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindUser(login string) (*User, error) {
	row := db.QueryRow(
		`SELECT id, login, email, pwhash
		 FROM users WHERE login = ?`, login)
	user := &User{}
	err := row.Scan(&user.Id, &user.Login, &user.Email, &user.pwhash)
	if err == sql.ErrNoRows {
		return nil, NotFound
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func AuthUser(login, pw string) (*User, error) {
	user, err := FindUser(login)
	if err == NotFound {
		return nil, InvalidUserOrPass
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.pwhash, []byte(pw))
	if err != nil {
		return nil, InvalidUserOrPass
	}

	return user, nil
}
