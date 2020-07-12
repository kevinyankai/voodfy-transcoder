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
	"strings"

	"github.com/Voodfy/voodfy-transcoder/internal/settings"
	"github.com/Voodfy/voodfy-transcoder/internal/utils"
)

var (
	// ExecFunc is command func.
	ExecFunc = ExecCmd
)

// Commands interface from ffmpeg
type Commands interface {
	RemoveAudioFromMP4(string, string, string) bool
	GenerateImageFromFrameVideo(string, string, string, string) bool
	GenerateWebpFromFrameVideo(string, string, string, string) bool
	LowDefinition(string, string, string) bool
	Transcode240p(string, string, string) bool
	Transcode360p(string, string, string) bool
	Transcode480p(string, string, string) bool
	Transcode720p(string, string, string) bool
	Transcode1080p(string, string, string) bool
	ThumbsPreviewGenerator(string, string, string, string) bool
	VTTGenerator(string, string, string, string) bool
	ExtractAudioFromMp4(string, string, string) bool
	SplitMp4IntoChunks(string, string, string) bool
	CheckIntegrityFromMp4s(string, string) bool
}

// Client instance of ffmpeg
type Client struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	User     string `json:"user"`
	Tracker  string `json:"tracker"`
	Position string `json:"position"`
	Duration string `json:"duration"`
	Language string `json:"language"`
}

// Specification struct to save info provided by ffprobe
type Specification struct {
	Format struct {
		BitRate        string `json:"bit_rate"`
		Duration       string `json:"duration"`
		Filename       string `json:"filename"`
		FormatLongName string `json:"format_long_name"`
		FormatName     string `json:"format_name"`
		NbPrograms     int    `json:"nb_programs"`
		NbStreams      int    `json:"nb_streams"`
		ProbeScore     int    `json:"probe_score"`
		Size           string `json:"size"`
		StartTime      string `json:"start_time"`
		Tags           struct {
			CompatibleBrands string `json:"compatible_brands"`
			Encoder          string `json:"encoder"`
			MajorBrand       string `json:"major_brand"`
			MinorVersion     string `json:"minor_version"`
		} `json:"tags"`
	} `json:"format"`
	Streams []struct {
		AvgFrameRate       string `json:"avg_frame_rate"`
		BitRate            string `json:"bit_rate"`
		BitsPerRawSample   string `json:"bits_per_raw_sample,omitempty"`
		ChromaLocation     string `json:"chroma_location,omitempty"`
		CodecLongName      string `json:"codec_long_name"`
		CodecName          string `json:"codec_name"`
		CodecTag           string `json:"codec_tag"`
		CodecTagString     string `json:"codec_tag_string"`
		CodecTimeBase      string `json:"codec_time_base"`
		CodecType          string `json:"codec_type"`
		CodedHeight        int    `json:"coded_height,omitempty"`
		CodedWidth         int    `json:"coded_width,omitempty"`
		DisplayAspectRatio string `json:"display_aspect_ratio,omitempty"`
		Disposition        struct {
			AttachedPic     int `json:"attached_pic"`
			CleanEffects    int `json:"clean_effects"`
			Comment         int `json:"comment"`
			Default         int `json:"default"`
			Dub             int `json:"dub"`
			Forced          int `json:"forced"`
			HearingImpaired int `json:"hearing_impaired"`
			Karaoke         int `json:"karaoke"`
			Lyrics          int `json:"lyrics"`
			Original        int `json:"original"`
			TimedThumbnails int `json:"timed_thumbnails"`
			VisualImpaired  int `json:"visual_impaired"`
		} `json:"disposition"`
		Duration          string `json:"duration"`
		DurationTs        int    `json:"duration_ts"`
		HasBFrames        int    `json:"has_b_frames,omitempty"`
		Height            int    `json:"height,omitempty"`
		Index             int    `json:"index"`
		IsAvc             string `json:"is_avc,omitempty"`
		Level             int    `json:"level,omitempty"`
		NalLengthSize     string `json:"nal_length_size,omitempty"`
		NbFrames          string `json:"nb_frames"`
		PixFmt            string `json:"pix_fmt,omitempty"`
		Profile           string `json:"profile"`
		RFrameRate        string `json:"r_frame_rate"`
		Refs              int    `json:"refs,omitempty"`
		SampleAspectRatio string `json:"sample_aspect_ratio,omitempty"`
		StartPts          int    `json:"start_pts"`
		StartTime         string `json:"start_time"`
		Tags              struct {
			HandlerName string `json:"handler_name"`
			Language    string `json:"language"`
		} `json:"tags"`
		TimeBase      string `json:"time_base"`
		Width         int    `json:"width,omitempty"`
		BitsPerSample int    `json:"bits_per_sample,omitempty"`
		ChannelLayout string `json:"channel_layout,omitempty"`
		Channels      int    `json:"channels,omitempty"`
		MaxBitRate    string `json:"max_bit_rate,omitempty"`
		SampleFmt     string `json:"sample_fmt,omitempty"`
		SampleRate    string `json:"sample_rate,omitempty"`
	} `json:"streams"`
}

