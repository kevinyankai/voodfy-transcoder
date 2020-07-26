package livepeerclient

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/Voodfy/voodfy-transcoder/pkg/logging"
	"gopkg.in/resty.v1"
)

// Client struct used to store livepeer client
type Client struct {
	OrchAddr       string
	OrchWebhookURL string
	Resty          *resty.Client
}

// NewClient func to return a instance from livepeer client
func NewClient(baseURL, orchWebhookURL, orchAddr string) *Client {
	client := resty.New()

	livepeer := &Client{
		Resty:          client,
		OrchWebhookURL: orchWebhookURL,
		OrchAddr:       orchAddr,
	}

	return livepeer
}

// PullToRemote the src file to be transcoded on livepeer
func (c *Client) PullToRemote(src, dst, profile, id string) bool {
	cmd := exec.Command("livepeer", "-pull", src, "-recordingDir", dst, "-transcodingOptions", profile, "-orchWebhookUrl", c.OrchWebhookURL, "-streamName", id, "-v", "99")
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw
	err := cmd.Start()
	if err != nil {
		logging.Info(fmt.Sprintf("cmd.Start() failed with '%s'\n", err))
		return false
	}

	err = cmd.Wait()
	if err != nil {
		logging.Info(fmt.Sprintf("cmd.Run() failed with %s\n", err))
		return false
	}
	return true
}

// PullToLocal the src file to be transcoded on livepeer
func (c *Client) PullToLocal(src, dst, profile, id string) bool {
	cmd := exec.Command("livepeer", "-pull", src, "-recordingDir", dst, "-transcodingOptions", profile, "-orchAddr", c.OrchAddr, "-streamName", id, "-v", "99")
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw
	err := cmd.Start()
	if err != nil {
		logging.Info(fmt.Sprintf("cmd.Start() failed with '%s'\n", err))
		return false
	}

	err = cmd.Wait()
	if err != nil {
		logging.Info(fmt.Sprintf("cmd.Run() failed with %s\n", err))
		return false
	}
	return true
}
