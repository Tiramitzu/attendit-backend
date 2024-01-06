package redisServices

import (
	models "attendit/backend/models/db"
	"attendit/backend/services"
	"context"
	"errors"
	"github.com/go-redis/cache/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func getUserSchedulesCacheKey(userId primitive.ObjectID) string {
	return "req:cache:user:schedules:" + userId.Hex()
}

func CacheUserSchedules(userId primitive.ObjectID, schedule *[]models.Schedule) {
func CacheUserSchedules(userId primitive.ObjectID, schedule *[]models.Schedule, page int) {
	if !services.Config.UseRedis {
		return
	}

	userScheduleCacheKey := getUserSchedulesCacheKey(userId, page)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   userScheduleCacheKey,
		Value: schedule,
		TTL:   time.Second * 30,
	})
}

func GetUserSchedulesFromCache(userId primitive.ObjectID, page int) (*[]models.Schedule, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	var schedule []models.Schedule
	userScheduleCacheKey := getUserSchedulesCacheKey(userId, page)
	err := services.GetRedisCache().Get(context.TODO(), userScheduleCacheKey, &schedule)
	return &schedule, err
}