// NewClient return a instance of ffmpeg
func NewClient() (c Client) {
	return Client{}
}

func getSrc(args ...string) (srcFile string) {
	var filename string
	var bucket = settings.ServerSetting.BucketMount

	filename = strings.Split(args[1], ".mp4")[0]
	srcFile = fmt.Sprintf("%s%s/%s/%s", bucket, args[2], args[3], filename)
	return
}

func getDst(args ...string) (dstFile string) {
	var bucket = settings.ServerSetting.BucketMount

	dstFile = fmt.Sprintf("%s%s/%s/", bucket, args[2], args[3])
	return
}

func getSourceIPFS(args ...string) (dstFile string) {
	var bucket = settings.ServerSetting.BucketMount

	dstFile = fmt.Sprintf("%s%s/%s/%s_ipfs/", bucket, args[2], args[3], args[0])
	os.MkdirAll(dstFile, 0777)
	return
}

// GetID return the id of the instance
func (c *Client) GetID() string {
	return c.ID
}

// Run execute the ffmpeg job
func Run(cmd Commands, fnc string, args ...string) bool {
	switch fnc {
	case "RemoveAudioFromMp4":
		srcFile := getSrc(args...)
		dstFile := getDst(args...)
		filename := fmt.Sprintf("%s.mp4", srcFile)
		return cmd.RemoveAudioFromMP4(filename, dstFile, args[0])
	case "ThumbsPreviewGenerator":
		srcFile := getSrc(args...)
		dstFile := getSourceIPFS(args...)
		r, _ := Execute(fmt.Sprintf("%s", srcFile))
		filename := fmt.Sprintf("%s.mp4", srcFile)
		return cmd.ThumbsPreviewGenerator(filename, dstFile, args[0], r.Format.Duration)
	case "GenerateImageFromFrameVideo":
		srcFile := getSrc(args...)
		dstFile := getSourceIPFS(args...)
		r, _ := Execute(fmt.Sprintf("%s", srcFile))
		filename := fmt.Sprintf("%s.mp4", srcFile)
		return cmd.GenerateImageFromFrameVideo(filename, dstFile, args[0], r.Format.Duration)
	case "ExtractAudioFromMp4":
		srcFile := getSrc(args...)
		dstFile := getSourceIPFS(args...)
		return cmd.ExtractAudioFromMp4(srcFile, dstFile, args[0])
	case "LowDefinition":
		srcFile := getSrc(args...)
		dstFile := getSourceIPFS(args...)
		filename := fmt.Sprintf("%s.mp4", srcFile)
		return cmd.LowDefinition(filename, dstFile, args[0])
	case "StandardFallback":
		srcFile := getSrc(args...)
		dstFile := getSourceIPFS(args...)
		filename := fmt.Sprintf("%s.mp4", srcFile)
		cmd.Transcode240p(filename, dstFile, args[0])
		return cmd.Transcode360p(filename, dstFile, args[0])
	case "MidDefinitionFallbackTask":
		srcFile := getSrc(args...)
		dstFile := getSourceIPFS(args...)
		filename := fmt.Sprintf("%s.mp4", srcFile)
		return cmd.Transcode480p(filename, dstFile, args[0])
	case "HDFallBackTask":
		srcFile := getSrc(args...)
		dstFile := getSourceIPFS(args...)
		filename := fmt.Sprintf("%s.mp4", srcFile)
		return cmd.Transcode720p(filename, dstFile, args[0])
	case "FullHDFallbackTask":
		srcFile := getSrc(args...)
		dstFile := getSourceIPFS(args...)
		filename := fmt.Sprintf("%s.mp4", srcFile)
		return cmd.Transcode1080p(filename, dstFile, args[0])
	}
	return false
}

