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

func getUserPaidLeavesCacheKey(userId primitive.ObjectID, page int) string {
	return "req:cache:paidLeaves:" + userId.Hex() + ":" + strconv.Itoa(page)
}

func CacheUserPaidLeaves(userId primitive.ObjectID, paidLeaves []*db.PaidLeave, page int) {
	if !services.Config.UseRedis {
		return
	}

	userPaidLeavesCacheKey := getUserPaidLeavesCacheKey(userId, page)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   userPaidLeavesCacheKey,
		Value: paidLeaves,
		TTL:   time.Minute,
	})
}

func GetUserPaidLeavesFromCache(userId primitive.ObjectID, page int) ([]*db.PaidLeave, error) {
	if !services.Config.UseRedis {
		return nil, nil
	}

	var paidLeaves []*db.PaidLeave
	userPaidLeavesCacheKey := getUserPaidLeavesCacheKey(userId, page)
	err := services.GetRedisCache().Get(context.TODO(), userPaidLeavesCacheKey, &paidLeaves)
	return paidLeaves, err
}

func getPaidLeavesCacheKey(page int) string {
	return "req:cache:paidLeaves:" + strconv.Itoa(page)
}

func CachePaidLeaves(paidLeaves []*db.PaidLeave, page int) {
	if !services.Config.UseRedis {
		return
	}

	paidLeavesCacheKey := getPaidLeavesCacheKey(page)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   paidLeavesCacheKey,
		Value: paidLeaves,
		TTL:   time.Minute,
	})
}

func GetPaidLeavesFromCache(page int) ([]*db.PaidLeave, error) {
	if !services.Config.UseRedis {
		return nil, nil
	}

	var paidLeaves []*db.PaidLeave
	paidLeavesCacheKey := getPaidLeavesCacheKey(page)
	err := services.GetRedisCache().Get(context.TODO(), paidLeavesCacheKey, &paidLeaves)
	return paidLeaves, err
}

func getPaidLeavesByStatusCacheKey(status string, page int) string {
	return "req:cache:paidLeaves:" + status + ":" + strconv.Itoa(page)
}

func CachePaidLeavesByStatus(status string, paidLeaves []*db.PaidLeave, page int) {
	if !services.Config.UseRedis {
		return
	}

	paidLeavesCacheKey := getPaidLeavesByStatusCacheKey(status, page)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   paidLeavesCacheKey,
		Value: paidLeaves,
		TTL:   time.Minute,
	})

}

func GetPaidLeavesByStatusFromCache(status string, page int) ([]*db.PaidLeave, error) {
	if !services.Config.UseRedis {
		return nil, nil
	}

	var paidLeaves []*db.PaidLeave
	paidLeavesCacheKey := "req:cache:paidLeaves:" + status + ":" + strconv.Itoa(page)
	err := services.GetRedisCache().Get(context.TODO(), paidLeavesCacheKey, &paidLeaves)
	return paidLeaves, err
}
