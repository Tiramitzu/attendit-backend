package services

import (
	db "attendit/backend/models/db"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindAttendanceById(id primitive.ObjectID) (*db.Attendance, error) {
	attendance := &db.Attendance{}
	err := mgm.Coll(attendance).FindByID(id, attendance)

	if err != nil {
		return nil, err
	}

	return attendance, nil
}

func AttendanceCheckIn(attendance *db.Attendance) (*db.Attendance, error) {
	err := mgm.Coll(attendance).Create(attendance)

	if err != nil {
		return nil, err
	}

	return attendance, nil
}

func AttendanceCheckOut(attendance *db.Attendance) (*db.Attendance, error) {
	err := mgm.Coll(attendance).Update(attendance)

	if err != nil {
		return nil, err
	}

	return attendance, nil
}
