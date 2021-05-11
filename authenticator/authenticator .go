package authenticator

import (
	"bytes"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"main/database"
)

type sessions map[string]hash

type hash []byte

type Authenticator struct {
	database interface{ database.Database }
	sessions sessions
}

func NewAuthenticator(database database.Database) Authenticator {
	return Authenticator{
		database,
		sessions{},
	}
}

func (a *Authenticator) NewUser(login, password string) error {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return errors.New("Hash gen error: " + err.Error())
	}
	return a.database.NewUser(login, string(hashBytes))
}

func (a *Authenticator) Login(login, password string) error {
	user, err := a.database.GetPassword(login)
	if err != nil {
		return errors.New("User " + login + " is not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.PassHash))
	if err != nil {
		return errors.New("Wrong password for user " + login)
	}

	return nil
}

func (a *Authenticator) GetSession(login string) (string, error) {
	token, ok := a.sessions[login]
	if ok {
		return string(token), nil
	}

	passwd, err := a.database.GetPassword(login)
	if err != nil {
		return "", errors.New("No such user " + login)
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(passwd.PassHash), 14)
	if err != nil {
		return "", errors.New("Hash gen error: " + err.Error())
	}

	return string(hashBytes), nil
}

func (a *Authenticator) CheckSession(login, token string) bool {
	rightHash, ok := a.sessions[login]
	if ok {
		return bytes.Equal(rightHash, []byte(token))
	}

	return false
}
