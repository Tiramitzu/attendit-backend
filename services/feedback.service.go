package services

import (
	db "attendit/backend/models/db"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetFeedbacks(userId primitive.ObjectID, isAdmin bool) ([]*db.Feedback, error) {
	var feedbacks []*db.Feedback
	if isAdmin {
		err := mgm.Coll(&db.Feedback{}).SimpleFind(&feedbacks, bson.M{})
		if err != nil {
			return nil, err
		}
	} else {
		err := mgm.Coll(&db.Feedback{}).SimpleFind(&feedbacks, bson.M{"userId": userId})
		if err != nil {
			return nil, err
		}
	}

	return feedbacks, nil
}

func GetTotalFeedbacks() (int64, error) {
	total, err := mgm.Coll(&db.Feedback{}).CountDocuments(mgm.Ctx(), bson.M{})

	if err != nil {
		return 0, err
	}

	return total, nil
}

func SendFeedback(feedback *db.Feedback) (*db.Feedback, error) {
	err := mgm.Coll(feedback).Create(feedback)

	if err != nil {
		return nil, err
	}

	return feedback, nil
}
