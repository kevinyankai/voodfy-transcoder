package task

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Voodfy/voodfy-transcoder/internal/ffmpeg"
	"github.com/Voodfy/voodfy-transcoder/internal/models"
	"github.com/Voodfy/voodfy-transcoder/internal/settings"
	"github.com/Voodfy/voodfy-transcoder/internal/utils"
	ipfsManager "github.com/Voodfy/voodfy-transcoder/pkg/ipfs"
	"github.com/Voodfy/voodfy-transcoder/pkg/livepeerclient"
	"github.com/Voodfy/voodfy-transcoder/pkg/logging"
	"github.com/Voodfy/voodfy-transcoder/pkg/voodfyapi"
)

var cl = ffmpeg.NewClient()

// RemoveAudioFromMp4Task ...
func RemoveAudioFromMp4Task(args ...string) error {
	ffmpeg.Run(&cl, "RemoveAudioFromMp4", args...)

	return nil
}

// ThumbsPreviewGeneratorTask ...
func ThumbsPreviewGeneratorTask(args ...string) error {
	ffmpeg.Run(&cl, "ThumbsPreviewGenerator", args...)
	return nil
}

// GenerateImageFromFrameVideoTask ...
func GenerateImageFromFrameVideoTask(args ...string) error {
	ffmpeg.Run(&cl, "GenerateImageFromFrameVideo", args...)
	return nil
}

// ExtractAudioFromMp4Task ...
func ExtractAudioFromMp4Task(args ...string) error {
	ffmpeg.Run(&cl, "ExtractAudioFromMp4", args...)

	return nil
}

// FallbackRenditionTask ...
func FallbackRenditionTask(args ...string) error {
	ffmpeg.Run(&cl, args[2], args...)
	return nil
}

// RenditionTask will send and receive the chunck transcoded by livepeer
func RenditionTask(args ...string) error {
	client := livepeerclient.NewClient("", args[5], args[6])

	if args[7] == "remote" {
		client.PullToRemote(args[0], args[1], args[2], args[3])
	} else {
		client.PullToLocal(args[0], args[1], args[2], args[3])
	}

	os.Remove(args[4])

	return nil
}

// SendDirToIPFSTask send final directory to ipfs
func SendDirToIPFSTask(args ...string) (string, error) {
	mg, err := ipfsManager.NewManager(settings.IPFSSetting.Gateway)
	logging.Info("Gateway ~>", mg.NodeAddress())

	if err != nil {
		utils.SendError("ipfsManager.NewManager", err)
	}

	// send the directory to ipfs
	cid, err := mg.AddDir(args[0])

	if err != nil {
		utils.SendError("mg.AddDir", err)
	}

	cVoodfyAPI := voodfyapi.NewClient()
	cVoodfyAPI.Endpoint = fmt.Sprintf("/v1/videos?resource=%s", args[1])
	videoID, err := cVoodfyAPI.GetVideoByResourceID(args[1], args[2])

	if err != nil {
		utils.SendError("cVoodfyAPI.UpdateCIDVideoByResourceID", err)
	}

	cVoodfyAPI.Endpoint = fmt.Sprintf("videos?resource=%s", args[1])
	err = cVoodfyAPI.UpdateCIDVideoByResourceID(videoID, cid, args[2])

	if err != nil {
		utils.SendError("cVoodfyAPI.UpdateCIDVideoByResourceID", err)
	}

	cVoodfyAPI.Endpoint = fmt.Sprintf("videos?resource=%s", args[1])
	err = cVoodfyAPI.UpdatePosterVideo(videoID, cid, args[2])

	if err != nil {
		utils.SendError("cVoodfyAPI.UpdateCIDVideoByResourceID", err)
	}

	directory := models.Directory{
		CID: cid,
		ID:  args[1],
	}

	cids, err := mg.List(cid)
	for _, c := range cids {
		resource := models.Resource{
			ID:   utils.EncodeMD5(c.Hash),
			Name: c.Name,
			CID:  c.Hash,
		}
		directory.Resources = append(directory.Resources, resource)
	}
	directory.Save()
	return cid, err
}

// LongRunningTask ...
func LongRunningTask() error {
	for i := 0; i < 10; i++ {
		log.Println(fmt.Sprintf("%d", 10-i))
		time.Sleep(1 * time.Second)
	}
	return nil
}
