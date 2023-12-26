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

func getUserAttendanceByCompanyCacheKey(userId primitive.ObjectID, companyId primitive.ObjectID) string {
	return "req:cache:user:attendance:" + userId.Hex() + ":" + companyId.Hex()
}

func CacheUserAttendancesByCompany(userId primitive.ObjectID, companyId primitive.ObjectID, attendance *[]models.Attendance) {
	if !services.Config.UseRedis {
		return
	}

	userAttendanceByCompanyCacheKey := getUserAttendanceByCompanyCacheKey(userId, companyId)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   userAttendanceByCompanyCacheKey,
		Value: attendance,
		TTL:   time.Minute,
	})
}

func GetUserAttendancesByCompanyFromCache(userId primitive.ObjectID, companyId primitive.ObjectID) (*[]models.Attendance, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	var attendances []models.Attendance
	userAttendanceByCompanyCacheKey := getUserAttendanceByCompanyCacheKey(userId, companyId)
	err := services.GetRedisCache().Get(context.TODO(), userAttendanceByCompanyCacheKey, &attendances)
	return &attendances, err
}

func getCompanyAttendancesCacheKey(companyId primitive.ObjectID) string {
	return "req:cache:company:attendances:" + companyId.Hex()
}

func CacheCompanyAttendances(companyId primitive.ObjectID, attendances []*models.Attendance) {
	if !services.Config.UseRedis {
		return
	}

	companyAttendancesCacheKey := getCompanyAttendancesCacheKey(companyId)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   companyAttendancesCacheKey,
		Value: attendances,
		TTL:   time.Minute,
	})
}

func GetCompanyAttendancesFromCache(companyId primitive.ObjectID) ([]*models.Attendance, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	var attendances []*models.Attendance
	companyAttendancesCacheKey := getCompanyAttendancesCacheKey(companyId)
	err := services.GetRedisCache().Get(context.TODO(), companyAttendancesCacheKey, &attendances)
	return attendances, err
}
