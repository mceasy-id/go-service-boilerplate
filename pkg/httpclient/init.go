package httpclient

import "github.com/hashicorp/go-retryablehttp"

func NewWithoutLog() *retryablehttp.Client {
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil

	return httpClient
}
