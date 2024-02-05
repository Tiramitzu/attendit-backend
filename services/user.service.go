package services

import (
	db "attendit/backend/models/db"
	"errors"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"net"
	"net/http"
	"strings"
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

// GetUserByToken find user by token
func GetUserByToken(token string) (*db.User, error) {
	user := &db.User{}
	tkn := &db.Token{}
	err := mgm.Coll(tkn).First(bson.M{"token": token}, tkn)
	if err != nil {
		return nil, errors.New("304: Not Modified")
	}

	err = mgm.Coll(user).First(bson.M{"_id": tkn.User}, user)
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
		return nil, errors.New("Email dan password tidak cocok")
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

func GetUsers(page int) ([]*db.UserWithPassword, error) {
	var users []*db.UserWithPassword
	opts := options.Find()
	opts.SetLimit(25)
	opts.SetSkip(int64((page - 1) * 25))
	err := mgm.Coll(&db.UserWithPassword{}).SimpleFind(&users, bson.M{}, opts)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetClientIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("IP not found")
}
