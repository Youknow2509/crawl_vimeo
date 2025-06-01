package services

import (
	"fmt"
	"github.com/youknow2509/crawl_vimeo/consts"
	"github.com/youknow2509/crawl_vimeo/utils"
)

/**
 * GetPathM3U8 retrieves the M3U8 path for a given video ID.
 * @param videoID The ID of the video.
 * @return The M3U8 path as a string.
 */
func GetPathM3U8(videoID string) string {

	url := consts.PLAYER_VIMEO_URL + videoID
	htmlContent, err := HttpGet(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return ""
	}
	// Extract the playerConfig from the HTML content
	playerConfig, err := utils.ExtractPlayerConfigFromHTML(htmlContent)
	if err != nil {
		fmt.Println("Error extracting playerConfig:", err)
		return ""
	}
	// fmt.Println("Extracted playerConfig:", playerConfig)
	// Parse the playerConfig JSON string into a map
	config, err := utils.ParsePlayerConfig(playerConfig)
	if err != nil {
		fmt.Println("Error parsing playerConfig:", err)
		return ""
	}
	// Extract the auth token from the parsed config
	res, err := utils.GetM3U8PathFromPlayerConfig(config)
	if err != nil {
		fmt.Println("Error extracting M3U8 path:", err)
		return ""
	}
	return res
}
