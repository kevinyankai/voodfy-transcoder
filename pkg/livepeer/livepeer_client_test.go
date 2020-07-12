package livepeer_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/streamgoinc/transcoder/pkg/livepeer"
	"github.com/tj/assert"
)

type LivepeerResourcesMock struct {
	statusCode            int
	livepeerResourcesFile string
}

func buildEmptyBody() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte("")))
}

func buildBody(livepeerMock *LivepeerResourcesMock, w http.ResponseWriter) {
	dataRelativePath := filepath.Join("testdata/", livepeerMock.livepeerResourcesFile)
	data, _ := ioutil.ReadFile(dataRelativePath)
	w.Write([]byte(data))
}

func (livepeerMock *LivepeerResourcesMock) httpHandler(w http.ResponseWriter, r *http.Request) {
	if livepeerMock.statusCode == 302 {
		w.Header().Add("Location", r.URL.String())
	}

	w.WriteHeader(livepeerMock.statusCode)
	buildBody(livepeerMock, w)
}

func TestHTTPLivepeerRequest(t *testing.T) {
	livepeerMock := &LivepeerResourcesMock{}

	livepeerServer := httptest.NewServer(http.HandlerFunc(livepeerMock.httpHandler))

	target := livepeer.NewClient(livepeerServer.URL)

	defer livepeerServer.Close()

	t.Run("TestlivepeerBroadcastersRequests", testLivepeerBroadcastersRequests(target, livepeerMock))
}

func testLivepeerBroadcastersRequests(target *livepeer.Client, livepeerMock *LivepeerResourcesMock) func(*testing.T) {
	testsScenarios := []struct {
		scenarioName                        string
		statusCode                          int
		expectedLivepeerBroadcasterResponse livepeer.Broadcasters
		livepeerPayoutResourceFile          string
	}{
		{
			scenarioName:               "ShouldReturnBroadcastersOK",
			statusCode:                 http.StatusOK,
			livepeerPayoutResourceFile: "broadcasters.json",
			expectedLivepeerBroadcasterResponse: livepeer.Broadcasters{
				livepeer.Broadcaster{
					Address: "https://chi-broadcaster-squirtle.livepeer-ac.live",
				},
				livepeer.Broadcaster{
					Address: "https://chi-broadcaster-charmander.livepeer-ac.live",
				},
				livepeer.Broadcaster{
					Address: "https://chi-broadcaster-bulbasaur.livepeer-ac.live",
				},
			},
		},
	}
	return func(t *testing.T) {
		for _, scn := range testsScenarios {
			t.Run(scn.scenarioName, func(t *testing.T) {
				livepeerMock.statusCode = scn.statusCode
				livepeerMock.livepeerResourcesFile = scn.livepeerPayoutResourceFile

				payload := []interface{}{"token"}

				broadcasterResponse := target.GetBroadcasters()

				assert.Equal(t, scn.expectedLivepeerBroadcasterResponse, broadcasterResponse)

			})
		}
	}
}
