package voodfyapi

import (
	"encoding/json"
	"fmt"

	"github.com/Voodfy/voodfy-transcoder/internal/settings"
	"gopkg.in/resty.v1"
)

// HTTPClient interface to encapsulate the logic behind client
type HTTPClient struct {
	Endpoint string                 `json:"endpoint"`
	Payload  map[string]interface{} `json:"payload"`
	BaseURL  string
}

// Client used to send requests to api
type Client interface {
	Patch(string) error
}

// NewClient return an instance of the client
func NewClient() HTTPClient {
	return HTTPClient{
		BaseURL: "https://publish.voodfy.com",
	}
}

// URL return the url
func (c *HTTPClient) URL() string {
	return fmt.Sprintf("%s%s", c.BaseURL, c.Endpoint)
}

// Post do request to voodfyAPI
func (c *HTTPClient) Post() error {
	_, err := resty.R().SetBody(c.Payload).Post(c.URL())
	return err
}

// GetVideoByResourceID get the video by resource id
func (c *HTTPClient) GetVideoByResourceID(rsID, token string) (videoID string, err error) {
	var response Response
	header := map[string]string{
		"Authorization": fmt.Sprintf("Token %s", token),
	}
	rsp, err := resty.R().SetHeaders(header).Get(c.URL())
	json.Unmarshal(rsp.Body(), &response)
	videoID = response.Result.Videos[0].ID
	return
}

// UpdateCIDVideoByResourceID update the video with cid from ipfs
func (c *HTTPClient) UpdateCIDVideoByResourceID(id, cid, token string) (err error) {
	header := map[string]string{
		"Authorization": fmt.Sprintf("Token %s", token),
	}
	c.Payload = map[string]interface{}{"cid": cid}
	c.Endpoint = fmt.Sprintf("/v1/videos/%s/cid", id)
	_, err = resty.R().SetBody(c.Payload).SetHeaders(header).Patch(c.URL())
	return
}

// UpdatePosterVideo update the video poster
func (c *HTTPClient) UpdatePosterVideo(id, cid, token string) (err error) {
	header := map[string]string{
		"Authorization": fmt.Sprintf("Token %s", token),
	}
	poster := fmt.Sprintf("%s/ipfs/%s/poster.jpg", settings.IPFSSetting.Origin, cid)
	c.Payload = map[string]interface{}{"poster": poster}
	c.Endpoint = fmt.Sprintf("/v1/videos/%s", id)
	_, err = resty.R().SetBody(c.Payload).SetHeaders(header).Patch(c.URL())
	return
}

// Token do request to get token
func (c *HTTPClient) Token() (string, error) {
	var response Response
	c.Endpoint = "/v1/auth-by-device"
	rsp, err := resty.R().SetBody(c.Payload).Post(c.URL())
	json.Unmarshal(rsp.Body(), &response)
	return response.Result.User.Token, err
}

// Signup do request to signup
func (c *HTTPClient) Signup() error {
	c.Endpoint = "/v1/users-by-device"
	_, err := resty.R().SetBody(c.Payload).Post(c.URL())
	return err
}

// Retrieve do request to retrieve user by secret
func (c *HTTPClient) Retrieve(secret string) (string, error) {
	var response Response
	c.Endpoint = "/v1/user-by-device"
	rsp, err := resty.R().SetQueryParams(map[string]string{
		"encodedDevice": secret,
	}).Get(c.URL())
	json.Unmarshal(rsp.Body(), &response)
	return response.Result.User.Device, err
}

// Powergate do request to retrieve powergate instance
func (c *HTTPClient) Powergate(secret string, premium bool) (Powergate, error) {
	var response Response
	c.Endpoint = "/v1/powergate"

	rsp, err := resty.R().Get(c.URL())

	json.Unmarshal(rsp.Body(), &response)
	return response.Result.Powergate, err
}

// Embed do request to create a video at Voodfy
func (c *HTTPClient) Embed(token, title, description, cid string) (Video, error) {
	var response Response
	c.Endpoint = "/v1/videos"
	c.Payload = map[string]interface{}{
		"title":       title,
		"description": description,
		"cid":         cid,
		"ipfs":        fmt.Sprintf("https://ipfs.voodfy.com/ipfs/%s", cid),
		"poster":      fmt.Sprintf("https://ipfs.voodfy.com/ipfs/%s/poster.jpg", cid),
	}

	header := map[string]string{
		"Authorization": fmt.Sprintf("Token %s", token),
	}
	rsp, err := resty.R().SetHeaders(header).SetBody(c.Payload).Post(c.URL())

	json.Unmarshal(rsp.Body(), &response)
	return response.Result.Video, err
}
