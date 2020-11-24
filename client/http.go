package client

import "net/http"

//creating new type HTTP client which is going to make requests to API and get a response
//* pointer to http.Client
type HTTPClient struct{
	client *http.Client
	BackendURI string
}
//constructor for a specific type
func NewHTTPClient(uri string) HTTPClient {
	return HTTPClient{
		BackendURI: uri,
		client: &http.Client{},
	}
}