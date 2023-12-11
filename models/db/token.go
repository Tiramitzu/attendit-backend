package models

import (
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	mgm.DefaultModel `bson:",inline"`
	User             primitive.ObjectID `json:"user" bson:"user"`
	Token            string             `json:"token" bson:"token"`
}

func (model *Token) GetResponseJson() gin.H {
	return gin.H{"token": model.Token}
}

func (model *Token) GetResponseString() string {
	return model.Token
}

func NewToken(userId primitive.ObjectID, tokenString string) *Token {
	return &Token{
		User:  userId,
		Token: tokenString,
	}
}

func (model *Token) CollectionName() string {
	return "tokens"
}

// You can override Collection functions or CRUD hooks
// https://github.com/Kamva/mgm#a-models-hooks
// https://github.com/Kamva/mgm#collections