// ExecCmd exec ffprobe command and return result of json.
func ExecCmd(fileName string) ([]byte, error) {
	return exec.Command("ffprobe",
		"-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", fmt.Sprintf("%s.mp4", fileName)).Output()
}

// Execute exec command and bind result to struct.
func Execute(fileName string) (r Specification, err error) {
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
func (c *Client) RemoveAudioFromMP4(filename, dstFile, id string) bool {
	var stdBuffer bytes.Buffer

	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-c", "copy", "-an", fmt.Sprintf("%s%s_without_audio.mp4", dstFile, id))
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("RemoveAudioFromMP4-cmd.Start() failed with '%s'\n", err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("RemoveAudioFromMP4-cmd.Run() failed with %s\n", err), err)
		return false
	}
	return true
}

// GenerateImageFromFrameVideo generate a jpg from mp4
func (c *Client) GenerateImageFromFrameVideo(filename, dstFile, id, duration string) bool {
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
		utils.SendError(fmt.Sprintf("GenerateImageFromFrameVideo-cmd.Start() failed with '%s'\n", err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("GenerateImageFromFrameVideo-cmd.Run() failed with %s\n", err), err)
		return false
	}
	return true
}

// GenerateWebpFromFrameVideo generate a jpg from mp4
func (c *Client) GenerateWebpFromFrameVideo(filename, dstFile, id, duration string) bool {
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
		log.Println(fmt.Sprintf("GenerateWebpFromFrameVideo-cmd.Start() failed with '%s'\n", err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		log.Println(fmt.Sprintf("GenerateWebpFromFrameVideo-cmd.Run() failed with %s\n", err), err)
		return false
	}
	log.Println("finish 240p ~> ", filename)

	return true
}

// LowDefinition low definition
func (c *Client) LowDefinition(filename, dstFile, id string) bool {
	var stdBuffer bytes.Buffer

	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:90'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "100k", "-an", fmt.Sprintf("%s%s_v1.mp4", dstFile, id), "-vf", "scale='-2:144'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "300k", "-an", fmt.Sprintf("%s%s_v2.mp4", dstFile, id))

	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw
	err := cmd.Start()

	if err != nil {
		utils.SendError(fmt.Sprintf("cmd.Start() failed with '%s'\n", err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("cmd.Run() failed with %s\n", err), err)
		return false
	}

	return true
}

// Transcode240p 240p
func (c *Client) Transcode240p(filename, dstFile, id string) bool {
	var stdBuffer bytes.Buffer
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:240'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "120k", "-an", fmt.Sprintf("%s%s_v3.mp4", dstFile, id))
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("Transcode240p-cmd.Start() failed with '%s'\n", err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("Transcode240p-cmd.Run() failed with %s\n", err), err)
		return false
	}
	return true
}

// Transcode360p 360p
func (c *Client) Transcode360p(filename, dstFile, id string) bool {
	var stdBuffer bytes.Buffer
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:360'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "284k", "-maxrate", "284k", "-bufsize", "568k", "-an", fmt.Sprintf("%s%s_v4.mp4", dstFile, id))
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw
	err := cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("Transcode360p-cmd.Start() failed with '%s'\n", err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("Transcode360p-cmd.Run() failed with %s\n", err), err)
		return false
	}
	return true
}

// Transcode480p 480p
func (c *Client) Transcode480p(filename, dstFile, id string) bool {
	var stdBuffer bytes.Buffer

	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:480'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "341k", "-maxrate", "341k", "-bufsize", "682k", "-an", fmt.Sprintf("%s%s_v5.mp4", dstFile, id))
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("Transcode480p-cmd.Start() failed with '%s'\n", err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("Transcode480p-cmd.Run() failed with %s\n", err), err)
		return false
	}

	return true
}

