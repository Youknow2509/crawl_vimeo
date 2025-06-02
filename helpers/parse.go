package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

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
	m3u8Link := strings.ReplaceAll(dataFile.Data.Element.VideoSource.Url, "api/video-encrypt-m3u8", "video")
	if m3u8Link == "" {
		return models.SectionVideo{}, fmt.Errorf("m3u8 link is empty")
	}
	
	// Create description with chapters
	description := utils.CreateDescriptionWithChapters(dataFile.Data.Element.TimestampVideo)
	// Link tài liệ
	if dataFile.Data.SubElement.MediaPath != "" {
		description += fmt.Sprintf("Tài liệu: %s\n \n", dataFile.Data.SubElement.MediaPath)
		fmt.Println("Link tài liệu:", dataFile.Data.SubElement.MediaPath)
	}

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