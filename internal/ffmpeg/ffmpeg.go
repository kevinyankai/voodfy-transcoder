package ffmpeg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/Voodfy/voodfy-transcoder/internal/models"
	"github.com/Voodfy/voodfy-transcoder/internal/utils"
)

var (
	// ExecFunc is command func.
	ExecFunc = ExecCmd
)

// Commands interface from ffmpeg
type Commands interface {
	RemoveAudioFromMP4(string, string) bool
	GenerateImageFromFrameVideo(string, string, string) bool
	GenerateWebpFromFrameVideo(string, string, string) bool
	Transcode90p(string, string) bool
	Transcode144p(string, string) bool
	Transcode240p(string, string) bool
	Transcode360p(string, string) bool
	Transcode480p(string, string) bool
	Transcode720p(string, string) bool
	Transcode1080p(string, string) bool
	ConvertToMp4(string, string) bool
	ThumbsPreviewGenerator(string, string, string) bool
	VTTGenerator(string, string, string) bool
	ExtractAudioFromMp4(string, string) bool
	CheckIntegrityFromMp4s(string, string) bool
}

// Client instance of ffmpeg
type Client struct{}

// NewClient return a instance of ffmpeg
func NewClient() (c Client) {
	return Client{}
}

// Run execute the ffmpeg job
func Run(cmd Commands, fnc string, args ...string) bool {
	switch fnc {
	case "FFprobe":
		r, _ := Execute(args[0])
		r.ID = args[1]
		r.Save()
	case "RemoveAudioFromMp4":
		return cmd.RemoveAudioFromMP4(args[0], args[1])
	case "ThumbsPreviewGenerator":
		r, _ := Execute(args[0])
		return cmd.ThumbsPreviewGenerator(args[0], args[1], r.Format.Duration)
	case "GenerateImageFromFrameVideo":
		r, _ := Execute(args[0])
		return cmd.GenerateImageFromFrameVideo(args[0], args[1], r.Format.Duration)
	case "ExtractAudioFromMp4":
		return cmd.ExtractAudioFromMp4(args[0], args[1])
	case "90p":
		cmd.Transcode90p(args[0], args[1])
	case "144p":
		cmd.Transcode144p(args[0], args[1])
	case "240p":
		cmd.Transcode240p(args[0], args[1])
	case "360p":
		cmd.Transcode360p(args[0], args[1])
	case "480p":
		cmd.Transcode480p(args[0], args[1])
	case "720p":
		cmd.Transcode720p(args[0], args[1])
	case "1080p":
		cmd.Transcode1080p(args[0], args[1])
	case "convertToMp4":
		cmd.ConvertToMp4(args[0], args[1])
	}
	return false
}

// ExecCmd exec ffprobe command and return result of json.
func ExecCmd(fileName string) ([]byte, error) {
	return exec.Command("ffprobe",
		"-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", fileName).Output()
}

// Execute exec command and bind result to struct.
func Execute(fileName string) (r models.Specification, err error) {
	out, err := ExecFunc(fileName)

	if err != nil {
		return r, err
	}

	if err := json.Unmarshal(out, &r); err != nil {
		return r, err
	}

	return r, nil
}

// RemoveAudioFromMP4 generate a mp4 without audion
func (c *Client) RemoveAudioFromMP4(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer

	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-c", "copy", "-an", dstFile)
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-RemoveAudioFromMP4-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-RemoveAudioFromMP4-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}
	return true
}

// GenerateImageFromFrameVideo generate a jpg from mp4
func (c *Client) GenerateImageFromFrameVideo(filename, dstFile, duration string) bool {
	var stdBuffer bytes.Buffer
	var position string
	position = "00:00:01"

	d, err := strconv.ParseFloat(duration, 64)

	if d >= 5.0 {
		position = "00:00:05"
	}

	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-ss", position, "-i", filename, "-vframes", "1", "-q:v", "1", fmt.Sprintf("%sposter.jpg", dstFile))
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err = cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-GenerateImageFromFrameVideo-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-GenerateImageFromFrameVideo-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}
	return true
}

