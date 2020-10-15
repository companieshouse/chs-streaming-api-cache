package handlers

import (
	. "github.com/companieshouse/chs-streaming-api-cache/cache"
	"github.com/companieshouse/chs-streaming-api-cache/logger"
	"github.com/companieshouse/chs-streaming-api-cache/offset"
	"log"
	"net/http"
	"sync"
)

type Subscribable interface {
	Subscribe() (chan string, error)
	Unsubscribe(chan string) error
}

type RequestHandler struct {
	broker       Subscribable
	cacheService CacheService
	key          string
	logger       logger.Logger
	offset       offset.Interface
	wg           *sync.WaitGroup
}

func NewRequestHandler(broker Subscribable, cacheService CacheService, logger logger.Logger, topic string) *RequestHandler {
	return &RequestHandler{
		broker:       broker,
		logger:       logger,
		cacheService: cacheService,
		key:          topic,
		offset:       offset.NewOffset(),
	}
}

func (h *RequestHandler) HandleRequest(writer http.ResponseWriter, request *http.Request) {

	o, _ := h.offset.Parse(request.URL.Query().Get("timepoint"))
	if o > 0 {
		h.processOffset(writer, o)
	} else {
		if h.processHttp(writer, request) {
			return
		}
	}
}

func (h *RequestHandler) processOffset(writer http.ResponseWriter, o int64) {
	// retrieve cached deltas for the given offset
	//TODO check offset is valid
	deltas, err := h.cacheService.Read(h.key, o)
	if err != nil {
		log.Fatal(err)
	}
	for _, delta := range deltas {
		_, _ = writer.Write([]byte(delta))
		writer.(http.Flusher).Flush()
		if h.wg != nil {
			h.wg.Done()
		}
	}
}

func (h *RequestHandler) processHttp(writer http.ResponseWriter, request *http.Request) bool {
	h.logger.InfoR(request, "User connected")
	subscription, _ := h.broker.Subscribe()
	for {
		select {
		case msg := <-subscription:
			_, _ = writer.Write([]byte(msg))
			writer.(http.Flusher).Flush()
			if h.wg != nil {
				h.wg.Done()
			}
		case <-request.Context().Done():
			_ = h.broker.Unsubscribe(subscription)
			h.logger.InfoR(request, "User disconnected")
			if h.wg != nil {
				h.wg.Done()
			}
			return true
		}
	}
	return false
}
