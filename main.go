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

	service.NewCacheService(cacheConfiguration).WithTopic(filingHistoryStream).WithPath("/filings").Start()
	service.NewCacheService(cacheConfiguration).WithTopic(companyProfileStream).WithPath("/companies").Start()
	service.NewCacheService(cacheConfiguration).WithTopic(companyInsolvencyStream).WithPath("/insolvency-cases").Start()
	service.NewCacheService(cacheConfiguration).WithTopic(companyChargesStream).WithPath("/charges").Start()
	service.NewCacheService(cacheConfiguration).WithTopic(companyOfficersStream).WithPath("/officers").Start()
	service.NewCacheService(cacheConfiguration).WithTopic(companyPSCStream).WithPath("/persons-with-significant-control").Start()

	svc.Router().Path("/healthcheck").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	svc.Start()
}
