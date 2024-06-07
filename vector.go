package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	vectorSuffix      = "/vector_stores"
	vectorFilesSuffix = "/files"
)

type Vector struct {
	ID         string      `json:"id"`
	Object     string      `json:"object"`
	CreatedAt  int64       `json:"created_at"`
	Name       *string     `json:"name,omitempty"`
	Bytes      int64       `json:"bytes"`
	FileCounts *FileCounts `json:"file_counts,omitempty"`
	httpHeader
}

type FileCounts struct {
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
	Failed     int `json:"failed"`
	Cancelled  int `json:"cancelled"`
	Total      int `json:"total"`
}

type VectorRequest struct {
	Name    *string   `json:"name,omitempty"`
	FileIDs *[]string `json:"file_ids,omitempty"`
}

// MarshalJSON provides a custom marshaller for the assistant request to handle the API use cases
// If Tools is nil, the field is omitted from the JSON.
// If Tools is an empty slice, it's included in the JSON as an empty array ([]).
// If Tools is populated, it's included in the JSON with the elements.
func (a VectorRequest) MarshalJSON() ([]byte, error) {

	return json.Marshal(a)
}

// AssistantsList is a list of assistants.
type VectorList struct {
	Vectors []Vector `json:"data"`
	LastID  *string  `json:"last_id"`
	FirstID *string  `json:"first_id"`
	HasMore bool     `json:"has_more"`
	httpHeader
}

type VectorDeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`

	httpHeader
}

type VectorFile struct {
	ID            string `json:"id"`
	Object        string `json:"object"`
	CreatedAt     int64  `json:"created_at"`
	UsageBytes    int64  `json:"usage_bytes"`
	VectorStoreID string `json:"vector_store_id"`
	Status        string `json:"status"`
	LastError     string `json:"last_error"`

	httpHeader
}

type VectorFileFileRequest struct {
	VectorStoreID string `json:"vector_store_id"`
	FileID        string `json:"file_id"`
}

type VectorFilesList struct {
	VectorFiles []VectorFile `json:"data"`

	httpHeader
}

// CreateVector creates a new vector.
func (c *Client) CreateVector(ctx context.Context, request VectorRequest) (response Vector, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL(vectorSuffix), withBody(request),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

// RetrieveAssistant retrieves an assistant.
func (c *Client) RetrieveVector(
	ctx context.Context,
	vectorID string,
) (response Vector, err error) {
	urlSuffix := fmt.Sprintf("%s/%s", vectorSuffix, vectorID)
	req, err := c.newRequest(ctx, http.MethodGet, c.fullURL(urlSuffix),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

// ModifyVector modifies an assistant.
func (c *Client) ModifyVector(
	ctx context.Context,
	vectorID string,
	request VectorRequest,
) (response Vector, err error) {
	urlSuffix := fmt.Sprintf("%s/%s", vectorSuffix, vectorID)
	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL(urlSuffix), withBody(request),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

// DeleteVector deletes an assistant.
func (c *Client) DeleteVector(
	ctx context.Context,
	vectorID string,
) (response VectorDeleteResponse, err error) {
	urlSuffix := fmt.Sprintf("%s/%s", vectorSuffix, vectorID)
	req, err := c.newRequest(ctx, http.MethodDelete, c.fullURL(urlSuffix),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

// ListVectors Lists the currently available assistants.
func (c *Client) ListVectors(
	ctx context.Context,
	limit *int,
	order *string,
	after *string,
	before *string,
) (response VectorList, err error) {
	urlValues := url.Values{}
	if limit != nil {
		urlValues.Add("limit", fmt.Sprintf("%d", *limit))
	}
	if order != nil {
		urlValues.Add("order", *order)
	}
	if after != nil {
		urlValues.Add("after", *after)
	}
	if before != nil {
		urlValues.Add("before", *before)
	}

	encodedValues := ""
	if len(urlValues) > 0 {
		encodedValues = "?" + urlValues.Encode()
	}

	urlSuffix := fmt.Sprintf("%s%s", vectorSuffix, encodedValues)
	req, err := c.newRequest(ctx, http.MethodGet, c.fullURL(urlSuffix),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

// CreateVectorFile creates a new assistant file.
func (c *Client) CreateVectorFile(
	ctx context.Context,
	vectorID string,
	request VectorFileFileRequest,
) (response AssistantFile, err error) {
	urlSuffix := fmt.Sprintf("%s/%s%s", vectorSuffix, vectorID, vectorFilesSuffix)
	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL(urlSuffix),
		withBody(request),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

// RetrieveAssistantFile retrieves an assistant file.
func (c *Client) RetrieveVectorFile(
	ctx context.Context,
	vectorId string,
	fileID string,
) (response AssistantFile, err error) {
	urlSuffix := fmt.Sprintf("%s/%s%s/%s", vectorSuffix, vectorId, vectorFilesSuffix, fileID)
	req, err := c.newRequest(ctx, http.MethodGet, c.fullURL(urlSuffix),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

// DeleteAssistantFile deletes an existing file.
func (c *Client) DeleteVectorFile(
	ctx context.Context,
	vectorID string,
	fileID string,
) (err error) {
	urlSuffix := fmt.Sprintf("%s/%s%s/%s", vectorSuffix, vectorID, vectorFilesSuffix, fileID)
	req, err := c.newRequest(ctx, http.MethodDelete, c.fullURL(urlSuffix),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, nil)
	return
}

// ListAssistantFiles Lists the currently available files for an assistant.
func (c *Client) ListVectrFiles(
	ctx context.Context,
	vectorID string,
	limit *int,
	order *string,
	after *string,
	before *string,
) (response VectorFilesList, err error) {
	urlValues := url.Values{}
	if limit != nil {
		urlValues.Add("limit", fmt.Sprintf("%d", *limit))
	}
	if order != nil {
		urlValues.Add("order", *order)
	}
	if after != nil {
		urlValues.Add("after", *after)
	}
	if before != nil {
		urlValues.Add("before", *before)
	}

	encodedValues := ""
	if len(urlValues) > 0 {
		encodedValues = "?" + urlValues.Encode()
	}

	urlSuffix := fmt.Sprintf("%s/%s%s%s", vectorSuffix, vectorID, vectorFilesSuffix, encodedValues)
	req, err := c.newRequest(ctx, http.MethodGet, c.fullURL(urlSuffix),
		withBetaAssistantVersion(c.config.AssistantVersion))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}
