package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Specification struct to save info provided by ffprobe
type Specification struct {
	ID     string `json:"id"`
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

// Save add specification to redis
func (s *Specification) Save() bool {
	InitDB()

	m, err := s.MarshalBinary()

	if err != nil {
		log.Println("err", err)
	}

	if err := db.Redis.Set(fmt.Sprintf("specification_%s", string(s.ID)), m, 0).Err(); err != nil {
		fmt.Printf("Unable to store example struct into redis due to: %s \n", err)
	}
	return true
}

// Update update specification to redis
func (s *Specification) Update() {
	InitDB()

	m, err := s.MarshalBinary()

	if err != nil {
		log.Println("err", err)
	}

	if err := db.Redis.Set(fmt.Sprintf("specification_%s", string(s.ID)), m, 0).Err(); err != nil {
		fmt.Printf("Unable to store example struct into redis due to: %s \n", err)
	}
}

// Get return a specification save on redis
func (s *Specification) Get() {
	var key string
	InitDB()

	keys, _ := db.Redis.Keys("*").Result()

	for _, k := range keys {
		if strings.Contains(k, "specification_") {
			key = k
		}
	}

	cacheData, cacheErr := db.Redis.Get(key).Result()

	if cacheErr == nil {
		if err := s.UnmarshalBinary([]byte(cacheData)); err != nil {
			fmt.Printf("Unable to unmarshal data into the new example struct due to: %s \n", err)
		}
	}
}

// MarshalBinary retrieve resource from binary
func (s *Specification) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalBinary bind specification save on redis
func (s *Specification) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, s); err != nil {
		return err
	}

	return nil
}
