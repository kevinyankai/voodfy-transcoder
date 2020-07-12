package livepeer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/Voodfy/voodfy-transcoder/internal/logging"
	"github.com/Voodfy/voodfy-transcoder/internal/settings"
	"github.com/go-resty/resty/v2"
)

// Client struct used to store livepeer client
type Client struct {
	Resty   *resty.Client
	BaseURL string
	URI     string
	Token   string
	Payload map[string]interface{}
}

func getBroadcaster() string {
	return settings.AppSetting.LivepeerBroadcaster
}

func getToken() string {
	return settings.AppSetting.LivepeerToken
}

func getBaseURL() string {
	return "https://livepeer.live/api/"
}

// NewClient func to return a instance from livepeer client
func NewClient(baseURL string) *Client {
	client := resty.New()
	token := getToken()

	livepeer := &Client{
		Token: token,
		Resty: client,
	}

	livepeer.BaseURL = getBaseURL()
	if os.Getenv("GIN_MODE") == "test" {
		livepeer.BaseURL = baseURL
	}

	return livepeer
}

// PullToRemote the src file to be transcoded on livepeer
func (c *Client) PullToRemote(src, dst, profile, id string) bool {
	cmd := exec.Command("livepeer", "-pull", src, "-recordingDir", dst, "-transcodingOptions", profile, "-orchWebhookUrl", settings.AppSetting.LivepeerBroadcaster, "-streamName", id, "-v", "99")
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
	broadcaster := getBroadcaster()
	cmd := exec.Command("livepeer", "-pull", src, "-recordingDir", dst, "-transcodingOptions", profile, "-orchAddr", broadcaster, "-streamName", id, "-v", "99")
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
