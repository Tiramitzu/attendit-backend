package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schedule struct {
	mgm.DefaultModel `bson:",inline"`
	UserId           primitive.ObjectID `json:"userId" bson:"userId"`
	Title            string             `json:"date" bson:"date"`
	StartTime        string             `json:"startTime" bson:"startTime"`
	EndTime          string             `json:"endTime" bson:"endTime"`
}

func NewSchedule(userId primitive.ObjectID, title string, startTime string, endTime string, status string) *Schedule {
	return &Schedule{
		UserId:    userId,
		Title:     title,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

func (model *Schedule) CollectionName() string {
	return "schedules"
}
