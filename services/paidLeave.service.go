package services

import (
	db "attendit/backend/models/db"
	"errors"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetActiveRequest(userId primitive.ObjectID) (*db.PaidLeave, error) {
	var paidLeave *db.PaidLeave
	err := mgm.Coll(paidLeave).First(bson.M{"userId": userId, "accepted": false}, paidLeave)

	if err != nil {
		return nil, errors.New("Gagal mencari permintaan cuti yang belum selesai")
	}

	return paidLeave, nil
}

func CreatePaidLeave(userId primitive.ObjectID, reason string) (*db.PaidLeave, error) {
	paidLeave := db.NewPaidLeave(userId, false, primitive.NilObjectID, reason)
	err := mgm.Coll(paidLeave).Create(paidLeave)

	if err != nil {
		return nil, errors.New("Gagal membuat permintaan cuti")
	}

	return paidLeave, nil
}
