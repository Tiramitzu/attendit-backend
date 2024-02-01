package redisServices

import (
	models "attendit/backend/models/db"
	"attendit/backend/services"
	"context"
	"errors"
	"github.com/go-redis/cache/v8"
	"time"
)

func getUsersCacheKey(page int) string {
	return "req:cache:users:" + string(rune(page))
}

func CacheUsers(page int, users []*models.User) {
	if !services.Config.UseRedis {
		return
	}

	usersCacheKey := getUsersCacheKey(page)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   usersCacheKey,
		Value: users,
		TTL:   time.Second * 30,
	})
}

func GetUsersFromCache(page int) ([]*models.User, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	var users []*models.User
	usersCacheKey := getUsersCacheKey(page)
	err := services.GetRedisCache().Get(context.TODO(), usersCacheKey, &users)
	return users, err
}
