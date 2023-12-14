package services

import (
	db "attendit/backend/models/db"
	"errors"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser create a user record
func CreateUser(username string, email string, plainPassword string, displayName string, phone string) (*db.User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("cannot generate hashed password")
	}

	user := db.NewUser(email, string(password), username, displayName, phone)
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

// FindUserById find user by id
func FindUserById(userId primitive.ObjectID) (*db.User, error) {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(userId, user)
	if err != nil {
		return nil, errors.New("304: Not Modified")
	}

	return user, nil
}

// FindUserByEmail find user by email
func FindUserByEmail(email string) (*db.User, error) {
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

func FindUserAttendanceByCompany(companyId primitive.ObjectID, userId primitive.ObjectID) (*[]db.Attendance, error) {
	attendances := &[]db.Attendance{}
	attendance := &db.Attendance{}
	_ = mgm.Coll(attendance).SimpleFind(bson.M{"companyId": companyId, "userId": userId}, attendances)

	return attendances, nil
}

func FindAttendanceByCompany(companyId primitive.ObjectID) (*db.Attendance, error) {
	attendance := &db.Attendance{}
	err := mgm.Coll(attendance).First(bson.M{"companyId": companyId}, attendance)
	if err != nil {
		return nil, errors.New("304: Not Modified")
	}

	return attendance, nil
}
