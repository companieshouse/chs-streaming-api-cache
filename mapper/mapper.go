package mapper

import (
	"fmt"
	"github.com/companieshouse/chs-streaming-api-cache/config"
)

const (
	servicePrefix     = "/streaming-api-cache"
	FilingHistoryPath = servicePrefix + "/filings"
	CompaniesPath     = servicePrefix + "/companies"
	InsolvencyPath    = servicePrefix + "/insolvency-cases"
	ChargesPath       = servicePrefix + "/charges"
	OfficersPath      = servicePrefix + "/officers"
	PSCsPath          = servicePrefix + "/persons-with-significant-control"
)

// A topic mapper that obtains topics for the specified resource kind from the app configuration model
type ConfigurationPathMapper struct {
	Paths        map[string]string
	DefaultTopic string
}

// Create a new ConfigurationPathMapper instance with all backend path mappings resolved from configuration
func New(cfg *config.Config) *ConfigurationPathMapper {
	mapper := &ConfigurationPathMapper{}
	mapper.Paths = map[string]string{
		FilingHistoryPath: cfg.StreamFilingsPath,
		CompaniesPath:     cfg.StreamCompaniesPath,
		InsolvencyPath:    cfg.StreamInsolvencyPath,
		ChargesPath:       cfg.StreamChargesPath,
		OfficersPath:      cfg.StreamOfficersPath,
		PSCsPath:          cfg.StreamPSCsPath,
	}
	return mapper
}

// Obtain a backend path corresponding to the given stream path
func (mapper *ConfigurationPathMapper) GetBackendPathForPath(path string) (string, error) {
	val, exists := mapper.Paths[path]
	if exists {
		return val, nil
	}
	return "", fmt.Errorf("resource path [%s] unhandled", path)
}
