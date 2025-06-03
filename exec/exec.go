package exec

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/youknow2509/crawl_vimeo/consts"
	"github.com/youknow2509/crawl_vimeo/helpers"
	"github.com/youknow2509/crawl_vimeo/initializes"
	"github.com/youknow2509/crawl_vimeo/models"
)

func Execute() {
	// Initialize all services and configurations
	initializes.Initialize()

	log.Println("Starting video processing execution...")

	// load `data/data.json` file
	fileData, err := os.Open("data/data.json")
	if err != nil {
		log.Fatalf("Failed to open data file: %v", err)
		return
	}
	defer fileData.Close()
	// create file success.json
	fileSuccess, err := os.Create("data/success.json")
	if err != nil {
		log.Fatalf("Failed to create success file: %v", err)
		return
	}
	defer fileSuccess.Close()
	// create file error.json
	fileError, err := os.Create("data/error.json")
	if err != nil {
		log.Fatalf("Failed to create error file: %v", err)
		return
	}
	defer fileError.Close()
	// Decode JSON data into a slice of SectionVideo
	var sectionVideos []models.SectionVideo
	if err := json.NewDecoder(fileData).Decode(&sectionVideos); err != nil {
		log.Fatalf("Failed to decode JSON data: %v", err)
		return
	}
	// 
	fmt.Println("Loaded", len(sectionVideos), "videos from data.json")
	ctx := context.Background()
	errorList := make([]models.DataError, 0)
	successList := make([]models.SectionVideo, 0)
	// loop
	for i, sectionVideo := range sectionVideos {
		i += 1
		fmt.Printf("Processing video %d/%d: %s\n", i, len(sectionVideos), sectionVideo.Title)
		res, err := helpers.ProcessVideo(ctx, sectionVideo)
		if err != nil {
			log.Printf("Error processing video %d: %v", i, err)
			
			errorData := models.DataError{
				CourseID: sectionVideo.CourseID,
				SectionID: sectionVideo.SectionID,
				M3u8Link: sectionVideo.M3u8Link,
			}
			// Write error data to error.json
			errorList = append(errorList, errorData)
		} else {
			log.Printf("Successfully processed video %d: %s", i, res.Title)
			successList = append(successList, *res)
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
		showDashboard(len(successList), len(errorList), len(sectionVideos), i)
	}

	// Final dashboard
	log.Println("Processing completed!")
	showFinalDashboard(len(successList), len(errorList), len(sectionVideos))
}

// 

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
