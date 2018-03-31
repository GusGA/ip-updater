package storage

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

var host = os.Getenv("REDIS_HOST")
var port = os.Getenv("REDIS_PORT")

var redisClient = redis.NewClient(&redis.Options{
	Addr:     host + ":" + port,
	Password: "", // no password set
	DB:       0,  // use default DB
})

func init() {
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Error connecting redis ", err)
	}
}

func SaveDomainData(domain string, data string) error {
	return redisClient.Set(domain, data, 0).Err()
}

func GetDomainData(domain string) (data string, err error) {
	data, err = redisClient.Get(domain).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return data, err
}

func Close() error {
	return redisClient.Close()
}
