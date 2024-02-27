package services

import (
	db "attendit/backend/models/db"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetFeedbacks(userId primitive.ObjectID, isAdmin bool, page int) ([]*db.Feedback, error) {
	var feedbacks []*db.Feedback
	var users []*db.UserWithoutProfPic

	if isAdmin {
		err := mgm.Coll(&db.Feedback{}).SimpleFind(&feedbacks, bson.M{}, options.Find().SetSkip(int64((page-1)*25)).SetLimit(25).SetSort(bson.M{"created_at": -1}))
		if err != nil {
			return nil, err
		}
		err = mgm.Coll(&db.User{}).SimpleFind(&users, bson.M{})
		if err != nil {
			return nil, err
		}

		for _, feedback := range feedbacks {
			for _, user := range users {
				if feedback.From == user.ID {
					feedback.User = user
					break
				}
			}
		}
	} else {
		err := mgm.Coll(&db.Feedback{}).SimpleFind(&feedbacks, bson.M{"userId": userId}, options.Find().SetSkip(int64((page-1)*25)).SetLimit(25).SetSort(bson.M{"created_at": -1}))
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

func UpdateFeedbackStatus(id primitive.ObjectID, status string) (*db.Feedback, error) {
	feedback := &db.Feedback{}
	err := mgm.Coll(feedback).FindByID(id, feedback)

	if err != nil {
		return nil, err
	}

	feedback.Status = status
	err = mgm.Coll(feedback).Update(feedback)

	if err != nil {
		return nil, err
	}

	return feedback, nil
}
