package handlers

import (
	. "github.com/companieshouse/chs-streaming-api-cache/cache"
	"github.com/companieshouse/chs-streaming-api-cache/logger"
	"github.com/companieshouse/chs-streaming-api-cache/offset"
	"github.com/companieshouse/chs.go/log"
	"net/http"
	"sync"
)

type Subscribable interface {
	Subscribe() (chan string, error)
	Unsubscribe(chan string) error
}

type RequestHandler struct {
	broker       Subscribable
	cacheService Cacheable
	key          string
	logger       logger.Logger
	offset       offset.Interface
	wg           *sync.WaitGroup
}

func NewRequestHandler(broker Subscribable, cacheService Cacheable, logger logger.Logger, topic string) *RequestHandler {
	return &RequestHandler{
		broker:       broker,
		logger:       logger,
		cacheService: cacheService,
		key:          topic,
		offset:       offset.NewOffset(),
	}
}

func (h *RequestHandler) HandleRequest(writer http.ResponseWriter, request *http.Request) {

	offset := request.URL.Query().Get("timepoint")
	o, err := h.offset.Parse(offset)
	if err != nil {
		h.logger.Info("Invalid offset requested", log.Data{"error": err, "offset":offset })
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	h.logger.Info("Retrieved offset from the url", log.Data{"timepoint": o, "topic": h.key})

	if o > 0 {
		h.processOffset(writer, o)
	}
	h.processHttp(writer, request)
	return
}

func (h *RequestHandler) processOffset(writer http.ResponseWriter, o int64) {
	//TODO check offset is valid
	h.logger.Info(" Retrieving cached deltas for the given offset", log.Data{"timepoint": o, "topic": h.key})
	deltas, err := h.cacheService.Read(h.key, o)
	if err != nil {
		h.logger.Error(err, log.Data{"timepoint": o, "topic": h.key})
		return
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
