package services

import (
	db "attendit/backend/models/db"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUsers(page int) ([]*db.User, error) {
	var users []*db.User
	opts := options.Find()
	opts.SetLimit(25)
	opts.SetSkip(int64((page - 1) * 25))
	err := mgm.Coll(&db.User{}).SimpleFind(&users, bson.M{}, opts)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func VerifyUser(userId primitive.ObjectID) (*db.User, error) {
	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}

	user.IsVerified = true
	err = mgm.Coll(&db.User{}).Update(user)

	if err != nil {
		return user, err
	}

	return user, nil
}

func GetUnVerifiedUsers() ([]*db.User, error) {
	var users []*db.User
	err := mgm.Coll(&db.User{}).SimpleFind(&users, bson.M{"isVerified": false})

	if err != nil {
		return nil, err
	}

	return users, nil
}

func DeleteUnVerifiedUsers() ([]*db.User, error) {
	var users []*db.User
	err := mgm.Coll(&db.User{}).SimpleFind(&users, bson.M{"isVerified": false})

	if err != nil {
		return nil, err
	}

	for _, user := range users {
		err = mgm.Coll(&db.User{}).Delete(user)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}
