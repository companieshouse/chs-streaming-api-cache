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
	path         string
	broker       Publishable
	httpClient   Doable
	username     string
	cacheService Cacheable
	key          string
	logger       logger.Logger
	wg           *sync.WaitGroup
}

type Publishable interface {
	Publish(msg string)
}

type Doable interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

// The result of the operation.
type Result struct {
	Data   string `json:"data"`
	Offset int64  `json:"offset"`
}

func NewClient(baseurl string, path string, broker Publishable, client Doable, username string, service Cacheable, key string, logger logger.Logger) *Client {
	return &Client{
		baseurl:      baseurl,
		path:         path,
		broker:       broker,
		httpClient:   client,
		username:     username,
		cacheService: service,
		key:          key,
		logger:       logger,
		wg:           nil,
	}
}

func (c *Client) Connect() {
	url := c.baseurl + c.path
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(c.username, "")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error(err, log.Data{})
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		c.logger.Info("Unable to stream from backend endpoint", log.Data{"endpoint": c.baseurl, "path": c.path, "Http Status": resp.StatusCode})
		panic("Unable to stream from backend endpoint")
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
