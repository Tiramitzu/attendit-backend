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
	Accepted         bool               `json:"accepted" bson:"accepted"`
	AcceptedBy       primitive.ObjectID `json:"acceptedBy" bson:"acceptedBy"`
}

func NewPaidLeave(userId primitive.ObjectID, accepted bool, acceptedBy primitive.ObjectID, reason string, startDate string, days int) *PaidLeave {
	return &PaidLeave{
		UserId:     userId,
		Reason:     reason,
		StartDate:  startDate,
		Days:       days,
		Accepted:   accepted,
		AcceptedBy: acceptedBy,
	}
}

func (model *PaidLeave) CollectionName() string {
	return "paidLeaves"
}
