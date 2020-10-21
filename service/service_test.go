package service

import (
	"github.com/companieshouse/chs-streaming-api-cache/config"
	"github.com/gorilla/pat"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCreateNewService(t *testing.T) {
	Convey("When a new service instance is constructed", t, func() {
		configuration := &CacheConfiguration{
			Configuration: &config.Config{
				RedisUrl:             "localhost:6379",
				CacheExpiryInSeconds: 2,
			},
			Router: pat.New(),
		}
		actual := NewCacheService(configuration)
		Convey("Then a new service instance reference should be returned", func() {
			So(actual, ShouldNotBeNil)
			So(actual.broker, ShouldNotBeNil)
			So(actual.router, ShouldEqual, configuration.Router)
		})
	})
}

func TestBindTopic(t *testing.T) {
	Convey("Given a new service instance has been constructed", t, func() {
		configuration := &CacheConfiguration{
			Configuration: &config.Config{
				RedisUrl:             "localhost:6379",
				CacheExpiryInSeconds: 2,
			},
			Router: pat.New(),
		}

		service := NewCacheService(configuration)
		Convey("When a topic is bound to it", func() {
			actual := service.WithTopic("topic")
			Convey("Then the topic should be added to the service", func() {
				So(actual, ShouldEqual, service)
				So(actual.topic, ShouldEqual, "topic")
			})
		})
	})
}


