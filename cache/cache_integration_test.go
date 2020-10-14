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
	envVariables.expiryInSeconds = 2

	redisCacheService = NewRedisCacheService("tcp", envVariables.redisURL, 10, envVariables.expiryInSeconds)

	return redisC
}

func stopContainer(container testcontainers.Container) {
	fmt.Println("Stopping container")
	container.Terminate(ctx)
}

func TestRedisCacheService_Create(t *testing.T) {
	Convey("When I create a cached entry", t, func() {
		const topic = "stream:test"
		err := redisCacheService.Create(topic, "{id : 123}", 20)
		Convey("Then the cached entry should be created", func() {
			if err != nil {
				t.Error("Failed: " + err.Error())
			}
		})
	})
}

func TestRedisCacheService_Read(t *testing.T) {
	Convey("Given an entry exists in the redis cache sortedSet", t, func() {
		const topic = "stream:test2"
		err := redisCacheService.Create(topic, "{id : 124}", 21)
		if err != nil {
			t.Error("Failed: " + err.Error())
		}
		Convey("When I fetch the cached entries", func() {
			actual, err := redisCacheService.Read(topic, 0)
			if err != nil {
				t.Error("Failed: " + err.Error())
			}
			Convey("Then the cached entries for the given topic should be found", func() {
				So(len(actual), ShouldEqual, 1)
				var expected = [1]string{"{id : 124}"}
				So(actual[0], ShouldEqual, expected[0])
			})
		})
	})
}

func TestRedisCacheService_ReadFromAGivenOffset(t *testing.T) {
	Convey("Given an entry exists in the redis cache sortedSet", t, func() {
		const topic = "stream:test3"
		for score := 10; score < 20; score++ {
			delta := fmt.Sprintf("{id : %d}", score)
			err := redisCacheService.Create(topic, delta, int64(score))
			if err != nil {
				t.Error("Failed: " + err.Error())
			}
		}
		Convey("When I fetch the cached entries for a given offset", func() {
			actualArray, err := redisCacheService.Read(topic, 15)
			if err != nil {
				t.Error("Failed: " + err.Error())
			}
			Convey("Then only the cached entries for the given offset should be found", func() {
				So(len(actualArray), ShouldEqual, 5)
				var expected = []string{"{id : 15}", "{id : 16}", "{id : 17}", "{id : 18}", "{id : 19}"}
				for index, actual := range actualArray {
					fmt.Println("Actual: ", index, actual)
					So(actual, ShouldEqual, expected[index])
				}
			})
		})
	})
}

func TestRedisCacheService_ReadDoesNotReturnExpiredEntries(t *testing.T) {
	Convey("Given an entry exists in the redis cache sortedSet", t, func() {
		const topic = "stream:test4"
		for score := 10; score < 20; score++ {
			delta := fmt.Sprintf("{id : %d}", score)
			err := redisCacheService.Create(topic, delta, int64(score))
			if err != nil {
				t.Error("Failed: " + err.Error())
			}
		}
		Convey("When the entries become expired", func() {
			fmt.Println("Waiting for cache entries to expire...")
			time.Sleep(time.Duration(envVariables.expiryInSeconds) * time.Second)
			Convey("Then the expired entries for the given offset should not be returned", func() {
				actualArray, err := redisCacheService.Read(topic, 10)
				if err != nil {
					t.Error("Failed: " + err.Error())
				}
				So(len(actualArray), ShouldEqual, 0)
			})
		})
	})
}
