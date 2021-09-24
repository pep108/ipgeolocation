package ipgeolocation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (c *Client) GetIPGeolocation(params *IPGeolocationRequest) (*IPGeolocationResult, error) {
	var (
		output IPGeolocationResult
		client = c.options.HTTPClient
	)

	// Validate input
	err := validateParams(params)
	if err != nil {
		return nil, err
	}
	if params.IP == "" {
		return nil, errors.New("missing ip")
	}

	// Prepare the key for authentication
	p := url.Values{
		"ip": {params.IP},
	}
	c.prepareUrlParams(&p, params)

	// Build the Request
	req, err := http.NewRequest("GET", c.options.Endpoint+"ipgeo?"+p.Encode(), nil)
	if err != nil {
		fmt.Printf("error creating request: %v", err.Error())
		return nil, err
	}

	// Make the Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error making request: %v", err.Error())

		// Make sure we close the body so we don't leak memory
		if resp != nil {
			resp.Body.Close()
		}

		return nil, err
	}

	// Always make sure the response body gets closed to avoid a memory leak
	defer func() {
		if resp != nil && resp.Body != nil {
			if err := resp.Body.Close(); err != nil {
				fmt.Printf("failed to close response body: %v", err.Error())
			}
		}
	}()

	// --------------------------------------------------------------------------------
	// Use ioutil.ReadAll to ensure that we read the entire response body.
	// If we fail to do so, the connection will not be reused, the file
	// descriptor will remain and the application will leak resources.
	// ref: https://stackoverflow.com/questions/17948827/reusing-http-connections-in-golang
	// --------------------------------------------------------------------------------
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body: %v", err.Error())
		return &output, err
	}

	err = json.Unmarshal(body, &output)
	if err != nil {
		fmt.Printf("error decoding response body: %v", err.Error())
		return &output, err
	}

	return &output, nil
}

func (c *Client) BulkIPGeolocation(params *IPGeolocationRequest) (*IPGeolocationResult, error) {
	var (
		output IPGeolocationResult
		client = c.options.HTTPClient
	)

	// Validate input
	err := validateParams(params)
	if err != nil {
		return nil, err
	}
	if len(params.IPs) == 0 {
		return nil, errors.New("missing ips")
	}

	// Prepare the key for authentication
	p := url.Values{}
	c.prepareUrlParams(&p, params)

	// Prepare the ip string
	ipbytes, err := json.Marshal(params.IPs)
	if err != nil {
		return nil, err
	}

	// Set up the body of the POST request
	postBody := map[string]string{
		"ips": string(ipbytes),
	}
	reqBody, err := json.Marshal(postBody)
	if err != nil {
		fmt.Printf("failed to marshal postBody: %v", err.Error())
		return nil, err
	}

	// Build the Request
	req, err := http.NewRequest("POST", c.options.Endpoint+"ipgeo-bulk?"+p.Encode(), bytes.NewReader(reqBody))
	if err != nil {
		fmt.Printf("error creating request: %v", err.Error())
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	// Make the Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error making request: %v", err.Error())

		// Make sure we close the body so we don't leak memory
		if resp != nil {
			resp.Body.Close()
		}

		return nil, err
	}

	// Always make sure the response body gets closed to avoid a memory leak
	defer func() {
		if resp != nil && resp.Body != nil {
			if err := resp.Body.Close(); err != nil {
				fmt.Printf("failed to close response body: %v", err.Error())
			}
		}
	}()

	// --------------------------------------------------------------------------------
	// Use ioutil.ReadAll to ensure that we read the entire response body.
	// If we fail to do so, the connection will not be reused, the file
	// descriptor will remain and the application will leak resources.
	// ref: https://stackoverflow.com/questions/17948827/reusing-http-connections-in-golang
	// --------------------------------------------------------------------------------
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body: %v", err.Error())
		return &output, err
	}

	err = json.Unmarshal(respBody, &output)
	if err != nil {
		fmt.Printf("error decoding response body: %v", err.Error())
		return &output, err
	}

	return &output, nil
}
