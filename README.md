The unofficial client library for accessing IP Geolocation APIs from Go. \
https://ipgeolocation.io

## Client Initialization

```go
ipgeo := ipgeolocation.New(ipgeolocation.Options{
	APIKey:     "API_KEY",
})
```

## API Examples

### Single IP Geolocation Lookup API

```go
func main() {
	// Creates an ipgeolocation instance with default http.Client
	ipgeo := ipgeolocation.New(ipgeolocation.Options{
		APIKey:     "API_KEY",
		Language:   "en",
	})

	// Make the request
	res, err := ipgeo.GetIPGeolocation(&IPGeolocationRequest{
		IP: "1.1.1.1",
		Fields: []{"geo"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Do something with the result
	fmt.Printf("result: %+v", res)

}
```

### Bulk IP Geolocation Lookup API

*paid subscription required*

```go
func main() {
	// Creates an ipgeolocation instance with default http.Client
	ipgeo := ipgeolocation.New(ipgeolocation.Options{
		APIKey:     "API_KEY",
	})

	// Make the request
	res, err := ipgeo.BulkIPGeolocation(&IPGeolocationRequest{
		IPs: []string{"1.1.1.1", "8.8.8.8"},
		Fields: []{"city", "state_prov", "zipcode"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Do something with the result
	fmt.Printf("result: %+v", res)
}
```