package services

import (
	db "attendit/backend/models/db"
	"errors"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser create a user record
func CreateUser(email string, plainPassword string, fullName string, phone string) (*db.User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("cannot generate hashed password")
	}

	user := db.NewUser(email, string(password), fullName, phone)
	err = mgm.Coll(user).Create(user)
	if err != nil {
		return nil, errors.New("cannot create new user")
	}

	return user, nil
}

func UpdateUser(user *db.User) (*db.User, error) {
	err := mgm.Coll(user).Update(user)
	if err != nil {
		return nil, errors.New("304: Not Modified")
	}

	return user, nil
}

// GetUserById find user by id
func GetUserById(userId primitive.ObjectID) (*db.User, error) {
	user := &db.User{}
	err := mgm.Coll(user).First(bson.M{"_id": userId}, user)
	if err != nil {
		return nil, errors.New("304: Not Modified")
	}

	return user, nil
}

// GetUserByEmail find user by email
func GetUserByEmail(email string) (*db.User, error) {
	user := &db.User{}
	err := mgm.Coll(user).First(bson.M{"email": email}, user)
	if err != nil {
		return nil, errors.New("Email and password don't match")
	}

	return user, nil
}

// CheckUserMail search user by email, return error if someone uses
func CheckUserMail(email string) error {
	user := &db.User{}
	userCollection := mgm.Coll(user)
	err := userCollection.First(bson.M{"email": email}, user)
	if err == nil {
		return errors.New("409: Conflict")
	}

	return nil
}

func GetUserAttendances(userId primitive.ObjectID, page int) (*[]db.Attendance, error) {
	var attendances *[]db.Attendance
	opts := options.Find()
	opts.SetLimit(25)
	opts.SetSkip(int64(page-1) * 25)
	err := mgm.Coll(&db.Attendance{}).SimpleFind(&attendances, bson.M{"userId": userId}, opts)

	if err != nil {
		return nil, err
	}

	return attendances, nil
}
