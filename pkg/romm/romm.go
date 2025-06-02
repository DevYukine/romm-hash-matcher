package romm

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"github.com/zhyee/zipstream"
	"go.uber.org/zap"
	"io"
	"resty.dev/v3"
	"romm-hash-matcher/internal/logging"
	"romm-hash-matcher/internal/model"
	"romm-hash-matcher/internal/ratelimit"
	"strconv"
	"strings"
)

type Client struct {
	client  *resty.Client
	baseUrl string
}

func NewClient(baseUrl string, username string, password string) *Client {
	return &Client{
		client: resty.New().SetBasicAuth(username, password).SetHeader("User-Agent", "romm-hash-matcher/1.0").SetBaseURL(baseUrl).SetDisableWarn(true),
	}
}

func (c Client) GetMetadataOfZippedRom(url string) (*model.InternalRom, error) {
	ratelimit.RommRatelimit.Take()

	resp, err := c.client.R().Get(url)
	if err != nil {
		logging.Logger.Fatal("Error fetching URL", zap.Error(err))
	}
	defer resp.Body.Close()

	if !resp.IsSuccess() {
		logging.Logger.Fatal("Unexpected HTTP status: %s", zap.Int("status", resp.StatusCode()))
	}

	zr := zipstream.NewReader(resp.Body)

	var fileCount int
	// Prepare hash writers
	md5HashWriter := md5.New()
	sha1HashWriter := sha1.New()
	sha256HashWriter := sha256.New()
	multiWriter := io.MultiWriter(md5HashWriter, sha1HashWriter, sha256HashWriter)

	var name string
	var size int64

	// Iterate entries in the ZIP stream
	for {
		entry, err := zr.GetNextEntry()
		if err == io.EOF {
			break
		}
		if err != nil {
			logging.Logger.Fatal("Error reading ZIP stream", zap.Error(err))
		}

		// Skip directories and nested files
		if entry.FileInfo().IsDir() || strings.Contains(entry.Name, "/") {
			continue
		}

		fileCount++
		if fileCount > 1 {
			logging.Logger.Warn("found more than one top-level file found in ZIP archive, skipping")
			return nil, nil
		}

		// Stream the single top-level file into hash writers
		reader, err := entry.Open()

		if err != nil {
			logging.Logger.Fatal("Error opening entry", zap.String("Name", entry.Name), zap.Error(err))
		}

		if _, err := io.Copy(multiWriter, reader); err != nil {
			logging.Logger.Fatal("Error hashing file '%s': %v", zap.String("Name", entry.Name), zap.Error(err))
		}

		name = entry.Name
		size = entry.FileInfo().Size()
	}

	if fileCount == 0 {
		logging.Logger.Fatal("Error: no top-level file found in ZIP archive")
	}

	md5Hash := fmt.Sprintf("%x", md5HashWriter.Sum(nil))
	sha1Hash := fmt.Sprintf("%x", sha1HashWriter.Sum(nil))
	sha256Hash := fmt.Sprintf("%x", sha256HashWriter.Sum(nil))

	return &model.InternalRom{
		Name:   name,
		Size:   size,
		MD5:    md5Hash,
		SHA1:   sha1Hash,
		SHA256: &sha256Hash,
	}, nil
}

func (c Client) ManuallyMatchRom(rommId int64, fsName string, searchResult *SearchResponse) error {
	ratelimit.RommRatelimit.Take()

	res, err := c.client.R().
		SetQueryParam("remove_cover", "false").
		SetQueryParam("unmatch_metadata", "false").
		SetFormData(map[string]string{
			"igdb_id":   strconv.FormatInt(*searchResult.IgdbId, 10),
			"name":      searchResult.Name,
			"fs_name":   fsName,
			"summary":   searchResult.Summary,
			"url_cover": searchResult.IgdbUrlCover,
		}).
		Put("api/roms/" + strconv.FormatInt(rommId, 10))

	if err != nil {
		return err
	}

	if !res.IsSuccess() {
		logging.Logger.Error("Unexpected HTTP status", zap.Int("status", res.StatusCode()))
	}

	logging.Logger.Debug("Manually matched ROM", zap.Int64("rommId", rommId), zap.Int64p("igdbId", searchResult.IgdbId))
	return nil
}

func (c Client) SearchMetadataByIgdbId(rommId int64, igdbId int64) (*SearchResponse, error) {
	ratelimit.RommRatelimit.Take()

	res, err := c.client.R().
		SetQueryParam("rom_id", strconv.FormatInt(rommId, 10)).
		SetQueryParam("search_term", strconv.FormatInt(igdbId, 10)).
		SetQueryParam("search_by", "Id").
		SetResult(&[]SearchResponse{}).
		Get("api/search/roms")

	if err != nil {
		return nil, err
	}

	if !res.IsSuccess() {
		logging.Logger.Error("Unexpected HTTP status", zap.Int("status", res.StatusCode()))
	}

	result := res.Result().(*[]SearchResponse)

	for _, searchResponse := range *result {
		if *searchResponse.IgdbId == igdbId {
			logging.Logger.Debug("Found matching IGDB ID", zap.Int64("rommId", rommId), zap.Int64("igdbId", igdbId))
			return &searchResponse, nil
		}
	}

	logging.Logger.Debug("No matching IGDB ID found", zap.Int64("rommId", rommId), zap.Int64("igdbId", igdbId))
	return nil, nil
}

func (c Client) GetUnmatchedRoms() (*[]Item, error) {
	ratelimit.RommRatelimit.Take()

	hasNextPage := true
	offset := 0

	romItems := make([]Item, 0)

	for hasNextPage {
		res, err := c.client.R().
			SetQueryParam("search_term", "").
			SetQueryParam("limit", "100").
			SetQueryParam("offset", strconv.FormatInt(int64(offset), 10)).
			SetQueryParam("order_by", "name").
			SetQueryParam("order_dir", "asc").
			SetQueryParam("unmatched_only", "true").
			SetQueryParam("matched_only", "false").
			SetQueryParam("favourites_only", "false").
			SetQueryParam("duplicates_only", "false").
			SetQueryParam("group_by_meta_id", "false").
			SetResult(&Rom{}).
			Get("api/roms")

		logging.Logger.Debug("requesting unmatched ROMs", zap.String("url", fmt.Sprintf("%s/api/roms", c.baseUrl)), zap.String("offset", strconv.FormatInt(int64(offset), 10)))

		if err != nil {
			logging.Logger.Error("Error fetching unmatched ROMs", zap.Error(err))
			return nil, err
		}

		result := res.Result().(*Rom)

		logging.Logger.Debug("Fetched page", zap.Int("offset", result.Offset), zap.Int("limit", result.Limit), zap.Int("total", result.Total), zap.Int("count", len(result.Items)))

		romItems = append(romItems, result.Items...)

		hasNextPage = (result.Offset + len(result.Items)) < result.Total
		offset += len(result.Items)
	}

	return &romItems, nil
}
