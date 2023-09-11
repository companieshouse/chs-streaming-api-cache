package config_test

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/companieshouse/chs-streaming-api-cache/config"
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// key constants
const (
	BINDADDRCONST             = `BIND_ADDRESS`
	CERTFILECONST             = `CERT_FILE`
	KEYFILECONST              = `KEY_FILE`
	CHSAPIKEYCONST            = `CHS_API_KEY`
	BACKENDURLCONST           = `STREAMING_BACKEND_URL`
	REDISURLCONST             = `REDIS_URL`
	REDISPOOLSIZECONST        = `REDIS_POOL_SIZE`
	CACHEEXPIRYINSECONDSCONST = `CACHE_EXPIRY_IN_SECONDS`
	STREAMFILINGSPATHCONST    = `STREAM_BACKEND_FILINGS_PATH`
	STREAMCOMPANIESPATHCONST  = `STREAM_BACKEND_COMPANIES_PATH`
	STREAMINSOLVENCYPATHCONST = `STREAM_BACKEND_INSOLVENCY_PATH`
	STREAMCHARGESPATHCONST    = `STREAM_BACKEND_CHARGES_PATH`
	STREAMOFFICERSPATHCONST   = `STREAM_BACKEND_OFFICERS_PATH`
	STREAMPSCSPATHCONST       = `STREAM_BACKEND_PSCS_PATH`
)

// value constants
const (
	bindAddrConst             = `bind-addr`
	certFileConst             = `cert-file`
	keyFileConst              = `key-file`
	chsApiKeyConst            = `chs-api-key`
	backEndUrlConst           = `streaming-backend-url`
	redisUrlConst             = `redis-url`
	redisPoolSizeConst        = 123
	cacheExpiryInSecondsConst = 456
	streamFilingsPathConst    = `stream-backend-filings-path`
	streamCompaniesPathConst  = `stream-backend-companies-path`
	streamInsolvencyPathConst = `stream-backend-insolvency-path`
	streamChargesPathConst    = `stream-backend-charges-path`
	streamOfficersPathConst   = `stream-backend-officers-path`
	streamPSCsPathConst       = `stream-backend-pscs-path`
)

func TestConfig(t *testing.T) {
	t.Parallel()
	os.Clearenv()
	var (
		err           error
		configuration *config.Config
		envVars       = map[string]string{
			BINDADDRCONST:             bindAddrConst,
			CERTFILECONST:             certFileConst,
			KEYFILECONST:              keyFileConst,
			CHSAPIKEYCONST:            chsApiKeyConst,
			BACKENDURLCONST:           backEndUrlConst,
			REDISURLCONST:             redisUrlConst,
			REDISPOOLSIZECONST:        strconv.Itoa(redisPoolSizeConst),
			CACHEEXPIRYINSECONDSCONST: strconv.Itoa(cacheExpiryInSecondsConst),
			STREAMFILINGSPATHCONST:    streamFilingsPathConst,
			STREAMCOMPANIESPATHCONST:  streamCompaniesPathConst,
			STREAMINSOLVENCYPATHCONST: streamInsolvencyPathConst,
			STREAMCHARGESPATHCONST:    streamChargesPathConst,
			STREAMOFFICERSPATHCONST:   streamOfficersPathConst,
			STREAMPSCSPATHCONST:       streamPSCsPathConst,
		}
		builtConfig = config.Config{
			BindAddress:          bindAddrConst,
			CertFile:             certFileConst,
			KeyFile:              keyFileConst,
			ChsApiKey:            chsApiKeyConst,
			BackEndUrl:           backEndUrlConst,
			RedisUrl:             redisUrlConst,
			RedisPoolSize:        redisPoolSizeConst,
			CacheExpiryInSeconds: cacheExpiryInSecondsConst,
			StreamFilingsPath:    streamFilingsPathConst,
			StreamCompaniesPath:  streamCompaniesPathConst,
			StreamInsolvencyPath: streamInsolvencyPathConst,
			StreamChargesPath:    streamChargesPathConst,
			StreamOfficersPath:   streamOfficersPathConst,
			StreamPSCsPath:       streamPSCsPathConst,
		}
		bindAddrRegex             = regexp.MustCompile(bindAddrConst)
		certFileRegex             = regexp.MustCompile(certFileConst)
		keyFileRegex              = regexp.MustCompile(keyFileConst)
		chsApiKeyRegex            = regexp.MustCompile(chsApiKeyConst)
		backEndUrlRegex           = regexp.MustCompile(backEndUrlConst)
		redisUrlRegex             = regexp.MustCompile(redisUrlConst)
		redisPoolSizeRegex        = regexp.MustCompile(strconv.Itoa(redisPoolSizeConst))
		cacheExpiryInSecondsRegex = regexp.MustCompile(strconv.Itoa(cacheExpiryInSecondsConst))
		streamFilingsPathRegex    = regexp.MustCompile(streamFilingsPathConst)
		streamCompaniesPathRegex  = regexp.MustCompile(streamCompaniesPathConst)
		streamInsolvencyPathRegex = regexp.MustCompile(streamInsolvencyPathConst)
		streamChargesPathRegex    = regexp.MustCompile(streamChargesPathConst)
		streamOfficersPathRegex   = regexp.MustCompile(streamOfficersPathConst)
		streamPSCsPathRegex       = regexp.MustCompile(streamPSCsPathConst)
	)

	// set test env variables
	for varName, varValue := range envVars {
		os.Setenv(varName, varValue)
		defer os.Unsetenv(varName)
	}

	Convey("Given an environment with no environment variables set", t, func() {

		Convey("Then configuration should be nil", func() {
			So(configuration, ShouldBeNil)
		})

		Convey("When the config values are retrieved", func() {

			Convey("Then there should be no error returned, and values are as expected", func() {
				configuration, err = config.Get()

				So(err, ShouldBeNil)
				So(configuration, ShouldResemble, &builtConfig)
			})

			Convey("The generated JSON string from configuration should not contain sensitive data", func() {
				jsonByte, err := json.Marshal(builtConfig)

				So(err, ShouldBeNil)
				So(bindAddrRegex.Match(jsonByte), ShouldEqual, true)
				So(certFileRegex.Match(jsonByte), ShouldEqual, false)
				So(keyFileRegex.Match(jsonByte), ShouldEqual, false)
				So(chsApiKeyRegex.Match(jsonByte), ShouldEqual, false)
				So(backEndUrlRegex.Match(jsonByte), ShouldEqual, true)
				So(redisUrlRegex.Match(jsonByte), ShouldEqual, true)
				So(redisPoolSizeRegex.Match(jsonByte), ShouldEqual, true)
				So(cacheExpiryInSecondsRegex.Match(jsonByte), ShouldEqual, true)
				So(streamFilingsPathRegex.Match(jsonByte), ShouldEqual, true)
				So(streamCompaniesPathRegex.Match(jsonByte), ShouldEqual, true)
				So(streamInsolvencyPathRegex.Match(jsonByte), ShouldEqual, true)
				So(streamChargesPathRegex.Match(jsonByte), ShouldEqual, true)
				So(streamOfficersPathRegex.Match(jsonByte), ShouldEqual, true)
				So(streamPSCsPathRegex.Match(jsonByte), ShouldEqual, true)
			})
		})
	})
}