// GenerateWebpFromFrameVideo generate a jpg from mp4
func (c *Client) GenerateWebpFromFrameVideo(filename, dstFile, duration string) bool {
	var stdBuffer bytes.Buffer
	var position string
	position = "00:00:01"
	d, err := strconv.ParseFloat(duration, 64)

	if d >= 5.0 {
		position = "00:00:05"
	}

	cmd := exec.Command("ffmpeg", "-hide_banner", "-i", filename, "-lossless", "0", "-ss", "00:00:00", "-t", position, "-s", "384x182", fmt.Sprintf("%sposter.webp", dstFile))
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err = cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-ConvertToMp4-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-ConvertToMp4-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}
	log.Println("finish 240p ~> ", filename)

	return true
}

// ConvertToMp4 convert mkv to mp4
func (c *Client) ConvertToMp4(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer

	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-c", "copy", dstFile)

	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw
	err := cmd.Start()

	if err != nil {
		utils.SendError(fmt.Sprintf("%s-ConvertToMp4-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-ConvertToMp4-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	return true
}

// Transcode90p low definition
func (c *Client) Transcode90p(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer

	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:90'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "100k", "-an", dstFile)

	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw
	err := cmd.Start()

	if err != nil {
		utils.SendError(fmt.Sprintf("%s-Transcode90p-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-Transcode90p-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	return true
}

// Transcode144p low definition
func (c *Client) Transcode144p(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer

	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:90'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "100k", "-an", dstFile)

	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw
	err := cmd.Start()

	if err != nil {
		utils.SendError(fmt.Sprintf("%s-Transcode144p-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-Transcode144p-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	return true
}

// Transcode240p 240p
func (c *Client) Transcode240p(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:240'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "120k", "-an", dstFile)
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-Transcode240p-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-Transcode240p-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}
	return true
}

// Transcode360p 360p
func (c *Client) Transcode360p(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:360'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "284k", "-maxrate", "284k", "-bufsize", "568k", "-an", dstFile)
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw
	err := cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-Transcode360p-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-Transcode360p-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}
	return true
}

// Transcode480p 480p
func (c *Client) Transcode480p(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer

	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:480'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "341k", "-maxrate", "341k", "-bufsize", "682k", "-an", dstFile)
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-Transcode480p-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-Transcode480p-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	return true
}

// Transcode720p 720p
func (c *Client) Transcode720p(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:720'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "765k", "-maxrate", "765k", "-bufsize", "1530k", "-an", dstFile)
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw
	err := cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-Transcode720p-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-Transcode720p-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}
	return true
}

// Transcode1080p 1080p
func (c *Client) Transcode1080p(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:1080'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "1579k", "-maxrate", "1579k", "-bufsize", "3158k", "-an", dstFile)

	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Start()

	if err != nil {
		utils.SendError(fmt.Sprintf("%s-Transcode1080p-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-Transcode1080p-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}
	return true
}

// ThumbsPreviewGenerator ...
func (c *Client) ThumbsPreviewGenerator(filename, dstFile, duration string) bool {
	var err error
	d, err := strconv.ParseFloat(duration, 64)
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-RemoveAudioFromMP4-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}
	columnsTotal := int(d) / 5 / 2
	cmd := exec.Command("thumbsgenerator", filename, "5", "126", "73", fmt.Sprintf("%d", columnsTotal), fmt.Sprintf("%s/thumbspreview.png", dstFile))

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw
	err = cmd.Start()

	if err != nil {
		utils.SendError(fmt.Sprintf("%s-ThumbsPreviewGenerator-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-ThumbsPreviewGenerator-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}
	return true
}

// VTTGenerator ...
func (c *Client) VTTGenerator(filename, dstFile, language string) bool {
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-f", "webvtt", fmt.Sprintf("%s_%s.vtt", dstFile, language))

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw
	err := cmd.Start()

	if err != nil {
		utils.SendError(fmt.Sprintf("%s-VTTGenerator.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-VTTGenerator-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}
	return true
}

// ExtractAudioFromMp4 generate a m4a extracting the audio from mp4
func (c *Client) ExtractAudioFromMp4(filename, dstFile string) bool {
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-vn", "-acodec", "copy", dstFile)
	err := cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-ExtractAudioFromMp4-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("%s-ExtractAudioFromMp4-cmd.Start() failed with '%s'\n", filename, err), err)
		return false
	}

	return true
}

// CheckIntegrityFromMp4s return a boolean about the duration of the video
func (c *Client) CheckIntegrityFromMp4s(source, output string) bool {
	s, _ := Execute(source)
	o, _ := Execute(output)

	if o.Format.Duration == s.Format.Duration {
		return true
	}

	return false
}
