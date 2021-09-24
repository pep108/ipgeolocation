package ipgeolocation

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type Credentials struct {
	APIKey string `json:"APIKey"`
}

func TestGetIPGeolocation(t *testing.T) {
	var (
		output   IPGeolocationResult
		testIP   = "1.1.1.1"
		apiCreds Credentials
	)

	// Read the credentials from file
	creds, err := ioutil.ReadFile("key.json")
	if err != nil {
		t.Error("error reading API credentials")
		return
	}
	err = json.Unmarshal(creds, &apiCreds)
	if err != nil {
		t.Error("failed to load API credentials")
		return
	}

	// Configure the default client
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
	client := &http.Client{
		Timeout:   time.Second * 50,
		Transport: netTransport,
	}

	// Prepare the key for authentication
	p := url.Values{
		"apiKey": {apiCreds.APIKey},
		"ip":     {testIP},
		"fields": {"city"},
	}

	// Build the request
	t.Logf("url: %+v", Endpoint+"ipgeo?"+p.Encode())

	// Build the Request
	req, err := http.NewRequest("GET", Endpoint+"ipgeo?"+p.Encode(), nil)
	if err != nil {
		t.Errorf("error creating request: %v", err.Error())
		return
	}

	// ---------------------------------------------------------------
	// Make the Request
	// see: https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	// ---------------------------------------------------------------
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("error making request: %v", err.Error())

		// Make sure we close the body so we don't leak memory
		if resp != nil {
			resp.Body.Close()
		}

		return
	}

	// Always make sure the response body gets closed to avoid a memory leak
	defer func() {
		if resp != nil && resp.Body != nil {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("failed to close response body: %v", err.Error())
			}
		}
	}()

	t.Logf("StatusCode: %v", resp.StatusCode)

	if resp.StatusCode != 200 {
		t.Fail()
	}

	// --------------------------------------------------------------------------------
	// Use ioutil.ReadAll to ensure that we read the entire response body.
	// If we fail to do so, the connection will not be reused, the file
	// descriptor will remain and the application will leak resources.
	// ref: https://stackoverflow.com/questions/17948827/reusing-http-connections-in-golang
	// --------------------------------------------------------------------------------
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("error reading response body: %v", err.Error())
		return
	}

	t.Logf("--------- ipgeo response body ---------")
	t.Logf("%v", string(resBody))

	err = json.Unmarshal(resBody, &output)
	if err != nil {
		t.Errorf("error decoding response body: %v", err.Error())
		return
	}

	return
}

func TestBulkIPGeolocation(t *testing.T) {
	var (
		output   IPGeolocationResult
		testIP   = []string{"1.1.1.1", "8.8.8.8"}
		apiCreds Credentials
	)

	// Read the credentials from file
	creds, err := ioutil.ReadFile("key.json")
	if err != nil {
		t.Error("error reading API credentials")
		return
	}
	err = json.Unmarshal(creds, &apiCreds)
	if err != nil {
		t.Error("failed to load API credentials")
		return
	}

	// Configure the default client
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
	client := &http.Client{
		Timeout:   time.Second * 50,
		Transport: netTransport,
	}

	// Prepare the key for authentication
	p := url.Values{
		"apiKey": {apiCreds.APIKey},
		"fields": {"city"},
	}

	// Build the request
	t.Logf("url: %+v", Endpoint+"ipgeo-bulk?"+p.Encode())

	// Prepare the ip string
	ipbytes, err := json.Marshal(testIP)
	if err != nil {
		t.Errorf("failed to marshal testIP: %v", err.Error())
		return
	}

	// Set up the body of the POST request
	postBody := map[string]string{
		"ips": string(ipbytes),
	}
	reqBody, err := json.Marshal(postBody)
	if err != nil {
		t.Errorf("failed to marshal postBody: %v", err.Error())
		return
	}

	// Build the Request
	req, err := http.NewRequest("POST", Endpoint+"ipgeo-bulk?"+p.Encode(), bytes.NewReader(reqBody))
	if err != nil {
		t.Errorf("error creating request: %v", err.Error())
		return
	}

	// Make the Request
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("error making request: %v", err.Error())

		// Make sure we close the body so we don't leak memory
		if resp != nil {
			resp.Body.Close()
		}

		return
	}

	// Always make sure the response body gets closed to avoid a memory leak
	defer func() {
		if resp != nil && resp.Body != nil {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("failed to close response body: %v", err.Error())
			}
		}
	}()

	t.Logf("StatusCode: %v", resp.StatusCode)

	if resp.StatusCode != 200 {
		t.Fail()
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("error reading response body: %v", err.Error())
		return
	}

	t.Logf("--------- ipgeo response body ---------")
	t.Logf("%v", string(resBody))

	err = json.Unmarshal(resBody, &output)
	if err != nil {
		t.Errorf("error decoding response body: %v", err.Error())
		return
	}

	return
}
