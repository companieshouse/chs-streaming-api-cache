package cache

import (
	"fmt"
	cache_config "github.com/companieshouse/chs-streaming-api-cache/config"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
	"time"
)

var envVariables struct{
	redisURL string
	expiryInSeconds int64
}

var redisCacheService CacheService

func setup(t *testing.T){
	fmt.Println("Main")
	config, err := cache_config.Get()
	if err != nil {
		panic(err)
	}
	envVariables.redisURL = config.RedisUrl
	envVariables.expiryInSeconds, err = strconv.ParseInt(config.CacheExpiryInSeconds, 10, 64)
	if err != nil {
		panic(err)
	}
	redisCacheService = NewRedisCacheService("tcp", envVariables.redisURL, 10, envVariables.expiryInSeconds)
}

func TestRedisCacheService_Create(t *testing.T) {
	Convey("Given an instance od the cache service", t, func() {
		setup(t)
		Convey("When I create a cached entry", func() {
			err := redisCacheService.Create("stream:test", "{id : 123}", 23)
			Convey("Then the cached entry should be created", func() {
				if err != nil{
					t.Error("Failed: " + err.Error())
				}
			})
		})
	})
}

func TestRedisCacheService_Read(t *testing.T) {
	Convey("Given an entry exists in the redis cache sortedSet", t, func() {
		setup(t)
		err := redisCacheService.Create("stream:test", "{id : 124}", 20)
		if err != nil{
			t.Error("Failed: " + err.Error())
		}
		Convey("When I fetch the cached entries fro a given offset", func() {
			actual, err := redisCacheService.Read("stream:test", 20)
			if err != nil{
				t.Error("Failed: " + err.Error())
			}
			Convey("Then only the cached entries for the given offset should be found", func() {
				So(len(actual), ShouldEqual, 1)
				var expected = [1]string{"{id : 124}"}
				So(actual[0], ShouldEqual, expected[0])
			})
		})
	})
}

func TestRedisCacheService_Delete(t *testing.T) {
	Convey("Given entries exist in the redis cache sortedSet", t, func() {
		setup(t)
		for score := 10; score < 20 ; score++  {
			err := redisCacheService.Create("stream:test3", "{id : 124}", int64(score))
			if err != nil{
				t.Error("Failed: " + err.Error())
			}
		}
		Convey("When a call is made to delete expired entries", func(){
			// allow entries to expire
			time.Sleep(time.Duration(envVariables.expiryInSeconds) * time.Second)
			err := redisCacheService.Delete("stream:test3")
			if err != nil{
				t.Error("Failed: " + err.Error())
			}
			Convey("Then the expired entries should be removed from the cache", func(){
				actual, err := redisCacheService.Read("stream:test3", 20)
				if err != nil{
					t.Error("Failed: " + err.Error())
				}
				So(len(actual), ShouldEqual, 0)
			})
		})
	})
}