package main

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)


var sessions = map[string]session{}

type session struct {
	Username string
	expiry time .Time
}

type sessionSend struct {
	Username string `json:username`
	Token string	`json:token`
	Expiry time .Time
}


func HashPassword(password string) (string) {
    bytes, err:= bcrypt.GenerateFromPassword([]byte(password), 14)
    if(err != nil){
    }
    return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

