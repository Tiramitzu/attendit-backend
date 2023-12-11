package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attendance struct {
	mgm.DefaultModel `bson:",inline"`
	UserId           primitive.ObjectID `json:"userId" bson:"userId"`
	CompanyId        primitive.ObjectID `json:"companyId" bson:"companyId"`
	IpAddress        string             `json:"ipAddress" bson:"ipAddress"`
	Date             string             `json:"date" bson:"date"`
}

func NewAttendance(userId primitive.ObjectID, companyId primitive.ObjectID, ipAddress string, date string) *Attendance {
	return &Attendance{
		UserId:    userId,
		CompanyId: companyId,
		IpAddress: ipAddress,
		Date:      date,
	}
}

func (model *Attendance) CollectionName() string {
	return "attendances"
}
