package redisServices

import (
	db "attendit/backend/models/db"
	"attendit/backend/services"
	"context"
	"github.com/go-redis/cache/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func getFeedbacksCacheKey(userId primitive.ObjectID, admin bool) string {
	if admin {
		return "req:cache:admin:feedbacks:admin"
	}
	return "req:cache:user:feedbacks:" + userId.Hex()
}

func CacheFeedbacks(userId primitive.ObjectID, messages []*db.Feedback, admin bool) {
	if !services.Config.UseRedis {
		return
	}

	messagesCacheKey := getFeedbacksCacheKey(userId, admin)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   messagesCacheKey,
		Value: messages,
		TTL:   time.Minute,
	})
}

func GetFeedbacksFromCache(userId primitive.ObjectID, admin bool) ([]*db.Feedback, error) {
	if !services.Config.UseRedis {
		return nil, nil
	}

	var messages []*db.Feedback
	messagesCacheKey := getFeedbacksCacheKey(userId, admin)
	err := services.GetRedisCache().Get(context.TODO(), messagesCacheKey, &messages)
	return messages, err
}
