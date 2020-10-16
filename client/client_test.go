package client

import (
	"github.com/companieshouse/chs-streaming-api-cache/broker"
	"github.com/companieshouse/chs.go/log"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"net/http"
	"strings"
	"sync"
	"testing"
)

type mockBroker struct {
	mock.Mock
}

func (b *mockBroker) Publish(msg string) {
	b.Called(msg)
}

type mockHttpClient struct {
	mock.Mock
}

func (c *mockHttpClient) Get(url string) (resp *http.Response, err error) {
	args := c.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

type mockCacheService struct {
	mock.Mock
}

func (s *mockCacheService) Create(key string, delta string, offset int64) error {
	args := s.Called(key, delta, offset)
	return args.Error(0)
}

func (s *mockCacheService) Read(key string, offset int64) ([]string, error) {
	args := s.Called(key, offset)
	return args.Get(0).([]string), args.Error(1)
}

type mockLogger struct {
	mock.Mock
}

func (l *mockLogger) Info(msg string, data ...log.Data) {
	panic("implement me")
}

func (l *mockLogger) InfoR(req *http.Request, message string, data ...log.Data) {
	panic("implement me")
}

func (l *mockLogger) Error(err error, data ...log.Data){
	l.Called(mock.Anything)
}

type mockBody struct {
	*strings.Reader
}

func (b *mockBody) Close() error {
	return nil
}

type mockLogger struct {
	mock.Mock
}

func (l *mockLogger) Info(msg string, data ...log.Data) {
	panic("implement me")
}

func (l *mockLogger) InfoR(req *http.Request, message string, data ...log.Data) {
	panic("implement me")
}

func (l *mockLogger) Error(err error, data ...log.Data){
	l.Called(mock.Anything)
}

func TestNewClient(t *testing.T) {
	Convey("given a new client instance is created", t, func() {
		actual := NewClient("baseurl", &broker.Broker{}, &http.Client{}, &mockCacheService{}, "key", &mockLogger{})
		Convey("then a new client should be created", func() {
			So(actual, ShouldNotBeNil)
			So(actual.baseurl, ShouldEqual, "baseurl")
			So(actual.broker, ShouldResemble, &broker.Broker{})
			So(actual.httpClient, ShouldResemble, &http.Client{})
		})
	})
}

func TestPublishToBroker(t *testing.T) {
	Convey("given a mock broker and http client is called", t, func() {
		broker := &mockBroker{}
		broker.On("Publish", mock.Anything).Return()
		httpClient := &mockHttpClient{}
		httpClient.On("Get", mock.Anything).Return(&http.Response{StatusCode: 200,
			Body: &mockBody{strings.NewReader("{\"data\":\"{\\\"greetings\\\":\\\"hello\\\"}\",\"offset\":43}\n")},
		}, nil)
		service := &mockCacheService{}
		service.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		logger := &mockLogger{}
		logger.On("Error", mock.Anything).Return(nil)
		client := NewClient("baseurl", broker, httpClient, service, "key", logger)
		client.wg = new(sync.WaitGroup)
		Convey("when a new message is published", func() {
			client.wg.Add(1)
			client.Connect()
			client.wg.Wait()
			Convey("Then the message should be written to the cache and forwarded to the broker", func() {
				So(service.AssertCalled(t, "Create", "key", "{\"greetings\":\"hello\"}", int64(43)), ShouldBeTrue)
				So(broker.AssertCalled(t, "Publish", "{\"greetings\":\"hello\"}"), ShouldBeTrue)
			})
		})
	})
}

//error handling if http client returns an error
