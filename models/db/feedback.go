package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feedback struct {
	mgm.DefaultModel `bson:",inline"`
	From             primitive.ObjectID `json:"from" bson:"from"`
	Content          string             `json:"content" bson:"content"`
	Status           string             `json:"status" bson:"status"`
}

func NewFeedback(from primitive.ObjectID, content string) *Feedback {
	return &Feedback{
		From:    from,
		Content: content,
		Status:  "pending",
	}
}

func (model *Feedback) CollectionName() string {
	return "feedbacks"
}
