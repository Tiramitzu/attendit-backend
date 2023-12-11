package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Member struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	Email            string `json:"email" bson:"email"`
	Phone            string `json:"phone" bson:"phone"`
	Role             string `json:"role" bson:"role"`
}

type Company struct {
	mgm.DefaultModel `bson:",inline"`
	Author           primitive.ObjectID `json:"author" bson:"author"`
	Name             string             `json:"name" bson:"name"`
	IPAddresses      []string           `json:"ipAddresses" bson:"ipAddresses"`
	CheckInTime      string             `json:"checkInTime" bson:"checkInTime"`
	CheckOutTime     string             `json:"checkOutTime" bson:"checkOutTime"`
	Members          []Member           `json:"members" bson:"members"`
}

type CompanyClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Type  string `json:"type"`
}

func NewCompany(author primitive.ObjectID, name string, ipAddress []string, checkInTime string, checkOutTime string, members []Member) *Company {
	for index, member := range members {
		members[index].ID = author
		members[index].CreatedAt = time.Now()
		members[index].UpdatedAt = time.Now()
		members[index].Name = member.Name
		members[index].Email = member.Email
		members[index].Phone = member.Phone
		members[index].Role = member.Role
	}

	return &Company{
		Author:       author,
		Name:         name,
		IPAddresses:  ipAddress,
		CheckInTime:  checkInTime,
		CheckOutTime: checkOutTime,
		Members:      members,
	}
}

func (model *Company) CollectionName() string {
	return "companies"
}
