package redisServices

import (
	models "attendit/backend/models/db"
	"attendit/backend/services"
	"context"
	"errors"
	"github.com/go-redis/cache/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

func getUserScheduleCacheKey(scheduleId primitive.ObjectID) string {
	return "req:cache:schedule:" + scheduleId.Hex()
}

func CacheSchedule(schedule *models.Schedule) {
	if !services.Config.UseRedis {
		return
	}

	scheduleCacheKey := getUserScheduleCacheKey(schedule.ID)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   scheduleCacheKey,
		Value: schedule,
		TTL:   time.Minute,
	})
}

func GetScheduleFromCache(scheduleId primitive.ObjectID) (*models.Schedule, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	var schedule models.Schedule
	scheduleCacheKey := getUserScheduleCacheKey(scheduleId)
	err := services.GetRedisCache().Get(context.TODO(), scheduleCacheKey, &schedule)
	return &schedule, err
}

func getUserSchedulesCacheKey(userId primitive.ObjectID, page int) string {
	return "req:cache:user:schedules:" + userId.Hex() + ":" + strconv.Itoa(page)
}

func CacheUserSchedules(userId primitive.ObjectID, schedule *[]models.Schedule, page int) {
	if !services.Config.UseRedis {
		return
	}

	userScheduleCacheKey := getUserSchedulesCacheKey(userId, page)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   userScheduleCacheKey,
		Value: schedule,
		TTL:   time.Minute,
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
