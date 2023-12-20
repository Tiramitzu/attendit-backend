package services

import (
	db "attendit/backend/models/db"
	"context"
	"github.com/kamva/mgm/v3/field"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindCompaniesByUserId(userId primitive.ObjectID) *[]db.Company {
	user, err := FindUserById(userId)

	if err != nil {
		return nil
	}

	companies := &[]db.Company{}
	for _, companyId := range user.Companies {
		company, _ := GetCompanyById(companyId)
		*companies = append(*companies, *company)
	}

	return companies
}

func GetCompanyById(companyId primitive.ObjectID) (*db.Company, error) {
	company := &db.Company{}
	err := mgm.Coll(company).First(bson.M{field.ID: companyId}, company)
	if err != nil {
		return nil, err
	}

	return company, nil
}

func FindMembersByCompanyId(companyId primitive.ObjectID, page int) *[]db.Member {
	members := &[]db.Member{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": companyId}
	opts := options.Find()
	opts.SetSkip(int64(page - 1))
	opts.SetLimit(10)

	cursor, _ := mgm.Coll(&db.Company{}).Aggregate(ctx, mongo.Pipeline{
		{{"$match", filter}},
		{{"$unwind", "$members"}},
		{{"$replaceRoot", bson.M{"newRoot": "$members"}}},
		{{"$skip", opts.Skip}},
		{{"$limit", opts.Limit}},
	})

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var member db.Member
		err := cursor.Decode(&member)
		if err != nil {
			return nil
		}
		*members = append(*members, member)
	}

	return members
}

func InsertMembersToCompany(companyId primitive.ObjectID, members []db.Member) error {
	company := &db.Company{}
	err := mgm.Coll(company).FindByID(companyId, company)
	if err != nil {
		return err
	}

	for _, member := range members {
		company.Members = append(company.Members, member)
	}

	err = mgm.Coll(company).Update(company)
	if err != nil {
		return err
	}

	return nil
}

func CreateCompany(company *db.Company) (*db.Company, error) {
	db.NewCompany(company.Author, company.Name, company.IPAddresses, company.CheckInTime, company.CheckOutTime, company.Members)
	err := mgm.Coll(company).Create(company)
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

func DeleteCompany(companyId primitive.ObjectID) error {
	company := &db.Company{}

	err := mgm.Coll(company).FindByID(companyId, company)
	if err != nil {
		return err
	}

	for _, member := range company.Members {
		user, _ := FindUserById(member.ID)
		for index, company := range user.Companies {
			if company == companyId {
				user.Companies = append(user.Companies[:index], user.Companies[index+1:]...)
			}
		}

		err = mgm.Coll(user).Update(user)
	}

	err = mgm.Coll(company).Delete(company)
	if err != nil {
		return err
	}

	return nil
}

func CreateInvitation(invitation *db.Invitation) (*db.Invitation, error) {
	db.NewInvitation(invitation.Author, invitation.UserID, invitation.CompanyID, invitation.Role)
	err := mgm.Coll(invitation).Create(invitation)

	if err != nil {
		return nil, err
	}

	return invitation, nil
}

func FindInvitationById(invitationId primitive.ObjectID) (*db.Invitation, error) {
	invitation := &db.Invitation{}
	err := mgm.Coll(invitation).FindByID(invitationId, invitation)
	if err != nil {
		return nil, err
	}

	return invitation, nil
}

func UpdateInvitation(invitation *db.Invitation) (*db.Invitation, error) {
	err := mgm.Coll(invitation).Update(invitation)
	if err != nil {
		return nil, err
	}

	return invitation, nil
}
