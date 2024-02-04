package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaidLeave struct {
	mgm.DefaultModel `bson:",inline"`
	UserId           primitive.ObjectID `json:"userId" bson:"userId"`
	Reason           string             `json:"reason" bson:"reason"`
	StartDate        string             `json:"startDate" bson:"startDate"`
	Days             int                `json:"days" bson:"days"`
	Status           int                `json:"status" bson:"status"`
	StatusBy         primitive.ObjectID `json:"statusBy" bson:"statusBy"`
	User             *User              `json:"user" bson:"user"`
}

func NewPaidLeave(userId primitive.ObjectID, status int, statusBy primitive.ObjectID, reason string, startDate string, days int) *PaidLeave {
	return &PaidLeave{
		UserId:    userId,
		Reason:    reason,
		StartDate: startDate,
		Days:      days,
		Status:    status,
		StatusBy:  statusBy,
	}
}

func (model *PaidLeave) CollectionName() string {
	return "paidLeaves"
}
