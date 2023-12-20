package services

import (
	models "attendit/backend/models/db"
	"context"
	"errors"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"sync"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB() {
	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, Config.MongodbDatabase, options.Client().ApplyURI(Config.MongodbUri))
	if err != nil {
		panic(err)
	}

	log.Println("Connected to MongoDB!")
}

var redisDefaultClient *redis.Client
var redisDefaultOnce sync.Once

var redisCache *cache.Cache
var redisCacheOnce sync.Once

func GetRedisDefaultClient() *redis.Client {
	redisDefaultOnce.Do(func() {
		redisDefaultClient = redis.NewClient(&redis.Options{
			Addr:     Config.RedisDefaultAddr,
			Password: Config.RedisDefaultPassword,
		})
	})

	return redisDefaultClient
}

func GetRedisCache() *cache.Cache {
	redisCacheOnce.Do(func() {
		redisCache = cache.New(&cache.Options{
			Redis:      GetRedisDefaultClient(),
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		})
	})

	return redisCache
}

func CheckRedisConnection() {
	redisClient := GetRedisDefaultClient()
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	log.Println("Connected to Redis!")
}

func getCompanyCacheKey(companyId primitive.ObjectID) string {
	return "req:cache:company:" + companyId.Hex()
}

func CacheOneCompany(company *models.Company) {
	if !Config.UseRedis {
		return
	}

	companyCacheKey := getCompanyCacheKey(company.ID)

	_ = GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   companyCacheKey,
		Value: company,
		TTL:   time.Minute,
	})
}

func GetCompanyFromCache(companyId primitive.ObjectID) (*models.Company, error) {
	if !Config.UseRedis {
		return nil, errors.New("no redis client, set USE_REDIS in .env")
	}

	note := &models.Company{}
	companyCacheKey := getCompanyCacheKey(companyId)
	err := GetRedisCache().Get(context.TODO(), companyCacheKey, note)
	return note, err
}
