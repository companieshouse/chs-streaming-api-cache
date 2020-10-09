package cache

import (
	"fmt"
	"github.com/mediocregopher/radix/v3"
	"log"
	"strconv"
	"time"
)

const (
	ZADD             = "ZADD"
	ZRANGEBYSCORE    = "ZRANGEBYSCORE"
	ZREMRANGEBYSCORE = "ZREMRANGEBYSCORE"
)

type CacheService interface {
	// Insert new entities into sorted sets with the offset number as the score
	Create(key string, delta string, score int64) error
	//Fetch a range of offsets using a specified offset number as the starting offset
	Read(key string, offset int64) ([]string, error)
	//Delete a range of offsets from the sorted sets using matching timestamps.
	Delete(key string) error
}

type RedisCacheService struct {
	pool            *radix.Pool
	expiryInSeconds time.Duration
}

func NewRedisCacheService(network string, url string, size int, expiryInSeconds int64) CacheService {

	pool, err := radix.NewPool(network, url, size)
	if err != nil {
		panic(err)
	}
	return &RedisCacheService{
		pool:            pool,
		expiryInSeconds: time.Duration(expiryInSeconds),
	}
}

func (r RedisCacheService) Create(key string, delta string, score int64) error {
	log.Printf("Creating new cache entry for key=%s", key)

	offset := strconv.FormatInt(score, 10)
	err := r.pool.Do(radix.Cmd(nil, ZADD, key, offset, delta))
	if err != nil {
		return err
	}
	log.Printf("Creating time to live cache entry for key=%s and offset=%s", key, offset)
	timestamp := strconv.FormatInt(time.Now().UnixNano()+int64(time.Nanosecond*time.Second*r.expiryInSeconds), 10)
	return r.pool.Do(radix.Cmd(nil, ZADD, key+":ttl", timestamp, offset))
}

func (r RedisCacheService) Read(key string, offset int64) ([]string, error) {
	var result []string
	err := r.pool.Do(radix.Cmd(&result, ZRANGEBYSCORE, key, strconv.FormatInt(offset, 10), "inf"))
	if err == nil {
		log.Printf("Retrieved %d cached entries for key=%s and offset=%d", len(result), key, offset)
	}
	for _, delta := range result {
		fmt.Println("Reading delta: ", delta)
	}

	return result, err
}

func (r RedisCacheService) Delete(key string) error {
	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	// find expired timestamp entries
	var result []string
	if err := r.pool.Do(radix.Cmd(&result, ZRANGEBYSCORE, key+":ttl", "-inf", now)); err != nil {
		return err
	}

	if len(result) > 0 {
		// remove expired timestamp entries
		var count int
		if err := r.pool.Do(radix.Cmd(&count, ZREMRANGEBYSCORE, key+":ttl", "-inf", now)); err != nil {
			return err
		}
		log.Printf("Removed %d expired entries for key %s\n", count, key+":ttl")
		// remove equivalent delta entries
		min := result[0]
		max := result[len(result)-1]
		log.Printf("Removing expired entries for %s in range %s to %s\n", key, min, max)
		if err := r.pool.Do(radix.Cmd(&count, ZREMRANGEBYSCORE, key, min, max)); err != nil {
			return err
		}
		log.Printf("Removed %d expired entries for key %s\n", count, key)
	} else {
		log.Printf("No expired entries for %s\n", key)
	}
	return nil
}
