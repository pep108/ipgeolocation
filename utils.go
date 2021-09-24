package ipgeolocation

import (
	"errors"
	"net/url"
)

// validateParams is used to make sure required fields are present
func validateParams(params *IPGeolocationRequest) error {
	if params == nil {
		return errors.New("missing input")
	}

	if params.IP == "" && len(params.IPs) == 0 {
		return errors.New("missing IP")
	}

	return nil
}

// prepareUrlParams adds values to the url params from the request parameters
func (c *Client) prepareUrlParams(p *url.Values, params *IPGeolocationRequest) {
	// Add the apiKey
	p.Add("apiKey", c.options.APIKey)

	// Add the fields
	addUrlParms(p, "fields", params.Fields)

	// Set the language
	if params.Language == "" {
		params.Language = c.options.Language
	}
	p.Add("lang", params.Language)

	// Add security
	if params.IncludeSecurity {
		p.Add("include", "security")
	}

	// Add Includes
	if params.IncludeHostname {
		p.Add("include", "hostname")
	}
	if params.IncludeLiveHostname {
		p.Add("include", "liveHostname")
	}
	if params.IncludeHostnameFallbackLive {
		p.Add("include", "hostnameFallbackLive")
	}
	if params.IncludeUseragent {
		p.Add("include", "useragent")
	}
}

// addUrlParms adds array values to the url parameters
func addUrlParms(p *url.Values, field string, values []string) {
	for _, v := range values {
		p.Add(field, v)
	}
}

// stringInArray checks if a string is present in []string
func stringInArray(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}
