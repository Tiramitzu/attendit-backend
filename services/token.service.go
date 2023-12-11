package services

import (
	"errors"
	"time"

	db "attendit/backend/models/db"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateToken create a new token record
func CreateToken(user *db.User) (*db.Token, error) {
	claims := &db.UserClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Subject:  user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(Config.JWTSecretKey))
	if err != nil {
		return nil, errors.New("cannot create access token")
	}

	tokenModel := db.NewToken(user.ID, tokenString)
	err = mgm.Coll(tokenModel).Create(tokenModel)
	if err != nil {
		return nil, errors.New("cannot save access token to db")
	}

	return tokenModel, nil
}

// DeleteTokenById delete token with id
func DeleteTokenById(tokenId primitive.ObjectID) error {
	ctx := mgm.Ctx()
	deleteResult, err := mgm.Coll(&db.Token{}).DeleteOne(ctx, bson.M{field.ID: tokenId})
	if err != nil || deleteResult.DeletedCount <= 0 {
		return errors.New("cannot delete token")
	}

	return nil
}

func GetTokenById(userId primitive.ObjectID) (*db.Token, error) {
	tokenModel := &db.Token{}
	err := mgm.Coll(tokenModel).First(
		bson.M{"user": userId},
		tokenModel,
	)
	if err != nil {
		return nil, errors.New("401: Unauthorized")
	}

	return tokenModel, nil
}

// GenerateAccessTokens generates "access" and "refresh" token for user
func GenerateAccessTokens(user *db.User) (*db.Token, error) {
	accessToken, err := CreateToken(user)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func VerifyToken(token string) (*db.Token, error) {
	claims := &db.UserClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Config.JWTSecretKey), nil
	})

	if err != nil {
		return nil, errors.New("401: Unauthorized")
	}

	tokenModel := &db.Token{}
	userId, _ := primitive.ObjectIDFromHex(claims.Subject)
	err = mgm.Coll(tokenModel).First(
		bson.M{"user": userId},
		tokenModel,
	)
	if err != nil {
		return nil, errors.New("401: Unauthorized")
	}

	return tokenModel, nil
}
