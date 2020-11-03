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
		StreamFilingsPath:    "/filings?timepoint=1",
		StreamCompaniesPath:  "/companies?timpoint=2",
		StreamInsolvencyPath: "/insolvencies",
		StreamChargesPath:    "/charges?timpoint=2",
		StreamOfficersPath:   "/officers?timepoint=1",
		StreamPSCsPath:       "/pscs",
	}
	pathMapper = *New(&cfg)

}

func TestUnitGetBackendPathForPathWithFilingHistoryReturnsMappedPath(t *testing.T) {
	initMapperTest()
	Convey("The provided path is for filing history", t, func() {
		Convey("The backend path is obtained for the path", func() {
			result, err := pathMapper.GetBackendPathForPath("/filings")
			Convey("The destination topic should be filing-history", func() {
				So(result, ShouldEqual, "/filings?timepoint=1")
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestUnitGetBackendPathForPathWithCompanyProfileReturnsMappedPath(t *testing.T) {
	initMapperTest()
	Convey("The provided path is for company profile", t, func() {
		Convey("The backend path is obtained for the path", func() {
			result, err := pathMapper.GetBackendPathForPath("/companies")
			Convey("The destination topic should be filing-history", func() {
				So(result, ShouldEqual, "/companies?timpoint=2")
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestUnitGetBackendPathForPathWithCompanyInsolvencyReturnsMappedPath(t *testing.T) {
	initMapperTest()
	Convey("The provided path is for company insolvency", t, func() {
		Convey("The backend path is obtained for the path", func() {
			result, err := pathMapper.GetBackendPathForPath("/insolvency-cases")
			Convey("The destination topic should be company insolvency", func() {
				So(result, ShouldEqual, "/insolvencies")
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestUnitGetBackendPathForPathWithCompanyChargesReturnsMappedPath(t *testing.T) {
	initMapperTest()
	Convey("The provided path is for company charges", t, func() {
		Convey("The backend path is obtained for the path", func() {
			result, err := pathMapper.GetBackendPathForPath("/charges")
			Convey("The destination topic should be company charges", func() {
				So(result, ShouldEqual, "/charges?timpoint=2")
				So(err, ShouldBeNil)
			})
		})
	})
}
func TestUnitGetBackendPathForPathWithCompanyOfficersReturnsMappedPath(t *testing.T) {
	initMapperTest()
	Convey("The provided path is for company officers", t, func() {
		Convey("The backend path is obtained for the path", func() {
			result, err := pathMapper.GetBackendPathForPath("/officers")
			Convey("The destination topic should be company officers", func() {
				So(result, ShouldEqual, "/officers?timepoint=1")
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestUnitGetBackendPathForPathWithCompanyPSCsReturnsMappedPath(t *testing.T) {
	initMapperTest()
	Convey("The provided path is for company officers", t, func() {
		Convey("The backend path is obtained for the path", func() {
			result, err := pathMapper.GetBackendPathForPath("/persons-with-significant-control")
			Convey("The destination topic should be company officers", func() {
				So(result, ShouldEqual, "/pscs")
				So(err, ShouldBeNil)
			})
		})
	})
}
