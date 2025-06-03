package hasheous

import (
	"fmt"
	"resty.dev/v3"
	"romm-hash-matcher/internal/http"
	"romm-hash-matcher/internal/model"
	"romm-hash-matcher/internal/ratelimit"
	"time"
)

const baseUrl = "https://hasheous.org/api/v1"

type Client struct {
	client *resty.Client
}

func NewClient() *Client {
	return &Client{
		client: resty.New().SetHeader("User-Agent", http.DefaultUserAgent).SetBaseURL(baseUrl),
	}
}

func (c Client) LookupByHash(rom model.InternalRom) (*LookupResponse, error) {
	ratelimit.HasheousRatelimit.Take()

	resp, err := c.client.R().
		AddRetryConditions(http.SimpleHttpRetryCondition).
		SetRetryCount(5).
		SetRetryWaitTime(5 * time.Second).
		SetBody(&LookupRequest{MD5: rom.MD5, ShA1: rom.SHA1}).
		SetResult(&LookupResponse{}).
		Post("/Lookup/ByHash")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == 404 {
		return nil, nil
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status())
	}

	return resp.Result().(*LookupResponse), nil

}
