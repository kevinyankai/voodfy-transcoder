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
	ThumbsPreviewGenerator(string, string, string) bool
	VTTGenerator(string, string, string) bool
	ExtractAudioFromMp4(string, string) bool
	CheckIntegrityFromMp4s(string, string) bool
}

// Client instance of ffmpeg
type Client struct{}

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

// Run execute the ffmpeg job
func Run(cmd Commands, fnc string, args ...string) bool {
	switch fnc {
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
	}
	return false
}

// ExecCmd exec ffprobe command and return result of json.
func ExecCmd(fileName string) ([]byte, error) {
	return exec.Command("ffprobe",
		"-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", fileName).Output()
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
func (c *Client) RemoveAudioFromMP4(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer

	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-c", "copy", "-an", dstFile)
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

// Transcode90p low definition
func (c *Client) Transcode90p(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer

	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:90'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "100k", "-an", dstFile)

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

// Transcode144p low definition
func (c *Client) Transcode144p(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer

	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:90'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "100k", "-an", dstFile)

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
func (c *Client) Transcode240p(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:240'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "120k", "-an", dstFile)
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
func (c *Client) Transcode360p(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:360'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "284k", "-maxrate", "284k", "-bufsize", "568k", "-an", dstFile)
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
func (c *Client) Transcode480p(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer

	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:480'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "341k", "-maxrate", "341k", "-bufsize", "682k", "-an", dstFile)
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
func (c *Client) Transcode720p(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:720'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "765k", "-maxrate", "765k", "-bufsize", "1530k", "-an", dstFile)
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
func (c *Client) Transcode1080p(filename, dstFile string) bool {
	var stdBuffer bytes.Buffer
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-movflags", "faststart", "-vf", "scale='-2:1080'", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-b:v", "1579k", "-maxrate", "1579k", "-bufsize", "3158k", "-an", dstFile)

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
func (c *Client) ThumbsPreviewGenerator(filename, dstFile, duration string) bool {
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
func (c *Client) VTTGenerator(filename, dstFile, language string) bool {
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-f", "webvtt", fmt.Sprintf("%s_%s.vtt", dstFile, language))

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
func (c *Client) ExtractAudioFromMp4(filename, dstFile string) bool {
	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", filename, "-vn", "-acodec", "copy", dstFile)
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

// CheckIntegrityFromMp4s return a boolean about the duration of the video
func (c *Client) CheckIntegrityFromMp4s(source, output string) bool {
	s, _ := Execute(source)
	o, _ := Execute(output)

	if o.Format.Duration == s.Format.Duration {
		return true
	}

	return false
}
