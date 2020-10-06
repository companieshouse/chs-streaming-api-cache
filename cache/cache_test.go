package cache

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const redisURL = "localhost:32771"
const expiryInSeconds int64 = 120

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
	Convey("Given an entry exists in the redis cache sortedSet", t, func() {
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

func TestRedisCacheService_Delete(t *testing.T) {
	Convey("Given entries exist in the redis cache sortedSet", t, func() {
		redisCacheService := NewRedisCacheService("tcp", redisURL, 10, 0)
		for score := 0; score < 10 ; score++  {
			err := redisCacheService.Create("stream:test3", "{id : 124}", int64(score))
			if err != nil{
				t.Error("Failed: " + err.Error())
			}
		}
		Convey("When a call is made to delete expired entries", func(){
			err := redisCacheService.Delete("stream:test3")
			if err != nil{
				t.Error("Failed: " + err.Error())
			}
			Convey("Then the expired entries should be removed from the cache", func(){
				actual, err := redisCacheService.Read("stream:test3", 10)
				if err != nil{
					t.Error("Failed: " + err.Error())
				}
				So(len(actual), ShouldEqual, 0)
			})
		})
	})
}
