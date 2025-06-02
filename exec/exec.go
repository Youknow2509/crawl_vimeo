package exec

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/youknow2509/crawl_vimeo/consts"
	"github.com/youknow2509/crawl_vimeo/global"
	"github.com/youknow2509/crawl_vimeo/initializes"
	"github.com/youknow2509/crawl_vimeo/models"
	"github.com/youknow2509/crawl_vimeo/utils"
)

func Execute() {
	// Initialize all services and configurations
	initializes.Initialize()

	log.Println("Starting video processing execution...")

	// Read all JSON files from both directories
	sectionVideos, err := readAllSectionVideos()
	if err != nil {
		log.Fatalf("Failed to read section videos: %v", err)
	}

	log.Printf("Found %d videos to process", len(sectionVideos))

	// Initialize result lists
	successList := []models.SectionVideo{}
	errorList := []models.DataError{}

	// Load existing success and error files
	if existingSuccess, err := loadExistingResults(consts.PATH_FILE_SUCCESS); err == nil {
		successList = existingSuccess
	}
	if existingErrors, err := loadExistingErrors(consts.PATH_FILE_ERROR); err == nil {
		errorList = existingErrors
	}

	// Process each video
	ctx := context.Background()
	totalVideos := len(sectionVideos)
	processedCount := 0

	for _, sectionVideo := range sectionVideos {
		processedCount++

		log.Printf("Processing video %d/%d: %s", processedCount, totalVideos, sectionVideo.Title)

		// Process the video
		success, err := processVideo(ctx, sectionVideo)
		if err != nil {
			log.Printf("Failed to process video %s: %v", sectionVideo.Title, err)

			// Add to error list
			errorData := models.DataError{
				CourseID:    sectionVideo.CourseID,
				SectionID:   sectionVideo.SectionID,
				SectionPath: sectionVideo.SectionPath,
				M3u8Link:    sectionVideo.M3u8Link,
			}
			errorList = append(errorList, errorData)
		} else {
			log.Printf("Successfully processed video %s with YouTube ID: %s", sectionVideo.Title, success.YtbVideoID)

			// Add to success list
			successList = append(successList, *success)
		}

		// Save results after each video
		if err := saveResults(successList, errorList); err != nil {
			log.Printf("Failed to save results: %v", err)
		}

		// Random sleep between 5-10 seconds
		sleepDuration := time.Duration(rand.Intn(6)+5) * time.Second
		log.Printf("Sleeping for %v before next video...", sleepDuration)
		time.Sleep(sleepDuration)

		// Show current dashboard
		showDashboard(len(successList), len(errorList), totalVideos, processedCount)
	}

	// Final dashboard
	log.Println("Processing completed!")
	showFinalDashboard(len(successList), len(errorList), totalVideos)
}

// readAllSectionVideos reads all JSON files and converts them to SectionVideo models
func readAllSectionVideos() ([]models.SectionVideo, error) {
	var sectionVideos []models.SectionVideo

	// Read from IELTS directory
	ieltsVideos, err := readVideosFromDirectory(consts.BASE_PATH_SECTION_IELTS)
	if err != nil {
		log.Printf("Warning: Failed to read IELTS videos: %v", err)
	} else {
		sectionVideos = append(sectionVideos, ieltsVideos...)
	}

	// Read from TOEIC directory
	toeicVideos, err := readVideosFromDirectory(consts.BASE_PATH_SECTION_TOEIC)
	if err != nil {
		log.Printf("Warning: Failed to read TOEIC videos: %v", err)
	} else {
		sectionVideos = append(sectionVideos, toeicVideos...)
	}

	return sectionVideos, nil
}

