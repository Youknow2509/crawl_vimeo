package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/youknow2509/crawl_vimeo/global"
	"github.com/youknow2509/crawl_vimeo/models"
	"github.com/youknow2509/crawl_vimeo/utils"
)

// parseJSONFile parses a JSON file and converts it to SectionVideo model
func parseJSONFile(filePath string) (models.SectionVideo, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return models.SectionVideo{}, fmt.Errorf("failed to read file: %w", err)
	}

	var dataFile models.DataFile
	if err := json.Unmarshal(data, &dataFile); err != nil {
		return models.SectionVideo{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	// 
	contentUrl := dataFile.Data.Element.Content
	if contentUrl == "" {
		return models.SectionVideo{}, fmt.Errorf("content URL is empty in file: %s", filePath)
	}
	// extract id video from content URL
	videoID := utils.ExtractVideoIDVimeoFromURL(contentUrl)

	// Extract M3U8 link from video source
	m3u8Link, err := global.M3U8_SERVICE.GetPathM3U8(videoID)
	if err != nil {
		return models.SectionVideo{}, fmt.Errorf("failed to get M3U8 link: %w", err)
	}

	// Create description with chapters
	description := utils.CreateDescriptionWithChapters(dataFile.Data.Element.TimestampVideo)

	sectionVideo := models.SectionVideo{
		CourseID:    dataFile.Data.CourseID,
		SectionID:   dataFile.Data.ID,
		SectionPath: filePath,
		M3u8Link:    m3u8Link,
		YtbVideoID:  "", // Will be filled after upload
		Chapters:    dataFile.Data.Element.TimestampVideo,
		Title:       dataFile.Data.Title,
		Description: description,
	}

	return sectionVideo, nil
}