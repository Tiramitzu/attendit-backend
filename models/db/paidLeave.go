package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaidLeave struct {
	mgm.DefaultModel `bson:",inline"`
	UserId           primitive.ObjectID `json:"userId" bson:"userId"`
	Accepted         bool               `json:"accepted" bson:"accepted"`
	AcceptedBy       primitive.ObjectID `json:"acceptedBy" bson:"acceptedBy"`
	Reason           string             `json:"reason" bson:"reason"`
}

func NewPaidLeave(userId primitive.ObjectID, accepted bool, acceptedBy primitive.ObjectID, reason string) *PaidLeave {
	return &PaidLeave{
		UserId:     userId,
		Accepted:   accepted,
		AcceptedBy: acceptedBy,
		Reason:     reason,
	}
}

func (model *PaidLeave) CollectionName() string {
	return "paidLeaves"
}
