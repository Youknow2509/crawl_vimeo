package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

// Structs để parse JSON
type PlaylistData struct {
	ClipID  string        `json:"clip_id"`
	BaseURL string        `json:"base_url"`
	Video   []VideoStream `json:"video"`
	Audio   []AudioStream `json:"audio"`
}

type VideoStream struct {
	ID                 string    `json:"id"`
	AvgID              string    `json:"avg_id"`
	BaseURL            string    `json:"base_url"`
	Format             string    `json:"format"`
	MimeType           string    `json:"mime_type"`
	Codecs             string    `json:"codecs"`
	Bitrate            int       `json:"bitrate"`
	AvgBitrate         int       `json:"avg_bitrate"`
	Duration           float64   `json:"duration"`
	Framerate          float64   `json:"framerate"`
	Width              int       `json:"width"`
	Height             int       `json:"height"`
	MaxSegmentDuration int       `json:"max_segment_duration"`
	InitSegment        string    `json:"init_segment"`
	InitSegmentURL     string    `json:"init_segment_url"`
	IndexSegment       string    `json:"index_segment"`
	Segments           []Segment `json:"segments"`
}

type AudioStream struct {
	ID           string    `json:"id"`
	AvgID        string    `json:"avg_id"`
	BaseURL      string    `json:"base_url"`
	Format       string    `json:"format"`
	MimeType     string    `json:"mime_type"`
	Codecs       string    `json:"codecs"`
	Bitrate      int       `json:"bitrate"`
	AvgBitrate   int       `json:"avg_bitrate"`
	Duration     float64   `json:"duration"`
	Channels     int       `json:"channels"`
	SampleRate   int       `json:"sample_rate"`
	AudioPrimary bool      `json:"audio_primary"`
	Segments     []Segment `json:"segments"`
}

type Segment struct {
	Start int    `json:"start"`
	End   int    `json:"end"`
	URL   string `json:"url"`
	Size  int    `json:"size"`
}

type AuthParams struct {
	Exp     string
	HMAC    string
	PathSig string
	R       string
}

type VimeoChunkGenerator struct {
	PlaylistData PlaylistData
	AuthParams   AuthParams
}

// NewVimeoChunkGenerator tạo generator mới
func NewVimeoChunkGenerator(playlistData PlaylistData, authParams AuthParams) *VimeoChunkGenerator {
	return &VimeoChunkGenerator{
		PlaylistData: playlistData,
		AuthParams:   authParams,
	}
}

// ExtractAuthParamsFromURL trích xuất các tham số auth từ URL
func ExtractAuthParamsFromURL(urlStr string) (AuthParams, error) {
	var params AuthParams

	// Extract expiry
	expRe := regexp.MustCompile(`exp=(\d+)`)
	if matches := expRe.FindStringSubmatch(urlStr); len(matches) > 1 {
		params.Exp = matches[1]
	}

	// Extract HMAC
	hmacRe := regexp.MustCompile(`hmac=([a-f0-9]+)`)
	if matches := hmacRe.FindStringSubmatch(urlStr); len(matches) > 1 {
		params.HMAC = matches[1]
	}

	// Extract pathsig
	pathsigRe := regexp.MustCompile(`pathsig=([^&]+)`)
	if matches := pathsigRe.FindStringSubmatch(urlStr); len(matches) > 1 {
		params.PathSig = matches[1]
	}

	// Extract r parameter
	rRe := regexp.MustCompile(`[?&]r=([^&]+)`)
	if matches := rRe.FindStringSubmatch(urlStr); len(matches) > 1 {
		params.R = matches[1]
	}

	return params, nil
}

// EncodeRange encode byte range theo format của Vimeo
func (vcg *VimeoChunkGenerator) EncodeRange(start, end int) string {
	rangeStr := fmt.Sprintf("range=%d-%d", start, end)
	encoded := base64.StdEncoding.EncodeToString([]byte(rangeStr))
	return encoded
}

