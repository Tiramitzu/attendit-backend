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

func getCompanyCacheKey(companyId primitive.ObjectID) string {
	return "req:cache:company:" + companyId.Hex()
}

func CacheOneCompany(company *models.Company) {
	if !services.Config.UseRedis {
		return
	}

	companyCacheKey := getCompanyCacheKey(company.ID)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   companyCacheKey,
		Value: company,
		TTL:   time.Minute,
	})
}

func GetCompanyFromCache(companyId primitive.ObjectID) (*models.Company, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	note := &models.Company{}
	companyCacheKey := getCompanyCacheKey(companyId)
	err := services.GetRedisCache().Get(context.TODO(), companyCacheKey, note)
	return note, err
}

func getCompaniesCacheKey(userId primitive.ObjectID) string {
	return "req:cache:companies" + userId.Hex()
}

func CacheCompanies(userId primitive.ObjectID, companies *[]models.Company) {
	if !services.Config.UseRedis {
		return
	}

	companiesCacheKey := getCompaniesCacheKey(userId)

	_ = services.GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   companiesCacheKey,
		Value: companies,
		TTL:   time.Minute,
	})
}

func GetCompaniesFromCache(userId primitive.ObjectID) (*[]models.Company, error) {
	if !services.Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	var companies *[]models.Company
	companiesCacheKey := getCompaniesCacheKey(userId)
	err := services.GetRedisCache().Get(context.TODO(), companiesCacheKey, &companies)
	return companies, err
}
