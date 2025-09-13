package restclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type restyClient struct {
	client        *resty.Client
	sslEnabled    bool
	retryCount    int
	retryWaitTime time.Duration
	timeout       time.Duration
}

// NewHTTPClient ...
func newRestyClient(sslEnabled bool) RestClientItf {
	return &restyClient{
		client:        resty.New(),
		sslEnabled:    sslEnabled,
		retryCount:    3,
		timeout:       5 * time.Second,
		retryWaitTime: 500 * time.Millisecond,
	}
}

func (p *restyClient) InvokeRequest(c context.Context, method string, url string, body interface{}, headers map[string]string, auth map[string]string, queryparams map[string]string, pathparams map[string]string, formData map[string]string) ([]byte, int, error) {
	var (
		request    *resty.Request
		response   *resty.Response
		resp       []byte
		err        error
		statusCode int
	)

	/* if timeout > 0 {
		request = resty.New().
			SetTimeout(time.Duration(timeout) * time.Millisecond).
			SetRetryCount(p.retryCount).
			SetRetryWaitTime(time.Duration(time.Duration(p.retryWaitTime)) * time.Millisecond).
			// Default (nil) implies exponential backoff with jitter
			SetRetryAfter(nil).
			R()
	} else  */{
		request = resty.New().
			SetTimeout(p.timeout * time.Millisecond).
			SetRetryCount(p.retryCount).
			SetRetryWaitTime(time.Duration(time.Duration(p.retryWaitTime)) * time.Millisecond).
			// Default (nil) implies exponential backoff with jitter
			SetRetryAfter(nil).
			R()
	}

	request = request.SetHeader("Content-Type", "application/json")
	//set authorization if available
	if len(auth) > 0 {
		if auth["username"] != "" {
			request = request.SetBasicAuth(auth["username"], auth["password"])
		} else if auth["token"] != "" {
			request = request.SetAuthToken(auth["token"])
		}
	}
	//add pathparams if available
	if len(pathparams) > 0 {
		for key, value := range pathparams {
			request = request.SetPathParam(key, value)
		}
	}
	//add queryparams
	if queryparams != nil {
		request = request.SetQueryParams(queryparams)
	}

	//set custom headers
	if len(headers) > 0 {
		for key, val := range headers {
			if len(strings.ReplaceAll(key, " ", "")) > 0 {
				request = request.SetHeader(key, val)
			}
		}
	}

	//Set form data
	if len(formData) > 0 {
		request = request.SetFormData(formData)
	}
	if body != "" {
		request = request.SetBody(body)
	}

	switch method {
	case resty.MethodGet:
		response, err = request.Get(url)
	case resty.MethodPost:
		response, err = request.Post(url)
	case resty.MethodPut:
		response, err = request.Put(url)
	case resty.MethodDelete:
		response, err = request.Delete(url)
	case resty.MethodPatch:
		response, err = request.Patch(url)
	default:
		err = fmt.Errorf("unsupported request method: %v", method)
	}

	if err == nil {
		resp = response.Body()
	}
	statusCode = response.StatusCode()

	return resp, statusCode, err
}