// GenerateVideoChunkURLs sinh URLs cho tất cả chunks của một quality
func (vcg *VimeoChunkGenerator) GenerateVideoChunkURLs(qualityID string) ([]string, error) {
	var videoStream *VideoStream

	// Tìm video stream theo ID
	for i := range vcg.PlaylistData.Video {
		if vcg.PlaylistData.Video[i].ID == qualityID {
			videoStream = &vcg.PlaylistData.Video[i]
			break
		}
	}

	if videoStream == nil {
		return nil, fmt.Errorf("video stream with ID %s not found", qualityID)
	}

	if len(videoStream.Segments) == 0 {
		return nil, fmt.Errorf("no segments found for video stream %s", qualityID)
	}

	var urls []string
	for _, segment := range videoStream.Segments {
		url := vcg.buildSegmentURL(*videoStream, segment)
		urls = append(urls, url)
	}

	return urls, nil
}

// GenerateAudioChunkURLs sinh URLs cho audio chunks
func (vcg *VimeoChunkGenerator) GenerateAudioChunkURLs(audioID string) ([]string, error) {
	var audioStream *AudioStream

	// Tìm audio stream theo ID
	for i := range vcg.PlaylistData.Audio {
		if vcg.PlaylistData.Audio[i].ID == audioID {
			audioStream = &vcg.PlaylistData.Audio[i]
			break
		}
	}

	if audioStream == nil {
		return nil, fmt.Errorf("audio stream with ID %s not found", audioID)
	}

	if len(audioStream.Segments) == 0 {
		return nil, fmt.Errorf("no segments found for audio stream %s", audioID)
	}

	var urls []string
	for _, segment := range audioStream.Segments {
		url := vcg.buildAudioSegmentURL(*audioStream, segment)
		urls = append(urls, url)
	}

	return urls, nil
}

// buildSegmentURL xây dựng URL cho một video segment cụ thể
func (vcg *VimeoChunkGenerator) buildSegmentURL(videoStream VideoStream, segment Segment) string {
	baseURL := "https://vod-adaptive-ak.vimeocdn.com"

	if videoStream.Format == "dash" {
		// Format DASH sử dụng range/prot/
		encodedRange := vcg.EncodeRange(segment.Start, segment.End)

		url := fmt.Sprintf("%s/exp=%s~acl=%%2F%s%%2F%%2A~hmac=%s/%s/v2/%s%s/avf/%s.mp4?pathsig=%s&r=%s&range=%d-%d",
			baseURL,
			vcg.AuthParams.Exp,
			vcg.PlaylistData.ClipID,
			vcg.AuthParams.HMAC,
			vcg.PlaylistData.ClipID,
			videoStream.BaseURL,
			encodedRange,
			videoStream.ID,
			vcg.AuthParams.PathSig,
			vcg.AuthParams.R,
			segment.Start,
			segment.End,
		)
		return url
	} else if videoStream.Format == "mp42" {
		// Format MP42 sử dụng remux/avf/ (nếu có segment index)
		url := fmt.Sprintf("%s/exp=%s~acl=%%2F%s%%2F%%2A~hmac=%s/%s/v2/%ssegment-%d.m4s?pathsig=%s&r=%s",
			baseURL,
			vcg.AuthParams.Exp,
			vcg.PlaylistData.ClipID,
			vcg.AuthParams.HMAC,
			vcg.PlaylistData.ClipID,
			videoStream.BaseURL,
			segment.Start, // Sử dụng start như segment index
			vcg.AuthParams.PathSig,
			vcg.AuthParams.R,
		)
		return url
	}

	return ""
}

// buildAudioSegmentURL xây dựng URL cho audio segment
func (vcg *VimeoChunkGenerator) buildAudioSegmentURL(audioStream AudioStream, segment Segment) string {
	baseURL := "https://vod-adaptive-ak.vimeocdn.com"
	encodedRange := vcg.EncodeRange(segment.Start, segment.End)

	url := fmt.Sprintf("%s/exp=%s~acl=%%2F%s%%2F%%2A~hmac=%s/%s/v2/%s%s/avf/%s.mp4?pathsig=%s&r=%s&range=%d-%d",
		baseURL,
		vcg.AuthParams.Exp,
		vcg.PlaylistData.ClipID,
		vcg.AuthParams.HMAC,
		vcg.PlaylistData.ClipID,
		audioStream.BaseURL,
		encodedRange,
		audioStream.ID,
		vcg.AuthParams.PathSig,
		vcg.AuthParams.R,
		segment.Start,
		segment.End,
	)
	return url
}

