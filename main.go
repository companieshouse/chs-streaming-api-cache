package main

import (
	"github.com/companieshouse/chs-streaming-api-cache/config"
	"github.com/companieshouse/chs-streaming-api-cache/service"
	chslog "github.com/companieshouse/chs.go/log"
	chsservice "github.com/companieshouse/chs.go/service"
	"github.com/companieshouse/chs.go/service/handlers/requestID"
	"github.com/justinas/alice"
	"net/http"
)

const (
	filingHistoryStream     = "stream-filing-history"
	companyProfileStream    = "stream-company-profile"
	companyInsolvencyStream = "stream-company-insolvency"
	companyChargesStream    = "stream-company-charges"
	companyOfficersStream   = "stream-company-officers"
	companyPSCStream        = "stream-company-psc"
)

func main() {
	chsservice.DefaultMiddleware = []alice.Constructor{requestID.Handler(20), chslog.Handler}

	config, err := config.Get()
	if err != nil {
		panic(err)
	}
	svc := chsservice.New(config.ServiceConfig())

	cacheConfiguration := &service.CacheConfiguration{
		Configuration: config,
		Router:        svc.Router(),
	}

	service.NewCacheService(cacheConfiguration).WithTopic(filingHistoryStream).WithPath("/streaming-api-cache/filings").Initialise().Start()
	service.NewCacheService(cacheConfiguration).WithTopic(companyProfileStream).WithPath("/streaming-api-cache/companies").Initialise().Start()
	service.NewCacheService(cacheConfiguration).WithTopic(companyInsolvencyStream).WithPath("/streaming-api-cache/insolvency-cases").Initialise().Start()
	service.NewCacheService(cacheConfiguration).WithTopic(companyChargesStream).WithPath("/streaming-api-cache/charges").Initialise().Start()
	service.NewCacheService(cacheConfiguration).WithTopic(companyOfficersStream).WithPath("/streaming-api-cache/officers").Initialise().Start()
	service.NewCacheService(cacheConfiguration).WithTopic(companyPSCStream).WithPath("/streaming-api-cache/persons-with-significant-control").Initialise().Start()

	svc.Router().Path("/healthcheck").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	svc.Start()
}
