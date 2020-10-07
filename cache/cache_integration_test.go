package cache

import (
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
	"time"
)

var envVariables struct {
	redisURL        string
	expiryInSeconds int64
}

var redisCacheService CacheService
var ctx context.Context

func TestMain(m *testing.M) {
	redisC := startContainer()
	// run tests
	exitVal := m.Run()
	stopContainer(redisC)
	os.Exit(exitVal)
}

func startContainer() testcontainers.Container {
	fmt.Println("Setup Tasks")
	// spin up redis container
	ctx = context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	redisHost, err := redisC.Host(ctx)
	if err != nil {
		panic(err)
	}
	redisPort, err := redisC.MappedPort(ctx, "6379/tcp")
	if err != nil {
		panic(err)
	}

	envVariables.redisURL = fmt.Sprintf("%s:%s", redisHost, redisPort.Port())
	envVariables.expiryInSeconds = int64(1)

	redisCacheService = NewRedisCacheService("tcp", envVariables.redisURL, 10, envVariables.expiryInSeconds)

	return redisC
}

func stopContainer(container testcontainers.Container) {
	container.Terminate(ctx)
}

func TestRedisCacheService_Create(t *testing.T) {
	Convey("When I create a cached entry", t, func() {
		err := redisCacheService.Create("stream:test", "{id : 123}", 23)
		Convey("Then the cached entry should be created", func() {
			if err != nil {
				t.Error("Failed: " + err.Error())
			}
		})
	})
}

func TestRedisCacheService_Read(t *testing.T) {
	Convey("Given an entry exists in the redis cache sortedSet", t, func() {
		err := redisCacheService.Create("stream:test", "{id : 124}", 20)
		if err != nil {
			t.Error("Failed: " + err.Error())
		}
		Convey("When I fetch the cached entries fro a given offset", func() {
			actual, err := redisCacheService.Read("stream:test", 20)
			if err != nil {
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
		for score := 10; score < 20; score++ {
			err := redisCacheService.Create("stream:test3", "{id : 124}", int64(score))
			if err != nil {
				t.Error("Failed: " + err.Error())
			}
		}
		Convey("When a call is made to delete expired entries", func() {
			// allow entries to expire
			time.Sleep(time.Duration(envVariables.expiryInSeconds) * time.Second)
			err := redisCacheService.Delete("stream:test3")
			if err != nil {
				t.Error("Failed: " + err.Error())
			}
			Convey("Then the expired entries should be removed from the cache", func() {
				actual, err := redisCacheService.Read("stream:test3", 20)
				if err != nil {
					t.Error("Failed: " + err.Error())
				}
				So(len(actual), ShouldEqual, 0)
			})
		})
	})
}