// GetAvailableQualities trả về danh sách các quality có sẵn
func (vcg *VimeoChunkGenerator) GetAvailableQualities() []VideoQuality {
	var qualities []VideoQuality

	for _, video := range vcg.PlaylistData.Video {
		if len(video.Segments) > 0 {
			quality := VideoQuality{
				ID:         video.ID,
				Resolution: fmt.Sprintf("%dx%d", video.Width, video.Height),
				Bitrate:    video.Bitrate,
				Format:     video.Format,
			}
			qualities = append(qualities, quality)
		}
	}

	return qualities
}

type VideoQuality struct {
	ID         string
	Resolution string
	Bitrate    int
	Format     string
}

// DownloadVideo tải video với quality được chọn
func (vcg *VimeoChunkGenerator) DownloadVideo(qualityID, outputFile string) error {
	urls, err := vcg.GenerateVideoChunkURLs(qualityID)
	if err != nil {
		return err
	}

	fmt.Printf("Downloading video with %d chunks...\n", len(urls))

	// Tạo file output
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Download từng chunk
	for i, chunkURL := range urls {
		fmt.Printf("Downloading chunk %d/%d\n", i+1, len(urls))

		resp, err := http.Get(chunkURL)
		if err != nil {
			return fmt.Errorf("failed to download chunk %d: %v", i+1, err)
		}

		chunkData, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return fmt.Errorf("failed to read chunk %d: %v", i+1, err)
		}

		_, err = file.Write(chunkData)
		if err != nil {
			return fmt.Errorf("failed to write chunk %d: %v", i+1, err)
		}
	}

	fmt.Printf("Video saved to %s\n", outputFile)
	return nil
}

func main() {
	// Load playlist data
	jsonFile, err := ioutil.ReadFile("playlist_515c7aa7-e99b-4f16-8c56-85fe2f10bdc7.json")
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		return
	}

	var playlistData PlaylistData
	err = json.Unmarshal(jsonFile, &playlistData)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Extract auth params từ URL gốc
	originalURL := "https://vod-adaptive-ak.vimeocdn.com/exp=1748765517~acl=%2F515c7aa7-e99b-4f16-8c56-85fe2f10bdc7%2F%2A~hmac=4588bd5185287e24866fb0dec98d7d4a965e184b53da76ef29f51d54b24a8e08/515c7aa7-e99b-4f16-8c56-85fe2f10bdc7/v2/playlist/av/primary/prot/cXNyPTE/playlist.json?omit=av1-hevc&pathsig=8c953e4f~pV8o72cHIR1_KqkDDRC_zxO9L9zCbtDlQkBdK9d8wDY&qsr=1&r=dXM%3D&rh=28Enmi"

	authParams, err := ExtractAuthParamsFromURL(originalURL)
	if err != nil {
		fmt.Printf("Error extracting auth params: %v\n", err)
		return
	}

	// Tạo generator
	generator := NewVimeoChunkGenerator(playlistData, authParams)

	// Hiển thị các quality có sẵn
	qualities := generator.GetAvailableQualities()
	fmt.Println("Available video qualities:")
	for i, quality := range qualities {
		fmt.Printf("%d. %s (%s) - %d kbps - %s\n",
			i+1, quality.Resolution, quality.ID, quality.Bitrate/1000, quality.Format)
	}

	// Sinh URLs cho video quality có segments (960x540)
	if len(qualities) > 0 {
		videoURLs, err := generator.GenerateVideoChunkURLs(qualities[0].ID)
		if err != nil {
			fmt.Printf("Error generating video URLs: %v\n", err)
			return
		}

		fmt.Printf("\nGenerated %d video chunk URLs:\n", len(videoURLs))
		for i, url := range videoURLs[:3] { // In 3 URL đầu tiên
			fmt.Printf("Chunk %d: %s\n", i+1, url)
		}

		// Uncomment để download video
		// err = generator.DownloadVideo(qualities[0].ID, "output_video.mp4")
		// if err != nil {
		//     fmt.Printf("Error downloading video: %v\n", err)
		// }
	}
}
