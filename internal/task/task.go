package task

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Voodfy/voodfy-transcoder/internal/ffmpeg"
	"github.com/Voodfy/voodfy-transcoder/internal/models"
	"github.com/Voodfy/voodfy-transcoder/internal/settings"
	"github.com/Voodfy/voodfy-transcoder/internal/utils"
	ipfsManager "github.com/Voodfy/voodfy-transcoder/pkg/ipfs"
	"github.com/Voodfy/voodfy-transcoder/pkg/livepeerclient"
	"github.com/Voodfy/voodfy-transcoder/pkg/logging"
	"github.com/Voodfy/voodfy-transcoder/pkg/powergate"
	cid "github.com/ipfs/go-cid"
	clusterApi "github.com/ipfs/ipfs-cluster/api"
	client "github.com/ipfs/ipfs-cluster/api/rest/client"
	multiaddr "github.com/multiformats/go-multiaddr"
)

var cl = ffmpeg.NewClient()

// ConvertToMp4Task ...
func ConvertToMp4Task(args ...string) error {
	ffmpeg.Run(&cl, "convertToMp4", args...)

	return nil
}

// FFprobeTask ...
func FFprobeTask(args ...string) error {
	ffmpeg.Run(&cl, "FFprobe", args...)

	return nil
}

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
	client := livepeerclient.NewClient()

	if settings.LivepeerSetting.Remote {
		client.PullToRemote(args[0], args[1], args[2], args[3])
	} else {
		client.PullToLocal(args[0], args[1], args[2], args[3])
	}

	return nil
}

// SendDirToIPFSTask send final directory to ipfs
func SendDirToIPFSTask(args ...string) (string, error) {
	var idx int
	idx = 2

	mg, err := ipfsManager.NewManager(settings.IPFSSetting.Gateway)
	logging.Info("Gateway ~>", mg.NodeAddress())

	if err != nil {
		utils.SendError("ipfsManager.NewManager", err)
	}

	ticker := time.NewTicker(settings.AppSetting.DelayWaitingIPFS * time.Second)
	select {
	case _ = <-ticker.C:
		if _, err := os.Stat(args[0]); !os.IsNotExist(err) {
			ticker.Stop()
			break
		}
	}

	entries, err := ioutil.ReadDir(args[0])
	if err != nil {
		utils.SendError("ioutil.ReadDir", err)
	}
	for _, entry := range entries {
		extension := filepath.Ext(entry.Name())
		if extension == ".mp4" {
			sourcePath := filepath.Join(args[0], entry.Name())
			newPath := filepath.Join(args[0], fmt.Sprintf("%s_v%d.mp4", args[1], idx))
			err := os.Rename(sourcePath, newPath)
			if err != nil {
				utils.SendError("os.Rename", err)
			}
			idx++
		}
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
		mg.Pin(c.Hash)
	}
	directory.Save()
	return cid, err
}

// PinDirToIPFSClusterTask send final directory to ipfs cluster
func PinDirToIPFSClusterTask(args ...string) error {
	var wait bool
	cfg := &client.Config{}
	addr, err := multiaddr.NewMultiaddr(settings.IPFSSetting.ClusterGateway)
	if err != nil {
		log.Fatal(err)
	}

	cfg.APIAddr = addr

	c, err := client.NewDefaultClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ci, err := cid.Decode(args[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("pinning: %s\t%s\n", ci, args[1])
	_, err = c.Pin(context.Background(), ci, clusterApi.PinOptions{Name: args[1]})
	if err != nil {
		log.Println(err)
	}
	if !wait {
	}
	_, err = client.WaitFor(context.Background(), c, client.StatusFilterParams{
		Cid:       ci,
		Target:    clusterApi.TrackerStatusPinned,
		CheckFreq: 5000 * time.Millisecond,
		Local:     false,
	})
	if err != nil {
		log.Println(err)
	}
	return err
}

// SendDirToFilecoinTask send final directory to filecoin
func SendDirToFilecoinTask(args ...string) ([]string, error) {
	var jids []string
	mg, err := ipfsManager.NewManager(settings.IPFSSetting.Gateway)
	logging.Info("Gateway ~>", mg.NodeAddress())

	if err != nil {
		utils.SendError("ipfsManager.NewManager", err)
	}

	cids, err := mg.List(args[0])
	for _, c := range cids {
		jids = append(jids, powergate.FFSPush(c.Hash, args[1], settings.AppSetting.HostedPowergateAddr))
	}
	return jids, err
}

// LongRunningTask ...
func LongRunningTask() error {
	for i := 0; i < 10; i++ {
		log.Println(fmt.Sprintf("%d", 10-i))
		time.Sleep(1 * time.Second)
	}
	return nil
}
