package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invitation struct {
	mgm.DefaultModel `bson:",inline"`
	Author           primitive.ObjectID `json:"author" bson:"author"`
	UserID           primitive.ObjectID `json:"userId" bson:"userId"`
	CompanyID        primitive.ObjectID `json:"companyId" bson:"companyId"`
	Role             string             `json:"role" bson:"role"`
	Status           string             `json:"status" bson:"status"`
}

func NewInvitation(author primitive.ObjectID, userId primitive.ObjectID, companyId primitive.ObjectID, role string) *Invitation {
	return &Invitation{
		Author:    author,
		UserID:    userId,
		CompanyID: companyId,
		Role:      role,
		Status:    "pending",
	}
}
