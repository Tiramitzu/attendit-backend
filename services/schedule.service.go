package services

import (
	models "attendit/backend/models/db"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserSchedules(userId primitive.ObjectID) (*[]models.Schedule, error) {
	var schedules []models.Schedule
	err := mgm.Coll(&models.Schedule{}).SimpleFind(&schedules, bson.M{"userId": userId})
	if err != nil {
		return nil, err
	}

	return &schedules, nil
}

func CreateSchedule(schedule *models.Schedule) (*models.Schedule, error) {
	err := mgm.Coll(schedule).Create(schedule)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}
