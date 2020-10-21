package cache

import (
	"fmt"
	"github.com/mediocregopher/radix/v3"
	"log"
	"strconv"
)

const (
	ZADD          = "ZADD"
	ZRANGEBYSCORE = "ZRANGEBYSCORE"
	SET           = "SET"
	GET           = "GET"
	EXPIRE        = "EXPIRE"
)

type Cacheable interface {
	// Insert new entities into sorted sets with the offset number as the score
	Create(key string, delta string, score int64) error
	//Fetch a range of offsets using a specified offset number as the starting offset
	Read(key string, offset int64) ([]string, error)
}

type RedisCacheService struct {
	pool            *radix.Pool
	expiryInSeconds int64
}

func NewRedisCacheService(network string, url string, size int, expiryInSeconds int64) Cacheable {

	pool, err := radix.NewPool(network, url, size)
	if err != nil {
		panic(err)
	}
	return &RedisCacheService{
		pool:            pool,
		expiryInSeconds: expiryInSeconds,
	}
}

func (r RedisCacheService) Create(key string, delta string, offset int64) error {
	log.Printf("Creating new cache entries for key=%s", key)

	offsetAsString := strconv.FormatInt(offset, 10)
	offsetsKey := key + ":offsets"
	log.Printf("Creating sorted set cache entry for key=%s and offset=%s", offsetsKey, offsetAsString)
	deltaKey := key + ":" + offsetAsString
	if err := r.pool.Do(radix.Cmd(nil, ZADD, offsetsKey, offsetAsString, deltaKey)); err != nil {
		return err
	}

	log.Printf("Creating new cache entry for key=%s", deltaKey)
	expirySeconds := fmt.Sprint(r.expiryInSeconds)
	if err := r.pool.Do(radix.Cmd(nil, SET, deltaKey, delta)); err != nil {
		return err
	}
	log.Printf("Setting expiry time of %s seconds for key=%s\n", expirySeconds, deltaKey)
	if err := r.pool.Do(radix.Cmd(nil, EXPIRE, deltaKey, expirySeconds)); err != nil {
		return err
	}
	return nil
}

func (r RedisCacheService) Read(key string, offset int64) ([]string, error) {
	offsetAsString := strconv.FormatInt(offset, 10)
	var offsets []string
	err := r.pool.Do(radix.Cmd(&offsets, ZRANGEBYSCORE, key+":offsets", offsetAsString, "inf"))
	if err == nil {
		log.Printf("Retrieved %d cached entries for key=%s and offset=%d", len(offsets), key, offset)
	}
	var deltas []string
	for _, offset := range offsets {
		var delta string
		fmt.Printf("Reading delta for offset: %s\n", offset)
		if err := r.pool.Do(radix.Cmd(&delta, GET, offset)); err != nil {
			return nil, err
		}
		if len(delta) > 0 {
			deltas = append(deltas, delta)
		}
	}
	return deltas, err
}
