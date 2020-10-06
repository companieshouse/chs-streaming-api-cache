package cache

import (
	"github.com/mediocregopher/radix/v3"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var redisURL = "localhost:32771"
var expiryInSeconds int64 = 120

func TestRedisCacheService_Create(t *testing.T) {
	Convey("Given an instance od the cache service", t, func() {
		redisCacheService := NewRedisCacheService("tcp", redisURL, 10, expiryInSeconds)
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
	Convey("Given an entry exits in the redis cache sortedSet", t, func() {
		redisCacheService := NewRedisCacheService("tcp", redisURL, 10, expiryInSeconds)
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

type mockRedisCacheService struct{
	pool *radix.Pool
}

