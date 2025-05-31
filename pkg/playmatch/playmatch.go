package playmatch

import (
	"go.uber.org/zap"
	"resty.dev/v3"
	"romm-hash-matcher/internal/logging"
	"romm-hash-matcher/internal/model"
	"romm-hash-matcher/internal/ratelimit"
	"strconv"
)

const baseUrl = "https://playmatch.retrorealm.dev/api"

type Client struct {
	client *resty.Client
}

func NewClient() *Client {
	return &Client{
		client: resty.New().SetBaseURL(baseUrl).SetHeader("User-Agent", "romm-hash-matcher/1.0"),
	}
}

func (c Client) GetIdBySlug(slug string) (*int64, error) {
	ratelimit.PlaymatchRatelimit.Take()

	resp, err := c.client.R().
		SetQueryParam("slug", slug).
		SetResult(&SlimmedGameResponse{}).
		Get("/igdb/game")

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		logging.Logger.Error("Unexpected HTTP status", zap.Int("status", resp.StatusCode()))
	}

	matchResponse := resp.Result().(*SlimmedGameResponse)

	return &matchResponse.Id, nil
}

func (c Client) IdentifyGame(rom *model.InternalRom) (*MatchResponse, error) {
	ratelimit.PlaymatchRatelimit.Take()

	resp, err := c.client.R().
		SetQueryParam("fileName", rom.Name).
		SetQueryParam("fileSize", strconv.FormatInt(rom.Size, 10)).
		SetQueryParam("md5", rom.MD5).
		SetQueryParam("sha1", rom.SHA1).
		SetQueryParam("sha256", rom.SHA256).
		SetResult(&MatchResponse{}).
		Get("/identify/ids")

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		logging.Logger.Error("Unexpected HTTP status", zap.Int("status", resp.StatusCode()))
	}

	matchResponse := resp.Result().(*MatchResponse)

	return matchResponse, nil
}
