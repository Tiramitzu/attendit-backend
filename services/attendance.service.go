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

func GetUserAttendances(userId primitive.ObjectID, page int) ([]db.Attendance, error) {
	var attendances []db.Attendance
	opts := options.Find()
	opts.SetLimit(25)
	opts.SetSkip(int64(page-1) * 25)
	opts.SetSort(bson.M{"createdAt": -1})
	err := mgm.Coll(&db.Attendance{}).SimpleFind(&attendances, bson.M{"userId": userId}, opts)

	if err != nil {
		return nil, err
	}

	return attendances, nil
}

func GetUserAttendancesByDate(userId primitive.ObjectID, fromDate string, toDate string, page int) ([]db.Attendance, error) {
	var attendances []db.Attendance
	opts := options.Find()
	opts.SetLimit(25)
	opts.SetSkip(int64(page-1) * 25)
	opts.SetSort(bson.M{"updatedAt": -1})
	err := mgm.Coll(&db.Attendance{}).SimpleFind(&attendances, bson.M{
		"userId": userId,
		"date": bson.M{
			"$gte": fromDate,
			"$lte": toDate,
		},
	}, opts)

	if err != nil {
		return nil, err
	}

	return attendances, nil
}

func GetAttendanceByUserAndDate(userId primitive.ObjectID, date string) (*db.Attendance, error) {
	attendance := &db.Attendance{}
	err := mgm.Coll(attendance).First(bson.M{"userId": userId, "date": date}, attendance)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}

		return nil, err
	}

	return attendance, nil
}

func GetAttendances(page int) ([]*db.Attendance, error) {
	var attendances []*db.Attendance
	var users []*db.User

	// Fetch attendances
	err := mgm.Coll(&db.Attendance{}).SimpleFind(&attendances, bson.M{}, options.Find().SetSkip(int64((page-1)*25)).SetLimit(25).SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return nil, err
	}

	// Fetch users
	err = mgm.Coll(&db.User{}).SimpleFind(&users, bson.M{})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	// Combine attendances and users
	for _, attendance := range attendances {
		for _, user := range users {
			if attendance.UserId == user.ID {
				attendance.User = user
				break
			}
		}
	}

	return attendances, nil
}

func GetAttendancesByDate(fromDate string, toDate string, page int) ([]db.Attendance, error) {
	var attendances []db.Attendance
	opts := options.Find()
	opts.SetLimit(25)
	opts.SetSkip(int64(page-1) * 25)
	opts.SetSort(bson.M{"createdAt": -1})
	err := mgm.Coll(&db.Attendance{}).SimpleFind(&attendances, bson.M{
		"date": bson.M{
			"$gte": fromDate,
			"$lte": toDate,
		},
	}, opts)

	if err != nil {
		return nil, err
	}

	return attendances, nil
}
