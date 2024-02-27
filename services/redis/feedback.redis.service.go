package redisServices

import (
	db "attendit/backend/models/db"
	"attendit/backend/services"
	"context"
	"github.com/go-redis/cache/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

func getFeedbacksCacheKey(userId primitive.ObjectID, admin bool, page int) string {
	if admin {
		return "req:cache:admin:feedbacks:admin"
	}
	return "req:cache:user:feedbacks:" + userId.Hex() + ":" + strconv.Itoa(page)
}

func CacheFeedbacks(userId primitive.ObjectID, messages []*db.Feedback, admin bool, page int) {
	if !services.Config.UseRedis {
		return
	}

	messagesCacheKey := getFeedbacksCacheKey(userId, admin, page)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   messagesCacheKey,
		Value: messages,
		TTL:   time.Minute,
	})
}

func GetFeedbacksFromCache(userId primitive.ObjectID, admin bool, page int) ([]*db.Feedback, error) {
	if !services.Config.UseRedis {
		return nil, nil
	}

	var messages []*db.Feedback
	messagesCacheKey := getFeedbacksCacheKey(userId, admin, page)
	err := services.GetRedisCache().Get(context.TODO(), messagesCacheKey, &messages)
	return messages, err
}