// readVideosFromDirectory reads all JSON files from a directory
func readVideosFromDirectory(dirPath string) ([]models.SectionVideo, error) {
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

		sectionVideos = append(sectionVideos, sectionVideo)
	}

	return sectionVideos, nil
}

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

	// Extract M3U8 link from video source
	m3u8Link := dataFile.Data.Element.VideoSource.URL
	if m3u8Link == "" {
		m3u8Link = dataFile.Data.Element.VideoSource.WebURL
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

// processVideo downloads M3U8 and uploads to YouTube
func processVideo(ctx context.Context, sectionVideo models.SectionVideo) (*models.SectionVideo, error) {
	// Generate safe filename
	safeTitle := generateSafeFilename(sectionVideo.Title)
	videoOutputPath := fmt.Sprintf("videos/%s.mp4", safeTitle)

	// Download M3U8 to MP4
	log.Printf("Downloading M3U8: %s", sectionVideo.M3u8Link)
	if err := global.M3U8_SERVICE.M3U8ToMP4(sectionVideo.M3u8Link, videoOutputPath); err != nil {
		return nil, fmt.Errorf("failed to download M3U8: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(videoOutputPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("video file was not created: %s", videoOutputPath)
	}

	// Upload to YouTube
	log.Printf("Uploading to YouTube: %s", sectionVideo.Title)
	video, err := global.YTB_SERVICE_ACTION.UploadVideoBase(
		ctx,
		sectionVideo.Title,
		sectionVideo.Description,
		"private", // Privacy status
		videoOutputPath,
	)
	if err != nil {
		// Clean up the downloaded file
		os.Remove(videoOutputPath)
		return nil, fmt.Errorf("failed to upload to YouTube: %w", err)
	}

	// Update section video with YouTube ID
	sectionVideo.YtbVideoID = video.Id

	// Clean up the downloaded file
	if err := os.Remove(videoOutputPath); err != nil {
		log.Printf("Warning: Failed to remove video file %s: %v", videoOutputPath, err)
	}

	return &sectionVideo, nil
}

// generateSafeFilename creates a safe filename from title
func generateSafeFilename(title string) string {
	// Remove or replace invalid characters
	reg := regexp.MustCompile(`[<>:"/\\|?*]`)
	safe := reg.ReplaceAllString(title, "_")

	// Limit length
	if len(safe) > 100 {
		safe = safe[:100]
	}

	// Remove spaces and replace with underscores
	safe = strings.ReplaceAll(safe, " ", "_")

	return safe
}

// loadExistingResults loads existing success results
func loadExistingResults(filePath string) ([]models.SectionVideo, error) {
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
func loadExistingErrors(filePath string) ([]models.DataError, error) {
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

// saveResults saves both success and error results to files
func saveResults(successList []models.SectionVideo, errorList []models.DataError) error {
	// Save success results
	successData, err := json.MarshalIndent(successList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal success data: %w", err)
	}

	if err := ioutil.WriteFile(consts.PATH_FILE_SUCCESS, successData, 0644); err != nil {
		return fmt.Errorf("failed to write success file: %w", err)
	}

	// Save error results
	errorData, err := json.MarshalIndent(errorList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal error data: %w", err)
	}

	if err := ioutil.WriteFile(consts.PATH_FILE_ERROR, errorData, 0644); err != nil {
		return fmt.Errorf("failed to write error file: %w", err)
	}

	return nil
}

// showDashboard displays current processing statistics
func showDashboard(successCount, errorCount, totalVideos, processedCount int) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("PROCESSING DASHBOARD")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Progress: %d/%d videos processed (%.1f%%)\n",
		processedCount, totalVideos, float64(processedCount)/float64(totalVideos)*100)
	fmt.Printf("Successful uploads: %d\n", successCount)
	fmt.Printf("Failed uploads: %d\n", errorCount)
	fmt.Printf("Remaining: %d\n", totalVideos-processedCount)
	fmt.Println(strings.Repeat("=", 50) + "\n")
}

// showFinalDashboard displays final processing statistics
func showFinalDashboard(successCount, errorCount, totalVideos int) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("FINAL PROCESSING REPORT")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Total videos processed: %d\n", totalVideos)
	fmt.Printf("Successful uploads: %d (%.1f%%)\n",
		successCount, float64(successCount)/float64(totalVideos)*100)
	fmt.Printf("Failed uploads: %d (%.1f%%)\n",
		errorCount, float64(errorCount)/float64(totalVideos)*100)
	fmt.Println()

	if successCount > 0 {
		fmt.Printf("Success! %d videos uploaded to YouTube\n", successCount)
	}
	if errorCount > 0 {
		fmt.Printf("%d videos failed to upload. Check error.json for details\n", errorCount)
	}

	fmt.Printf("Results saved to:\n")
	fmt.Printf("   - Success: %s\n", consts.PATH_FILE_SUCCESS)
	fmt.Printf("   - Errors: %s\n", consts.PATH_FILE_ERROR)
	fmt.Println(strings.Repeat("=", 60))
}
