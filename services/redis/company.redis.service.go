package redisServices

import (
	models "attendit/backend/models/db"
	"attendit/backend/services"
	"context"
	"errors"
	"github.com/go-redis/cache/v8"
	"time"
)

func getCompanyCacheKey() string {
	return "req:cache:company"
}

func CacheCompany(company *models.Company) {
	if !services.Config.UseRedis {
		return
	}

	companyCacheKey := getCompanyCacheKey()

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   companyCacheKey,
		Value: company,
		TTL:   time.Second * 30,
	})
}

func GetCompanyFromCache() (*models.Company, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	note := &models.Company{}
	companyCacheKey := getCompanyCacheKey()
	err := services.GetRedisCache().Get(context.TODO(), companyCacheKey, note)
	return note, err
}

}

	if !services.Config.UseRedis {
		return
	}


	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
	})
}

	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

}
