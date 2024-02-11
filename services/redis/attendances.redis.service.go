package redisServices

import (
	"attendit/backend/models"
	db "attendit/backend/models/db"
	"attendit/backend/services"
	"context"
	"errors"
	"github.com/go-redis/cache/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

func getUserAttendanceByCompanyCacheKey(userId primitive.ObjectID, page int) string {
	return "req:cache:user:attendance:" + userId.Hex() + ":" + strconv.Itoa(page)
}

func CacheUserAttendancesByCompany(userId primitive.ObjectID, attendance []db.Attendance, page int) {
	if !services.Config.UseRedis {
		return
	}

	userAttendanceByCompanyCacheKey := getUserAttendanceByCompanyCacheKey(userId, page)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   userAttendanceByCompanyCacheKey,
		Value: attendance,
		TTL:   time.Minute,
	})
}

func GetUserAttendancesFromCache(userId primitive.ObjectID, page int) ([]db.Attendance, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	var attendances []db.Attendance
	userAttendanceByCompanyCacheKey := getUserAttendanceByCompanyCacheKey(userId, page)
	err := services.GetRedisCache().Get(context.TODO(), userAttendanceByCompanyCacheKey, &attendances)
	return attendances, err
}

func getAttendancesCacheKey(page int) string {
	return "req:cache:company:attendances:" + strconv.Itoa(page)
}

func CacheAttendances(page int, attendances []*db.Attendance) {
	if !services.Config.UseRedis {
		return
	}

	companyAttendancesCacheKey := getAttendancesCacheKey(page)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   companyAttendancesCacheKey,
		Value: attendances,
		TTL:   time.Minute,
	})
}

func GetAttendancesFromCache(page int) ([]*db.Attendance, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	var attendances []*db.Attendance
	companyAttendancesCacheKey := getAttendancesCacheKey(page)
	err := services.GetRedisCache().Get(context.TODO(), companyAttendancesCacheKey, &attendances)
	return attendances, err
}

func getAttendanceTotalCacheKey() string {
	return "req:cache:company:attendance:total"
}

func CacheAttendanceTotal(total models.AttendanceTotal) {
	if !services.Config.UseRedis {
		return
	}

	attendanceTotalCacheKey := getAttendanceTotalCacheKey()

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   attendanceTotalCacheKey,
		Value: total,
		TTL:   time.Minute,
	})
}

func GetAttendanceTotalFromCache() (models.AttendanceTotal, error) {
	if !services.Config.UseRedis {
		return models.AttendanceTotal{}, errors.New("no redis client, set USE_REDIS in .env")
	}

	var total models.AttendanceTotal
	attendanceTotalCacheKey := getAttendanceTotalCacheKey()
	err := services.GetRedisCache().Get(context.TODO(), attendanceTotalCacheKey, &total)
	return total, err
}

func getAttendanceByDateCacheKey(fromDate string, toDate string, page int) string {
	return "req:cache:company:attendance:" + fromDate + ":" + toDate + ":" + strconv.Itoa(page)
}

func CacheAttendancesByDate(fromDate string, toDate string, page int, attendances []*db.Attendance) {
	if !services.Config.UseRedis {
		return
	}

	attendanceByDateCacheKey := getAttendanceByDateCacheKey(fromDate, toDate, page)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   attendanceByDateCacheKey,
		Value: attendances,
		TTL:   time.Minute,
	})
}

func GetAttendancesByDateFromCache(fromDate string, toDate string, page int) ([]*db.Attendance, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	var attendances []*db.Attendance
	attendanceByDateCacheKey := getAttendanceByDateCacheKey(fromDate, toDate, page)
	err := services.GetRedisCache().Get(context.TODO(), attendanceByDateCacheKey, &attendances)
	return attendances, err
}

func getAttendanceTotalByDateCacheKey(fromDate string, toDate string) string {
	return "req:cache:company:attendance:total:" + fromDate + ":" + toDate
}

func CacheAttendanceTotalByDate(fromDate string, toDate string, total models.AttendanceTotal) {
	if !services.Config.UseRedis {
		return
	}

	attendanceTotalByDateCacheKey := getAttendanceTotalByDateCacheKey(fromDate, toDate)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   attendanceTotalByDateCacheKey,
		Value: total,
		TTL:   time.Minute,
	})
}

func GetAttendanceTotalByDateFromCache(fromDate string, toDate string) (models.AttendanceTotal, error) {
	if !services.Config.UseRedis {
		return models.AttendanceTotal{}, errors.New("no redis client, set USE_REDIS in .env")
	}

	var total models.AttendanceTotal
	attendanceTotalByDateCacheKey := getAttendanceTotalByDateCacheKey(fromDate, toDate)
	err := services.GetRedisCache().Get(context.TODO(), attendanceTotalByDateCacheKey, &total)
	return total, err
}
