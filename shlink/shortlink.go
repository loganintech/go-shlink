package shlink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type DeviceLongURLs struct {
	Android *string `json:"android,omitempty"`
	Ios     *string `json:"ios,omitempty"`
	Desktop *string `json:"desktop,omitempty"`
}

type CreateShortlinkRequest struct {
	LongUrl         string         `json:"longUrl,omitempty"`
	DeviceLongURLs  DeviceLongURLs `json:"deviceLongUrls,omitempty"`
	ValidSince      *time.Time     `json:"validSince,omitempty"`
	ValidUntil      *time.Time     `json:"validUntil,omitempty"`
	MaxVisits       int            `json:"maxVisits,omitempty"`
	Tags            []string       `json:"tags,omitempty"`
	Title           string         `json:"title,omitempty"`
	Crawlable       bool           `json:"crawlable,omitempty"`
	ForwardQuery    bool           `json:"forwardQuery,omitempty"`
	CustomSlug      string         `json:"customSlug,omitempty"`
	FindIfExists    bool           `json:"findIfExists,omitempty"`
	Domain          string         `json:"domain,omitempty"`
	ShortCodeLength int            `json:"shortCodeLength,omitempty"`
}

type ShortLinkVisitsSummary struct {
	Total   int `json:"total,omitempty"`
	NonBots int `json:"nonBots,omitempty"`
	Bots    int `json:"bots,omitempty"`
}

type ShortLinkMeta struct {
	ValidSince *time.Time `json:"validSince,omitempty"`
	ValidUntil *time.Time `json:"validUntil,omitempty"`
	MaxVisits  int        `json:"maxVisits,omitempty"`
}

type ShortLink struct {
	ShortCode      string                 `json:"shortCode,omitempty"`
	ShortUrl       string                 `json:"shortUrl,omitempty"`
	LongUrl        string                 `json:"longUrl,omitempty"`
	DeviceLongURLs DeviceLongURLs         `json:"deviceLongUrls,omitempty"`
	DateCreated    *time.Time             `json:"dateCreated,omitempty"`
	VisitsSummary  ShortLinkVisitsSummary `json:"visitsSummary,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Meta           ShortLinkMeta          `json:"meta,omitempty"`
	Domain         string                 `json:"domain,omitempty"`
	Title          string                 `json:"title,omitempty"`
	Crawlable      bool                   `json:"crawlable,omitempty"`
	VisitsCount    int64                  `json:"visitsCount,omitempty"`
}

func (c *Client) CreateShortlink(request *CreateShortlinkRequest) (*ShortLink, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("POST", "/rest/v3/short-urls", &buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		var shortLink ShortLink

		err = json.Unmarshal(respBytes, &shortLink)
		if err != nil {
			return nil, err
		}
		return &shortLink, nil
	}

	if resp.StatusCode == 400 {
		var createErr CreateShortlinkError
		err = json.Unmarshal(respBytes, &createErr)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error occurred creating shortlink: [%s] invalid elements: %s", createErr.Title, strings.Join(createErr.InvalidElements, ","))
	}

	var unknownError UnknownError
	err = json.Unmarshal(respBytes, &unknownError)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("unknown api error: %s", unknownError.Detail)
}

type CreateShortlinkError struct {
	Title           string   `json:"title,omitempty"`
	Type            string   `json:"type,omitempty"`
	Detail          string   `json:"detail,omitempty"`
	Status          int      `json:"status,omitempty"`
	InvalidElements []string `json:"invalidElements,omitempty"`
}

type UnknownError struct {
	Type   string `json:"type,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Status int    `json:"status,omitempty"`
}
