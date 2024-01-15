package services

import (
	models "attendit/backend/models/db"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUserSchedules(userId primitive.ObjectID, page int) (*[]models.Schedule, error) {
	var schedules []models.Schedule
	opts := options.Find()
	opts.SetLimit(25)
	opts.SetSkip(int64((page - 1) * 25))
	opts.SetSort(bson.M{"createdAt": -1})
	err := mgm.Coll(&models.Schedule{}).SimpleFind(&schedules, bson.M{"userId": userId}, opts)

	if err != nil {
		return nil, err
	}

	return &schedules, nil
}

func GetScheduleById(id primitive.ObjectID) (*models.Schedule, error) {
	schedule := &models.Schedule{}
	err := mgm.Coll(schedule).FindByID(id, schedule)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func CreateSchedule(schedule *models.Schedule) (*models.Schedule, error) {
	err := mgm.Coll(schedule).Create(schedule)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}
