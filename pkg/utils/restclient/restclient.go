package restclient

import (
	"context"
)

type RestClientItf interface {
	InvokeRequest(c context.Context, method string, url string, body interface{}, headers map[string]string, auth map[string]string, queryparams map[string]string, pathparams map[string]string, formData map[string]string) ([]byte, int, error)
}

func GetRestClientItf(sslEnabled bool) RestClientItf {
	return newRestyClient(sslEnabled)
}
