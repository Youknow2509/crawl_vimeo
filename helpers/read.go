package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/youknow2509/crawl_vimeo/consts"
	"github.com/youknow2509/crawl_vimeo/models"
	"github.com/youknow2509/crawl_vimeo/utils"
)

// readAllSectionVideos reads all JSON files and converts them to SectionVideo models
func ReadAllSectionVideos() ([]models.SectionVideo, error) {
	var sectionVideos []models.SectionVideo

	// Read from TOEIC directory
	toeicVideos, err := readVideosFromDirectory(consts.BASE_PATH_SECTION_TOEIC, "[TOEIC]")
	if err != nil {
		log.Printf("Warning: Failed to read TOEIC videos: %v", err)
	} else {
		sectionVideos = append(sectionVideos, toeicVideos...)
	}

	// Read from IELTS directory
	ieltsVideos, err := readVideosFromDirectory(consts.BASE_PATH_SECTION_IELTS, "[IELTS]")
	if err != nil {
		log.Printf("Warning: Failed to read IELTS videos: %v", err)
	} else {
		sectionVideos = append(sectionVideos, ieltsVideos...)
	}

	return sectionVideos, nil
}

// readVideosFromDirectory reads all JSON files from a directory
func readVideosFromDirectory(dirPath string, typeVideo string) ([]models.SectionVideo, error) {
	var sectionVideos []models.SectionVideo

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		sectionVideo, err := parseJSONFile(filePath)
		if err != nil {
			log.Printf("Warning: Failed to parse file %s: %v", filePath, err)
			continue
		}
		// Set the title of video
		sectionVideo.Title = fmt.Sprintf("%s %s", typeVideo, sectionVideo.Title)
		// Set description of video
		sectionVideo.Description = utils.CreateDescriptionWithChapters(sectionVideo.Chapters)
		// add to list
		sectionVideos = append(sectionVideos, sectionVideo)
		// sleep to avoid rate limit
	}

	return sectionVideos, nil
}

