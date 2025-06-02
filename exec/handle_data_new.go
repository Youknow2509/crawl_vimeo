package exec

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/youknow2509/crawl_vimeo/consts"
	"github.com/youknow2509/crawl_vimeo/models"
)

func CountDataWithDataNew() {
	log.Println("Counting videos in data.json and data_new.json...")

	// Load data.json
	dataVideos, err := loadDataFromFile("data/data.json")
	if err != nil {
		log.Printf("Warning: Could not load data.json (file may not exist): %v", err)
		return
	}

	// Load data_new.json
	dataNewVideos, err := loadDataFromFile("data/data_new.json")
	if err != nil {
		log.Printf("Warning: Could not load data_new.json (file may not exist): %v", err)
		return
	}

	if len(dataVideos) == len(dataNewVideos) {
		log.Println("data.json == data_new.json == %d", len(dataVideos))
		return
	} else {
		log.Printf("data.json:: %d --- data_new.json:: %d", len(dataVideos), len(dataNewVideos))
	}
}

func CountSuccessWithSuccessNew() {
	log.Println("Counting successful videos...")

	// Load success_new.json
	successNewVideos, err := loadDataFromFile("data/success_new.json")
	if err != nil {
		log.Printf("Warning: Could not load success_new.json (file may not exist): %v", err)
		return
	}

	// load success.json
	successVideos, err := loadDataFromFile(consts.PATH_FILE_SUCCESS)
	if err != nil {
		log.Printf("Warning: Could not load success.json (file may not exist): %v", err)
		return
	}
	// 
	if len(successNewVideos) == len(successVideos) {
		log.Println("success_new.json == success.json == %d", len(successNewVideos))
		return
	} else {
		log.Printf("success_new.json:: %d --- success.json:: %d", len(successNewVideos), len(successVideos))
	}
}

func HandleFileData() {
	log.Println("Creating new data file (data_new.json)...")

	// Load success.json
	successVideos, err := loadDataFromFile(consts.PATH_FILE_SUCCESS)
	if err != nil {
		log.Printf("Warning: Could not load success.json (file may not exist): %v", err)
		successVideos = []models.SectionVideo{}
	}

	// var
	successDataNew := []models.SectionVideo{}
	data_new := []models.SectionVideo{}

	// Create map of successful videos (course_id + section_id as key)
	successMap := make(map[string]models.SectionVideo)
	for _, video := range successVideos {
		if video.YtbVideoID != "" {
			key := createVideoKey(video.CourseID, video.SectionID)
			successMap[key] = video
			successDataNew = append(successDataNew, video)
		}
	}
	// write successDataNew to success_new.json
	if err := saveDataToFile(successDataNew, "data/success_new.json"); err != nil {
		fmt.Errorf("failed to save success_new.json: %w", err)
		return
	}

	// Load data.json
	dataAll, err := loadDataFromFile("data/data.json")
	if err != nil {
		fmt.Errorf("failed to load data.json: %w", err)
		return
	}

	for _, video := range dataAll {
		key := createVideoKey(video.CourseID, video.SectionID)
		if _, exists := successMap[key]; !exists {
			// If the video is not in successMap, add it to data_new
			data_new = append(data_new, video)
		}
	}

	// Save data_new to data_new.json
	if err := saveDataToFile(data_new, "data/data_new.json"); err != nil {
		fmt.Printf("failed to save data_new.json: %w", err)
		return
	}
}

// loadDataFromFile loads SectionVideo data from a JSON file
func loadDataFromFile(filepath string) ([]models.SectionVideo, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer file.Close()

	var videos []models.SectionVideo
	if err := json.NewDecoder(file).Decode(&videos); err != nil {
		return nil, fmt.Errorf("failed to decode JSON from %s: %w", filepath, err)
	}
	return videos, nil
}

// saveDataToFile saves SectionVideo data to a JSON file
func saveDataToFile(videos []models.SectionVideo, filepath string) error {
	data, err := json.MarshalIndent(videos, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	if err := ioutil.WriteFile(filepath, data, 0644); err != nil {
		fmt.Errorf("failed to write file %s: %w", filepath, err)
		return err
	}

	return nil
}

// createVideoKey creates a unique key from course_id and section_id
func createVideoKey(courseID, sectionID int) string {
	return fmt.Sprintf("%d_%d", courseID, sectionID)
}