// Transcode720p 720p
func (c *Client) Transcode720p(filename, dstFile, id string) bool {
	var stdBuffer bytes.Buffer
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:720'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "765k", "-maxrate", "765k", "-bufsize", "1530k", "-an", fmt.Sprintf("%s%s_v6.mp4", dstFile, id))
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw
	err := cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("Transcode720p-cmd.Start() failed with '%s'\n", err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("Transcode720p-cmd.Run() failed with %s\n", err), err)
		return false
	}
	return true
}

// Transcode1080p 1080p
func (c *Client) Transcode1080p(filename, dstFile, id string) bool {
	var stdBuffer bytes.Buffer
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:1080'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "1579k", "-maxrate", "1579k", "-bufsize", "3158k", "-an", fmt.Sprintf("%s%s_v7.mp4", dstFile, id))

	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Start()

	if err != nil {
		utils.SendError(fmt.Sprintf("Transcode1080p-cmd.Start() failed with '%s'\n", err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("Transcode1080p-cmd.Run() failed with %s\n", err), err)
		return false
	}
	return true
}

// ThumbsPreviewGenerator ...
func (c *Client) ThumbsPreviewGenerator(filename, dstFile, id, duration string) bool {
	var err error
	d, err := strconv.ParseFloat(duration, 64)
	if err != nil {
		utils.SendError(fmt.Sprintf("ThumnailPreviewGenerator-err-%s", filename), err)
	}
	columnsTotal := int(d) / 5 / 2
	cmd := exec.Command("thumbsgenerator", filename, "5", "126", "73", fmt.Sprintf("%d", columnsTotal), fmt.Sprintf("%s/thumbspreview.png", dstFile))

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw
	err = cmd.Start()

	if err != nil {
		utils.SendError(fmt.Sprintf("ThumbnailPreviewGenerator-cmd.Start() failed with '%s'\n", err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("ThumbnailPreviewGenerator-cmd.Run() failed with %s\n", err), err)
		return false
	}
	return true
}

// VTTGenerator ...
func (c *Client) VTTGenerator(filename, dstFile, id, language string) bool {
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-f", "webvtt", fmt.Sprintf("%s%s_%s.vtt", dstFile, id, language))

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw
	err := cmd.Start()

	if err != nil {
		utils.SendError(fmt.Sprintf("VTTGenerator-cmd.Start() failed with '%s'\n", err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("VTTGenerator-cmd.Run() failed with %s\n", err), err)
		return false
	}
	return true
}

// ExtractAudioFromMp4 generate a m4a extracting the audio from mp4
func (c *Client) ExtractAudioFromMp4(filename, dstFile, id string) bool {
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", fmt.Sprintf("%s.mp4", filename), "-vn", "-acodec", "copy", fmt.Sprintf("%s/%s_a1.m4a", dstFile, id))
	err := cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("ExtractAudioFromMp4-cmd.Start() failed with '%s'\n", err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("ExtractAudioFromMp4-cmd.Run() failed with %s\n", err), err)
		return false
	}

	return true
}

// SplitMp4IntoChunks generate chunck parts of mp4
func (c *Client) SplitMp4IntoChunks(filename, dstFile, id string) bool {
	var stdBuffer bytes.Buffer
	os.MkdirAll(fmt.Sprintf("%schunks", dstFile), 0777)
	os.MkdirAll(fmt.Sprintf("%schunks_livepeer", dstFile), 0777)

	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", fmt.Sprintf("%s.mp4", filename), "-acodec", "aac", "-f", "segment", "-vcodec", "copy", "-reset_timestamps", "0", "-map", "0", fmt.Sprintf("%schunks/%s", dstFile, "output%03d.mp4"))

	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Start()
	if err != nil {
		utils.SendError(fmt.Sprintf("SplitMp4IntoChunks-cmd.Start() failed with '%s'\n", err), err)
		return false
	}

	err = cmd.Wait()
	if err != nil {
		utils.SendError(fmt.Sprintf("SplitMp4IntoChunks-cmd.Run() failed with %s\n", err), err)
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
