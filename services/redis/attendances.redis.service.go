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

func getUserAttendanceByCompanyCacheKey(userId primitive.ObjectID) string {
	return "req:cache:user:attendance:" + userId.Hex()
}

func CacheUserAttendancesByCompany(userId primitive.ObjectID, attendance *[]models.Attendance) {
	if !services.Config.UseRedis {
		return
	}

	userAttendanceByCompanyCacheKey := getUserAttendanceByCompanyCacheKey(userId)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   userAttendanceByCompanyCacheKey,
		Value: attendance,
		TTL:   time.Minute,
	})
}

func GetUserAttendancesFromCache(userId primitive.ObjectID) (*[]models.Attendance, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	var attendances []models.Attendance
	userAttendanceByCompanyCacheKey := getUserAttendanceByCompanyCacheKey(userId)
	err := services.GetRedisCache().Get(context.TODO(), userAttendanceByCompanyCacheKey, &attendances)
	return &attendances, err
}

func getCompanyAttendancesCacheKey(page int) string {
	return "req:cache:company:attendances:" + string(rune(page))
}

func CacheCompanyAttendances(page int, attendances []*models.Attendance) {
	if !services.Config.UseRedis {
		return
	}

	companyAttendancesCacheKey := getCompanyAttendancesCacheKey(page)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   companyAttendancesCacheKey,
		Value: attendances,
		TTL:   time.Minute,
	})
}

func GetCompanyAttendancesFromCache(page int) ([]*models.Attendance, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	var attendances []*models.Attendance
	companyAttendancesCacheKey := getCompanyAttendancesCacheKey(page)
	err := services.GetRedisCache().Get(context.TODO(), companyAttendancesCacheKey, &attendances)
	return attendances, err
}
