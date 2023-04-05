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

type InvalidShortlinkDataError struct {
	Title           string   `json:"title,omitempty"`
	Type            string   `json:"type,omitempty"`
	Detail          string   `json:"detail,omitempty"`
	Status          int      `json:"status,omitempty"`
	InvalidElements []string `json:"invalidElements,omitempty"`
}

type ShortcodeNotFoundError struct {
	Detail    string `json:"detail"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Status    int    `json:"status"`
	ShortCode string `json:"shortCode"`
}

type UnknownError struct {
	Type   string `json:"type,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Status int    `json:"status,omitempty"`
}

type ModifyShortlinkRequest struct {
	LongUrl        string         `json:"longUrl"`
	DeviceLongUrls DeviceLongURLs `json:"deviceLongUrls"`
	ValidSince     string         `json:"validSince"`
	ValidUntil     string         `json:"validUntil"`
	MaxVisits      int            `json:"maxVisits"`
	Tags           []string       `json:"tags"`
	Title          string         `json:"title"`
	Crawlable      bool           `json:"crawlable"`
	ForwardQuery   bool           `json:"forwardQuery"`
}

type CannotDeleteShortlink struct {
	Title     string `json:"title"`
	Type      string `json:"type"`
	Detail    string `json:"detail"`
	Status    int    `json:"status"`
	ShortCode string `json:"shortCode"`
	Threshold int    `json:"threshold"`
}

func handleCommonErrors(respBytes []byte, statusCode int) error {
	if statusCode == 200 || statusCode == 201 {
		return nil
	}

	if statusCode == 404 {
		var modifyErr ShortcodeNotFoundError
		err := json.Unmarshal(respBytes, &modifyErr)
		if err != nil {
			return err
		}
		return fmt.Errorf("error occurred with shortlink: not found")
	}

	if statusCode == 400 {
		var modifyError InvalidShortlinkDataError
		err := json.Unmarshal(respBytes, &modifyError)
		if err != nil {
			return err
		}
		return fmt.Errorf("error occurred with shortlink: [%s] invalid elements: %s", modifyError.Title, strings.Join(modifyError.InvalidElements, ","))
	}

	if statusCode == 422 {
		var modifyError CannotDeleteShortlink
		err := json.Unmarshal(respBytes, &modifyError)
		if err != nil {
			return err
		}
		return fmt.Errorf("error occurred with shortlink: cannot delete threshold is too high %d", modifyError.Threshold)
	}

	var unknownError UnknownError
	err := json.Unmarshal(respBytes, &unknownError)
	if err != nil {
		return err
	}
	return fmt.Errorf("unknown api error: %s", unknownError.Detail)
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

	return nil, handleCommonErrors(respBytes, resp.StatusCode)
}

func (c *Client) GetShortlink(shortCode string) (*ShortLink, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/rest/v3/short-urls/%s", shortCode), nil)
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

	return nil, handleCommonErrors(respBytes, resp.StatusCode)
}

func (c *Client) UpdateShortlink(shortCode string, request *ModifyShortlinkRequest) (*ShortLink, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("PATCH", fmt.Sprintf("/rest/v3/short-urls/%s", shortCode), &buf)
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

	return nil, handleCommonErrors(respBytes, resp.StatusCode)
}

func (c *Client) DeleteShortlink(shortCode string, request *ModifyShortlinkRequest) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}

	resp, err := c.doRequest("DELETE", fmt.Sprintf("/rest/v3/short-urls/%s", shortCode), &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == 204 {
		return nil
	}

	return handleCommonErrors(respBytes, resp.StatusCode)
}
