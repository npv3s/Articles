package authenticator

import (
	"bytes"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"main/database"
	"strconv"
)

type sessions map[int]hash

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

func (a *Authenticator) NewUser(login, password string) (*int, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, errors.New("Hash gen error: " + err.Error())
	}
	return a.database.NewUser(login, string(hashBytes))
}

func (a *Authenticator) Login(login, password string) (*int, error) {
	user, err := a.database.GetUserByLogin(login)
	if err != nil {
		return nil, errors.New("User " + login + " is not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password))
	if err != nil {
		return nil, errors.New("Wrong password for user " + login)
	}

	return &user.Id, nil
}

func (a *Authenticator) GetSession(userId int) (string, error) {
	token, ok := a.sessions[userId]
	if ok {
		return string(token), nil
	}

	passwd, err := a.database.GetUserById(userId)
	if err != nil {
		return "", errors.New("No user with id " + strconv.Itoa(userId))
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(passwd.PassHash), 14)
	if err != nil {
		return "", errors.New("Hash gen error: " + err.Error())
	}

	a.sessions[userId] = hashBytes

	return string(hashBytes), nil
}

func (a *Authenticator) CheckSession(userId int, token string) bool {
	rightHash, ok := a.sessions[userId]
	if ok {
		return bytes.Equal(rightHash, []byte(token))
	}

	return false
}
