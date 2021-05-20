package authenticator

import (
	"bytes"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"main/database"
	"math/rand"
	"strconv"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

type sessions map[int]hash

type hash []byte

type Authenticator struct {
	database interface{ database.Database }
	sessions sessions
	admins []int
}

func NewAuthenticator(database database.Database) (*Authenticator, error) {
	admins, err := database.GetAdmins()
	if err != nil {
		return nil, err
	}

	return &Authenticator{
		database,
		sessions{},
		admins,
	}, nil
}

func (a *Authenticator) GenPassword(password string) ([]byte, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, errors.New("Hash gen error: " + err.Error())
	}

	return hashBytes, nil
}

func (a *Authenticator) NewUser(login, password string) (*int, error) {
	hashBytes, err := a.GenPassword(password)
	if err != nil {
		return nil, err
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

	_, err := a.database.GetUserById(userId)
	if err != nil {
		return "", errors.New("No user with id " + strconv.Itoa(userId))
	}

	hashBytes, err := bcrypt.GenerateFromPassword(randStringBytes(16), 5)
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

func (a *Authenticator) IsAdmin(userId int) bool {
	for _, uId := range a.admins {
		if uId == userId {
			return true
		}
	}
	return false
}
