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

func getUserAttendanceByCompanyCacheKey(userId primitive.ObjectID, page int) string {
	return "req:cache:user:attendance:" + userId.Hex() + ":" + string(rune(page))
}

func CacheUserAttendancesByCompany(userId primitive.ObjectID, attendance []models.Attendance, page int) {
	if !services.Config.UseRedis {
		return
	}

	userAttendanceByCompanyCacheKey := getUserAttendanceByCompanyCacheKey(userId, page)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   userAttendanceByCompanyCacheKey,
		Value: attendance,
		TTL:   time.Second * 30,
	})
}

func GetUserAttendancesFromCache(userId primitive.ObjectID, page int) ([]models.Attendance, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	var attendances []models.Attendance
	userAttendanceByCompanyCacheKey := getUserAttendanceByCompanyCacheKey(userId, page)
	err := services.GetRedisCache().Get(context.TODO(), userAttendanceByCompanyCacheKey, &attendances)
	return attendances, err
}

func getAttendancesCacheKey(page int) string {
	return "req:cache:company:attendances:" + string(rune(page))
}

func CacheAttendances(page int, attendances []*models.Attendance) {
	if !services.Config.UseRedis {
		return
	}

	companyAttendancesCacheKey := getAttendancesCacheKey(page)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   companyAttendancesCacheKey,
		Value: attendances,
		TTL:   time.Second * 30,
	})
}

func GetAttendancesFromCache(page int) ([]*models.Attendance, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	var attendances []*models.Attendance
	companyAttendancesCacheKey := getAttendancesCacheKey(page)
	err := services.GetRedisCache().Get(context.TODO(), companyAttendancesCacheKey, &attendances)
	return attendances, err
}
