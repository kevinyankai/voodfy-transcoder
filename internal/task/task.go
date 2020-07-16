package task

import (
	"fmt"
	"os"
	"time"

	"github.com/Voodfy/voodfy-transcoder/internal/ffmpeg"
	"github.com/Voodfy/voodfy-transcoder/internal/logging"
	"github.com/Voodfy/voodfy-transcoder/internal/models"
	"github.com/Voodfy/voodfy-transcoder/internal/settings"
	"github.com/Voodfy/voodfy-transcoder/internal/utils"
	ipfsManager "github.com/Voodfy/voodfy-transcoder/pkg/ipfs"
	"github.com/Voodfy/voodfy-transcoder/pkg/livepeer"
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

// StandardFallbackTask ...
func StandardFallbackTask(args ...string) error {
	ffmpeg.Run(&cl, "StandardFallback", args...)

	return nil
}

// MidDefinitionFallbackTask ...
func MidDefinitionFallbackTask(args ...string) error {
	ffmpeg.Run(&cl, "MidDefinitionFallback", args...)

	return nil
}

// HDFallBackTask ...
func HDFallBackTask(args ...string) error {
	ffmpeg.Run(&cl, "HDFallBackTask", args...)

	return nil
}

// FullHDFallbackTask ...
func FullHDFallbackTask(args ...string) error {
	ffmpeg.Run(&cl, "FullHDFallbackTask", args...)

	return nil
}

// LowDefinitionTask ...
func LowDefinitionTask(args ...string) error {
	ffmpeg.Run(&cl, "LowDefinition", args...)

	return nil
}

// RenditionTask will send and receive the chunck transcoded by livepeer
func RenditionTask(args ...string) error {
	var bucket = settings.ServerSetting.BucketMount

	client := livepeer.NewClient("")
	srcFile := fmt.Sprintf("%s%s/%s/%s_without_audio.mp4", bucket, args[0], args[1], args[2])
	dstFiles := fmt.Sprintf("%s%s/%s/%s_ipfs", bucket, args[0], args[1], args[2])
	profile := fmt.Sprintf("%s%s", bucket, args[3])

	if settings.AppSetting.LivepeerMode == "remote" {
		client.PullToRemote(srcFile, dstFiles, profile, args[2])
	} else {
		client.PullToLocal(srcFile, dstFiles, profile, args[2])
	}

	os.Remove(fmt.Sprintf("%s/%s_source.mp4", dstFiles, args[2]))

	return nil
}

// SendDirToIPFSTask send final directory to ipfs
func SendDirToIPFSTask(args ...string) (string, error) {
	var bucket = settings.ServerSetting.BucketMount
	dstFiles := fmt.Sprintf("%s%s/%s/%s_ipfs", bucket, args[0], args[1], args[2])
	mg, err := ipfsManager.NewManager(settings.IPFSSetting.Gateway)
	logging.Info("Gateway ~>", mg.NodeAddress())

	if err != nil {
		utils.SendError("ipfsManager.NewManager", err)
	}

	// send the directory to ipfs
	cid, err := mg.AddDir(dstFiles)

	if err != nil {
		utils.SendError("mg.AddDir", err)
	}

	directory := models.Directory{
		CID: cid,
		ID:  args[2],
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
		logging.Info(fmt.Sprintf("%d", 10-i))
		time.Sleep(1 * time.Second)
	}
	return nil
}
