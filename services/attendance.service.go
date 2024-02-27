package services

import (
	"attendit/backend/models"
	db "attendit/backend/models/db"
	"encoding/base64"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

func CheckOutAllAttendances() error {
	company := &db.Company{}
	err := mgm.Coll(company).First(bson.M{}, company)
	if err != nil {
		return err
	}
	_, err = mgm.Coll(&db.Attendance{}).UpdateMany(mgm.Ctx(), bson.M{"checkOut": ""}, bson.M{"$set": bson.M{"checkOut": company.CheckOutTime}})

	if err != nil {
		return err
	}

	fmt.Println("All attendances checked out")

	return nil
}

func GetUserAttendances(userId primitive.ObjectID, page int) ([]db.Attendance, error) {
	var attendances []db.Attendance
	opts := options.Find()
	opts.SetLimit(25)
	opts.SetSkip(int64(page-1) * 25)
	opts.SetSort(bson.M{"created_at": -1})
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
	opts.SetSort(bson.M{"created_at": -1})
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

func GetUserTotalAttendances(userId primitive.ObjectID) (models.AttendanceTotal, error) {
	location, _ := time.LoadLocation("Asia/Jakarta")

	totalAll, err := mgm.Coll(&db.Attendance{}).CountDocuments(mgm.Ctx(), bson.M{"userId": userId})
	if err != nil {
		return models.AttendanceTotal{}, err
	}

	day := time.Now().Day()
	month := time.Now().Month()
	year, week := time.Now().ISOWeek()

	startWeek := WeekStart(year, week).Format("02-01-2006")
	startWeekDay, _ := strconv.Atoi(startWeek[:2])
	businessWeekDay := startWeekDay

	businessWeekDays := 0
	for i := startWeekDay; i <= day; i++ {
		Day := time.Date(year, month, i, 0, 0, 0, 0, location)
		if Day.Weekday() != time.Saturday && Day.Weekday() != time.Sunday {
			if Day.Format("02-01-2006") <= time.Now().Format("02-01-2006") {
				businessWeekDay++
				businessWeekDays++
			}
		}
	}

	t := time.Date(year, month, 32, 0, 0, 0, 0, location)
	daysInMonth := 32 - t.Day()
	businessDays := 0
	for i := 1; i <= daysInMonth; i++ {
		Day := time.Date(year, month, i, 0, 0, 0, 0, location)
		if Day.Weekday() != time.Saturday && Day.Weekday() != time.Sunday {
			if Day.Format("02-01-2006") <= time.Now().Format("02-01-2006") {
				businessDays++
			}
		}
	}

	totalPresentToday, err := mgm.Coll(&db.Attendance{}).CountDocuments(mgm.Ctx(), bson.M{
		"userId":     userId,
		"created_at": primitive.NewDateTimeFromTime(time.Date(year, month, day, 0, 0, 0, 0, location)),
	})

	totalPresentWeek, err := mgm.Coll(&db.Attendance{}).CountDocuments(mgm.Ctx(), bson.M{
		"userId": userId,
		"created_at": bson.M{
			"$gte": primitive.NewDateTimeFromTime(time.Date(year, month, startWeekDay, 0, 0, 0, 0, location)),
			"$lte": primitive.NewDateTimeFromTime(time.Date(year, month, businessWeekDay, 23, 59, 59, 1e9-1, location)),
		},
	})
	if err != nil {
		return models.AttendanceTotal{}, err
	}

	totalPresentMonth, err := mgm.Coll(&db.Attendance{}).CountDocuments(mgm.Ctx(), bson.M{
		"userId": userId,
		"created_at": bson.M{
			"$gte": primitive.NewDateTimeFromTime(time.Date(year, month, 1, 0, 0, 0, 0, location)),
			"$lte": primitive.NewDateTimeFromTime(time.Date(year, month, day, 23, 59, 59, 1e9-1, location)),
		},
	})
	if err != nil {
		return models.AttendanceTotal{}, err
	}

	return models.AttendanceTotal{
		All: totalAll,
		Today: models.AttendanceWM{
			Present: totalPresentToday,
			Absent:  totalAll - totalPresentToday,
		},
		Weekly: models.AttendanceWM{
			Present: totalPresentWeek,
			Absent:  (totalAll * int64(businessWeekDays)) - totalPresentWeek,
		},
		Monthly: models.AttendanceWM{
			Present: totalPresentMonth,
			Absent:  (totalAll * int64(businessDays)) - totalPresentMonth,
		},
	}, nil
}

func GetTotalAttendances() (models.AttendanceTotal, error) {
	location, _ := time.LoadLocation("Asia/Jakarta")

	totalAll, err := mgm.Coll(&db.Attendance{}).CountDocuments(mgm.Ctx(), bson.M{})
	if err != nil {
		return models.AttendanceTotal{}, err
	}

	day := time.Now().Day()
	month := time.Now().Month()
	year, week := time.Now().ISOWeek()

	startWeek := WeekStart(year, week).Format("02-01-2006")
	startWeekDay, _ := strconv.Atoi(startWeek[:2])
	businessWeekDay := startWeekDay

	businessWeekDays := 0
	for i := startWeekDay; i <= day; i++ {
		Day := time.Date(year, month, i, 0, 0, 0, 0, location)
		if Day.Weekday() != time.Saturday && Day.Weekday() != time.Sunday {
			if Day.Format("02-01-2006") <= time.Now().Format("02-01-2006") {
				businessWeekDay++
				businessWeekDays++
			}
		}
	}

	t := time.Date(year, month, 32, 0, 0, 0, 0, location)
	daysInMonth := 32 - t.Day()
	businessDays := 0
	for i := 1; i <= daysInMonth; i++ {
		Day := time.Date(year, month, i, 0, 0, 0, 0, location)
		if Day.Weekday() != time.Saturday && Day.Weekday() != time.Sunday {
			if Day.Format("02-01-2006") <= time.Now().Format("02-01-2006") {
				businessDays++
			}
		}
	}

	totalUsers, err := mgm.Coll(&db.User{}).CountDocuments(mgm.Ctx(), bson.M{
		"accessLevel": 0,
	})
	if err != nil {
		return models.AttendanceTotal{}, err
	}

	totalPresentToday, err := mgm.Coll(&db.Attendance{}).CountDocuments(mgm.Ctx(), bson.M{
		"created_at": primitive.NewDateTimeFromTime(time.Date(year, month, day, 0, 0, 0, 0, location)),
	})

	totalPresentWeek, err := mgm.Coll(&db.Attendance{}).CountDocuments(mgm.Ctx(), bson.M{
		"created_at": bson.M{
			"$gte": primitive.NewDateTimeFromTime(time.Date(year, month, startWeekDay, 0, 0, 0, 0, location)),
			"$lte": primitive.NewDateTimeFromTime(time.Date(year, month, businessWeekDay, 23, 59, 59, 1e9-1, location)),
		},
	})
	if err != nil {
		return models.AttendanceTotal{}, err
	}

	totalPresentMonth, err := mgm.Coll(&db.Attendance{}).CountDocuments(mgm.Ctx(), bson.M{
		"created_at": bson.M{
			"$gte": primitive.NewDateTimeFromTime(time.Date(year, month, 1, 0, 0, 0, 0, location)),
			"$lte": primitive.NewDateTimeFromTime(time.Date(year, month, day, 23, 59, 59, 1e9-1, location)),
		},
	})
	if err != nil {
		return models.AttendanceTotal{}, err
	}

	return models.AttendanceTotal{
		All: totalAll,
		Today: models.AttendanceWM{
			Present: totalPresentToday,
			Absent:  totalUsers - totalPresentToday,
		},
		Weekly: models.AttendanceWM{
			Present: totalPresentWeek,
			Absent:  (totalUsers * int64(businessWeekDays)) - totalPresentWeek,
		},
		Monthly: models.AttendanceWM{
			Present: totalPresentMonth,
			Absent:  (totalUsers * int64(businessDays)) - totalPresentMonth,
		},
	}, nil
}

func GetAttendances(page int) ([]*db.Attendance, error) {
	var attendances []*db.Attendance
	var users []*db.UserWithoutProfPic

	// Fetch attendances
	err := mgm.Coll(&db.Attendance{}).SimpleFind(&attendances, bson.M{}, options.Find().SetSkip(int64((page-1)*25)).SetLimit(25).SetSort(bson.M{"created_at": -1}))
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

func GetTotalAttendancesByDate(fromDate string, toDate string) (models.AttendanceTotal, error) {
	totalAll, err := mgm.Coll(&db.Attendance{}).CountDocuments(mgm.Ctx(), bson.M{
		"date": bson.M{
			"$gte": fromDate,
			"$lte": toDate,
		},
	})
	if err != nil {
		return models.AttendanceTotal{}, err
	}

	totalUsers, err := mgm.Coll(&db.User{}).CountDocuments(mgm.Ctx(), bson.M{
		"accessLevel": 0,
	})
	if err != nil {
		return models.AttendanceTotal{}, err
	}

	from, _ := time.Parse("02-01-2006", fromDate)
	to, _ := time.Parse("02-01-2006", toDate)
	days := int(to.Sub(from).Hours() / 24)

	totalAbsent := totalUsers*int64(days) - totalAll

	return models.AttendanceTotal{
		All: totalAll,
		Weekly: models.AttendanceWM{
			Present: totalAll,
			Absent:  totalAbsent,
		},
	}, nil
}

func GetAttendancesByDate(fromDate string, toDate string, page int) ([]*db.Attendance, error) {
	var attendances []*db.Attendance
	opts := options.Find()
	opts.SetLimit(25)
	opts.SetSkip(int64(page-1) * 25)
	opts.SetSort(bson.M{"created_at": -1})
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

func WeekStart(year, week int) time.Time {
	location, _ := time.LoadLocation("Asia/Jakarta")

	// Start from the first day of the year:
	t := time.Date(year, 1, 1, 0, 0, 0, 0, location)

	// Roll forward to Monday of the first week:
	if wd := t.Weekday(); wd != time.Monday {
		daysToAdd := time.Monday - wd
		if daysToAdd < 0 {
			daysToAdd += 7 // If it's before Monday, roll forward to the next Monday
		}
		t = t.AddDate(0, 0, int(daysToAdd))
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

func SaveImage(base string, user *db.User, imagesDir string, types string) (string, error) {
	// Assume user.Photo contains base64-encoded JPEG data
	// Decode the base64 string to obtain the image data
	photoData := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base))

	// Decode the JPEG data into an image
	pht, err := jpeg.Decode(photoData)
	if err != nil {
		log.Println("Error decoding JPEG:", err)
		return "", err
	}

	// Create the directory if it doesn't exist
	err = os.MkdirAll(imagesDir, os.ModePerm)
	if err != nil {
		log.Println("Error creating directory:", err)
		return "", err
	}

	photoPath := filepath.Join(imagesDir, "user-"+types+"-"+user.ID.Hex()+".jpg")
	file, err := os.Create(photoPath)
	if err != nil {
		log.Println("Error creating image file:", err)
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	// Write the image data to the file
	err = jpeg.Encode(file, pht, nil)
	if err != nil {
		log.Println("Error encoding image to JPEG:", err)
		return "", err
	}

	return photoPath, nil
}
