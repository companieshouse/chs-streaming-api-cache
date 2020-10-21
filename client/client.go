package client

import (
	"bufio"
	"encoding/json"
	. "github.com/companieshouse/chs-streaming-api-cache/cache"
	"github.com/companieshouse/chs-streaming-api-cache/logger"
	"github.com/companieshouse/chs.go/log"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	baseurl      string
	broker       Publishable
	httpClient   Gettable
	cacheService Cacheable
	key          string
	logger       logger.Logger
	wg           *sync.WaitGroup
}

type Publishable interface {
	Publish(msg string)
}

type Gettable interface {
	Get(url string) (resp *http.Response, err error)
}

//The result of the operation.
type Result struct {
	Data   string `json:"data"`
	Offset int64  `json:"offset"`
}

func NewClient(baseurl string, broker Publishable, client Gettable, service Cacheable, key string, logger logger.Logger) *Client {
	return &Client{
		baseurl,
		broker,
		client,
		service,
		key,
		logger,
		nil,
	}
}

func (c *Client) Connect() {
	resp, err := c.httpClient.Get(c.baseurl)
	if err != nil {
		c.logger.Error(err, log.Data{})
		panic(err)
	}
	body := resp.Body
	reader := bufio.NewReader(body)
	go c.loop(reader)
}

func (c *Client) loop(reader *bufio.Reader) {

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			c.logger.Error(err, log.Data{})
			continue
		}
		result := &Result{}
		err = json.Unmarshal(line, result)
		if err != nil {
			c.logger.Error(err, log.Data{})
			continue
		}
		err = c.cacheService.Create(c.key, result.Data, result.Offset)
		if err != nil {
			c.logger.Error(err, log.Data{})
			continue
		}
		c.broker.Publish(result.Data)
		if c.wg != nil {
			c.wg.Done()
		}
		time.Sleep(300)
	}
}

func (c *Client) Run() {
	c.Connect()
}