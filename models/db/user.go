package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/kamva/mgm/v3"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string `json:"email" bson:"email"`
	Password         string `json:"-" bson:"password"`
	FullName         string `json:"fullName" bson:"fullName"`
	Phone            string `json:"phone" bson:"phone"`
	AccessLevel      int    `json:"accessLevel" bson:"accessLevel"`
	IsVerified       bool   `json:"isVerified" bson:"isVerified"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Type  string `json:"type"`
}

func NewUser(email string, password string, fullName string, phone string) *User {
	return &User{
		Email:       email,
		Password:    password,
		FullName:    fullName,
		Phone:       phone,
		AccessLevel: 0,
		IsVerified:  false,
	}
}

func (model *User) CollectionName() string {
	return "users"
}
