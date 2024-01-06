package services

import (
	db "attendit/backend/models/db"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAttendanceById(id primitive.ObjectID) (*db.Attendance, error) {
	attendance := &db.Attendance{}
	err := mgm.Coll(attendance).First(bson.M{"_id": id}, attendance)

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

func GetAttendanceByUserAndDate(userId primitive.ObjectID, date string) (*db.Attendance, error) {
	attendance := &db.Attendance{}
	err := mgm.Coll(attendance).First(bson.M{"userId": userId, "date": date}, attendance)
	if err != nil {
		return nil, err
	}

	return attendance, nil
}

func GetAttendancesByCompany(page int) ([]*db.Attendance, error) {
	var attendances []*db.Attendance
	opts := options.Find()
	opts.SetLimit(25)
	opts.SetSkip(int64((page - 1) * 25))
	err := mgm.Coll(&db.Attendance{}).SimpleFind(&attendances, bson.M{}, opts)

	if err != nil {
		return nil, err
	}

	return attendances, nil
}
