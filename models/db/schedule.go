package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schedule struct {
	mgm.DefaultModel `bson:",inline"`
	UserId           primitive.ObjectID `json:"userId" bson:"userId"`
	Title            string             `json:"title" bson:"title"`
	StartTime        string             `json:"startTime" bson:"startTime"`
	EndTime          string             `json:"endTime" bson:"endTime"`
	Date             string             `json:"date" bson:"date"`
}

func NewSchedule(userId primitive.ObjectID, title string, startTime string, endTime string, date string) *Schedule {
	return &Schedule{
		UserId:    userId,
		Title:     title,
		StartTime: startTime,
		EndTime:   endTime,
		Date:      date,
	}
}

func (model *Schedule) CollectionName() string {
	return "schedules"
}
