package impl

import (
	"github.com/youknow2509/crawl_vimeo/consts"
	"github.com/youknow2509/crawl_vimeo/services"
	"github.com/youknow2509/crawl_vimeo/utils"
)

// struct
type M3U8Service struct {
	osExec services.IOSExecutor
}

// #############################################

// GetPathM3U8 implements services.IM3u8.
func (m *M3U8Service) GetPathM3U8(videoID string) (string, error) {
	// Create url for the video ID
	url := consts.PLAYER_VIMEO_URL + videoID
	// Get HTML content from the URL
	htmlContent, err := utils.HttpGet(url)
	if err != nil {
		return "", err
	}
	// Extract the playerConfig from the HTML content
	playerConfig, err := utils.ExtractPlayerConfigFromHTML(htmlContent)
	if err != nil {
		return "", err
	}
	// Parse the playerConfig JSON string into a map
	config, err := utils.ParsePlayerConfig(playerConfig)
	if err != nil {
		return "", err
	}
	// Extract the auth token from the parsed config
	res, err := utils.GetM3U8PathFromPlayerConfig(config)
	if err != nil {
		return "", err
	}
	// Return the M3U8 path
	return res, nil
}

// M3U8ToMP4 implements services.IM3u8.
func (m *M3U8Service) M3U8ToMP4(inputUrl string, outPath string) error {
	err := m.osExec.CreateDirIfNotExists(outPath)
	if err != nil {
		return err
	}

	// ffmpeg command to convert M3U8 to MP4
	err = m.osExec.ConvertM3U8ToMP4(inputUrl, outPath)
	if err != nil {
		return err
	}
	return nil
}

// #############################################

// implementation of IM3u8 interface
func NewM3u8Service(osExec services.IOSExecutor) services.IM3u8 {
	return &M3U8Service{
		osExec: osExec,
	}
}
