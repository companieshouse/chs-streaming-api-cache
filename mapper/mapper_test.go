package mapper

import (
	"github.com/companieshouse/chs-streaming-api-cache/config"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	pathMapper ConfigurationPathMapper
	cfg        config.Config
)

func initMapperTest() {
	cfg = config.Config{
		StreamFilingsPath:    "/backend/filings?timepoint=1",
		StreamCompaniesPath:  "/backend/companies?timpoint=2",
		StreamInsolvencyPath: "/backend/insolvencies",
		StreamChargesPath:    "/backend/charges?timpoint=2",
		StreamOfficersPath:   "/backend/officers?timepoint=1",
		StreamPSCsPath:       "/backend/pscs",
	}
	pathMapper = *New(&cfg)

}

func TestUnitGetBackendPathForPathWithFilingHistoryReturnsMappedPath(t *testing.T) {
	initMapperTest()
	Convey("The provided path is for filing history", t, func() {
		Convey("The backend path is obtained for the path", func() {
			result, err := pathMapper.GetBackendPathForPath("/streaming-api-cache/filings")
			Convey("The destination topic should be filing-history", func() {
				So(result, ShouldEqual, "/backend/filings?timepoint=1")
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestUnitGetBackendPathForPathWithCompanyProfileReturnsMappedPath(t *testing.T) {
	initMapperTest()
	Convey("The provided path is for company profile", t, func() {
		Convey("The backend path is obtained for the path", func() {
			result, err := pathMapper.GetBackendPathForPath("/streaming-api-cache/companies")
			Convey("The destination topic should be filing-history", func() {
				So(result, ShouldEqual, "/backend/companies?timpoint=2")
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestUnitGetBackendPathForPathWithCompanyInsolvencyReturnsMappedPath(t *testing.T) {
	initMapperTest()
	Convey("The provided path is for company insolvency", t, func() {
		Convey("The backend path is obtained for the path", func() {
			result, err := pathMapper.GetBackendPathForPath("/streaming-api-cache/insolvency-cases")
			Convey("The destination topic should be company insolvency", func() {
				So(result, ShouldEqual, "/backend/insolvencies")
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestUnitGetBackendPathForPathWithCompanyChargesReturnsMappedPath(t *testing.T) {
	initMapperTest()
	Convey("The provided path is for company charges", t, func() {
		Convey("The backend path is obtained for the path", func() {
			result, err := pathMapper.GetBackendPathForPath("/streaming-api-cache/charges")
			Convey("The destination topic should be company charges", func() {
				So(result, ShouldEqual, "/backend/charges?timpoint=2")
				So(err, ShouldBeNil)
			})
		})
	})
}
func TestUnitGetBackendPathForPathWithCompanyOfficersReturnsMappedPath(t *testing.T) {
	initMapperTest()
	Convey("The provided path is for company officers", t, func() {
		Convey("The backend path is obtained for the path", func() {
			result, err := pathMapper.GetBackendPathForPath("/streaming-api-cache/officers")
			Convey("The destination topic should be company officers", func() {
				So(result, ShouldEqual, "/backend/officers?timepoint=1")
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestUnitGetBackendPathForPathWithCompanyPSCsReturnsMappedPath(t *testing.T) {
	initMapperTest()
	Convey("The provided path is for company officers", t, func() {
		Convey("The backend path is obtained for the path", func() {
			result, err := pathMapper.GetBackendPathForPath("/streaming-api-cache/persons-with-significant-control")
			Convey("The destination topic should be company officers", func() {
				So(result, ShouldEqual, "/backend/pscs")
				So(err, ShouldBeNil)
			})
		})
	})
}
