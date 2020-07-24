package task

import (
	"fmt"
	"log"
	"time"

	"github.com/Voodfy/voodfy-transcoder/internal/ffmpeg"
	"github.com/Voodfy/voodfy-transcoder/internal/logging"
	"github.com/Voodfy/voodfy-transcoder/internal/models"
	"github.com/Voodfy/voodfy-transcoder/internal/settings"
	"github.com/Voodfy/voodfy-transcoder/internal/utils"
	ipfsManager "github.com/Voodfy/voodfy-transcoder/pkg/ipfs"
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
