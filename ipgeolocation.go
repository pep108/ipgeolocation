package ipgeolocation

import (
	"errors"
	"net"
	"net/http"
	"time"
)

// IPGeolocationAPI Developr Documentation
// https://ipgeolocation.io/documentation/ip-geolocation-api.html

const Endpoint = "https://api.ipgeolocation.io/"

// New returns a pointer to a IPGeolocation config
func New(config Options) (*Client, error) {

	// Supported languages
	languages := []string{"en", "de", "ru", "ja", "fr", "cn", "es", "cs", "it"}

	opts := Options{
		APIKey:     config.APIKey,
		Endpoint:   Endpoint,
		HTTPClient: config.HTTPClient,
	}

	if opts.Language == "" {
		opts.Language = "en"
	} else {
		// Make sure the language is supported
		found := stringInArray(languages, opts.Language)
		if !found {
			return nil, errors.New("language not supported")
		}
	}

	// -----------------------------------------------
	// Initialize a Transport and http.Client for reuse.
	// The Client's Transport typically has internal state (cached TCP connections),
	// so Clients should be reused instead of created as needed.
	// Clients are safe for concurrent use by multiple goroutines.
	// references:
	// 	- https://golang.org/pkg/net/http/#Client
	// -----------------------------------------------
	if opts.HTTPClient == nil {
		var netTransport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   50 * time.Second,
				KeepAlive: 90 * time.Second,
			}).Dial,
			MaxIdleConnsPerHost:   2,
			MaxIdleConns:          20,
			IdleConnTimeout:       time.Duration(90) * time.Second,
			TLSHandshakeTimeout:   time.Duration(10) * time.Second,
			ExpectContinueTimeout: time.Duration(1) * time.Second,
			DisableKeepAlives:     true,
		}
		opts.HTTPClient = &http.Client{
			Timeout:   time.Second * 50,
			Transport: netTransport,
		}
	}

	client := &Client{
		options: opts,
	}

	return client, nil
}
