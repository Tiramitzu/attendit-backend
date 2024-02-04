package services

import (
	db "attendit/backend/models/db"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func GetCompany() (*db.Company, error) {
	company := &db.Company{}
	err := mgm.Coll(company).First(bson.M{}, company)
	if err != nil {
		return nil, err
	}

	return company, nil
}

func UpdateCompany(company *db.Company) (*db.Company, error) {
	err := mgm.Coll(company).Update(company)
	if err != nil {
		return nil, err
	}

	return company, nil
}
