package models

import (
	"github.com/kamva/mgm/v3"
)

type Company struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string   `json:"name" bson:"name"`
	IPAddresses      []string `json:"ipAddresses" bson:"ipAddresses"`
	Locations        []string `json:"locations" bson:"locations"`
	CheckInTime      string   `json:"checkInTime" bson:"checkInTime"`
	CheckOutTime     string   `json:"checkOutTime" bson:"checkOutTime"`
}

func (model *Company) CollectionName() string {
	return "companies"
}
