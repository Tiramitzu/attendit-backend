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

func GetAttendanceByUserAndDateAndCompany(userId primitive.ObjectID, date string, companyId primitive.ObjectID) (*db.Attendance, error) {
	attendance := &db.Attendance{}
	err := mgm.Coll(attendance).First(bson.M{"userId": userId, "date": date, "companyId": companyId}, attendance)
	if err != nil {
		return nil, err
	}

	return attendance, nil
}

func GetAttendancesByCompany(companyId primitive.ObjectID, page int) ([]*db.Attendance, error) {
	var attendances []*db.Attendance
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"companyId": companyId}
	opts := options.Find()
	opts.SetSkip(int64(page - 1))
	opts.SetLimit(10)

	cursor, err := mgm.Coll(&db.Attendance{}).Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.New("304: Not Modified")
	}

	for cursor.Next(ctx) {
		attendance := &db.Attendance{}
		err := cursor.Decode(attendance)
		if err != nil {
			return nil, errors.New("304: Not Modified")
		}

		attendances = append(attendances, attendance)
	}

	return attendances, nil
}
