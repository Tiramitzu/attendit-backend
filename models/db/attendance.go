package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attendance struct {
	mgm.DefaultModel `bson:",inline"`
	UserId           primitive.ObjectID `json:"userId" bson:"userId"`
	IpAddress        string             `json:"ipAddress" bson:"ipAddress"`
	Status           string             `json:"status" bson:"status"`
	Date             string             `json:"date" bson:"date"`
	CheckIn          string             `json:"checkIn" bson:"checkIn"`
	CheckOut         string             `json:"checkOut" bson:"checkOut"`
}

func NewAttendance(userId primitive.ObjectID, ipAddress string, date string, status string, checkIn string, checkOut string) *Attendance {
	return &Attendance{
		UserId:    userId,
		IpAddress: ipAddress,
		Status:    status,
		Date:      date,
		CheckIn:   checkIn,
		CheckOut:  checkOut,
	}
}

func (model *Attendance) CollectionName() string {
	return "attendances"
}
