package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/companieshouse/chs-streaming-api-cache/cache"
	"net/http"
	"os"
	"sync"
	"time"
)

type Client struct {
	baseurl string
	broker  Publishable
	client  Gettable
	service cache.CacheService
	key string
	wg *sync.WaitGroup
}

type Publishable interface {
	Publish (msg string)
}

type Gettable interface {
	Get(url string) (resp *http.Response, err error)
}

//The result of the operation.
type Result struct{
	Data   string `json:"data"`
	Offset int64 `json:"offset"`
}

func NewClient(baseurl string, broker Publishable, client Gettable, service cache.CacheService, key string) *Client{
	return &Client {
		baseurl,
		broker,
		client,
		service,
		key,
		nil,
	}
}

func (c *Client) Connect(){
	resp, _ := c.client.Get(c.baseurl)
	body := resp.Body
	reader := bufio.NewReader(body)
	go c.loop(reader)
}

func (c *Client) loop(reader *bufio.Reader) {

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "error during resp.Body read:%s\n", err)
			continue
		}
		result := &Result{}
		json.Unmarshal(line, result)
		c.service.Create(c.key, result.Data, result.Offset)
		c.broker.Publish(result.Data)
		if c.wg != nil {
			c.wg.Done()
		}
		time.Sleep(600)
	}
}
