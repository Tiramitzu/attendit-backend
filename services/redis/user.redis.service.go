package redisServices

import (
	db "attendit/backend/models/db"
	"attendit/backend/services"
	"context"
	"errors"
	"github.com/go-redis/cache/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func getUserCacheKey(userId primitive.ObjectID) string {
	return "req:cache:user:" + userId.Hex()
}

func CacheUser(user *db.User) {
	if !services.Config.UseRedis {
		return
	}

	userCacheKey := getUserCacheKey(user.ID)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   userCacheKey,
		Value: user,
		TTL:   time.Second * 30,
	})
}

func GetUserFromCache(userId primitive.ObjectID) (*db.User, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	user := &db.User{}
	userCacheKey := getUserCacheKey(userId)
	err := services.GetRedisCache().Get(context.TODO(), userCacheKey, user)
	return user, err
}
