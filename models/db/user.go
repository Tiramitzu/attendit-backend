package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string               `json:"email" bson:"email"`
	Password         string               `json:"-" bson:"password"`
	DisplayName      string               `json:"displayName" bson:"displayName"`
	Phone            string               `json:"phone" bson:"phone"`
	UserName         string               `json:"username" bson:"username"`
	MailVerified     bool                 `json:"mail_verified" bson:"mail_verified"`
	Companies        []primitive.ObjectID `json:"companies" bson:"companies"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Type  string `json:"type"`
}

func NewUser(email string, password string, username string, displayName string, phone string) *User {
	return &User{
		Email:        email,
		Password:     password,
		DisplayName:  displayName,
		UserName:     username,
		Phone:        phone,
		MailVerified: false,
	}
}

func (model *User) CollectionName() string {
	return "users"
}

// You can override Collection functions or CRUD hooks
// https://github.com/Kamva/mgm#a-models-hooks
// https://github.com/Kamva/mgm#collections
