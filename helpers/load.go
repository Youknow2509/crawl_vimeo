package helpers

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/youknow2509/crawl_vimeo/models"
)

// loadExistingResults loads existing success results
func LoadExistingResults(filePath string) ([]models.SectionVideo, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []models.SectionVideo{}, nil
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []models.SectionVideo{}, nil
	}

	var results []models.SectionVideo
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// loadExistingErrors loads existing error results
func LoadExistingErrors(filePath string) ([]models.DataError, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []models.DataError{}, nil
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []models.DataError{}, nil
	}

	var results []models.DataError
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	return results, nil
}