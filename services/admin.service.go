package services

import (
	db "attendit/backend/models/db"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
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
