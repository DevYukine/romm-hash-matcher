package http

import "resty.dev/v3"

func SimpleHttpRetryCondition(resp *resty.Response, err error) bool {
	if err != nil {
		return true
	}

	// Retry on HTTP status codes 500-599 (Server Errors)
	if resp.StatusCode() >= 500 && resp.StatusCode() < 600 {
		return true
	}

	return false
}
