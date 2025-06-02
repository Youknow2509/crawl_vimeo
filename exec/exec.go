package exec

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/youknow2509/crawl_vimeo/helpers"
	"github.com/youknow2509/crawl_vimeo/initializes"
)

func Execute() {
	// Initialize all services and configurations
	initializes.Initialize()

	log.Println("Starting video processing execution...")

	// Read all JSON files from both directories
	sectionVideos, err := helpers.ReadAllSectionVideos()
	if err != nil {
		log.Fatalf("Failed to read section videos: %v", err)
	}

	log.Printf("Found %d videos to process", len(sectionVideos))

	// create json file save sectionVideos
	// Tạo file, nếu đã có sẽ ghi đè
	file, err := os.Create("output.json")
	if err != nil {
		fmt.Println("Không thể tạo file:", err)
		return
	}
	defer file.Close()

	// Ghi dữ liệu dưới dạng JSON, thụt lề đẹp
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(sectionVideos); err != nil {
		fmt.Println("Lỗi ghi JSON:", err)
		return
	}
	fmt.Println("Ghi dữ liệu ra file output.json thành công!")
}
	// // Initialize result lists
	// successList := []models.SectionVideo{}
	// errorList := []models.DataError{}

	// Load existing success and error files
// 	if existingSuccess, err := helpers.LoadExistingResults(consts.PATH_FILE_SUCCESS); err == nil {
// 		successList = existingSuccess
// 	}
// 	if existingErrors, err := helpers.LoadExistingErrors(consts.PATH_FILE_ERROR); err == nil {
// 		errorList = existingErrors
// 	}

// 	// Process each video
// 	ctx := context.Background()
// 	totalVideos := len(sectionVideos)
// 	processedCount := 0

// 	for _, sectionVideo := range sectionVideos {
// 		processedCount++

// 		log.Printf("Processing video %d/%d: %s", processedCount, totalVideos, sectionVideo.Title)

// 		// Process the video
// 		success, err := helpers.ProcessVideo(ctx, sectionVideo)
// 		if err != nil {
// 			log.Printf("Failed to process video %s: %v", sectionVideo.Title, err)

// 			// Add to error list
// 			errorData := models.DataError{
// 				CourseID:    sectionVideo.CourseID,
// 				SectionID:   sectionVideo.SectionID,
// 				SectionPath: sectionVideo.SectionPath,
// 				M3u8Link:    sectionVideo.M3u8Link,
// 			}
// 			errorList = append(errorList, errorData)
// 		} else {
// 			log.Printf("Successfully processed video %s with YouTube ID: %s", sectionVideo.Title, success.YtbVideoID)

// 			// Add to success list
// 			successList = append(successList, *success)
// 		}

// 		// Save results after each video
// 		if err := saveResults(successList, errorList); err != nil {
// 			log.Printf("Failed to save results: %v", err)
// 		}

// 		// Random sleep between 5-10 seconds
// 		sleepDuration := time.Duration(rand.Intn(6)+5) * time.Second
// 		log.Printf("Sleeping for %v before next video...", sleepDuration)
// 		time.Sleep(sleepDuration)

// 		// Show current dashboard
// 		showDashboard(len(successList), len(errorList), totalVideos, processedCount)
// 	}

// 	// Final dashboard
// 	log.Println("Processing completed!")
// 	showFinalDashboard(len(successList), len(errorList), totalVideos)
// }


// // saveResults saves both success and error results to files
// func saveResults(successList []models.SectionVideo, errorList []models.DataError) error {
// 	// Save success results
// 	successData, err := json.MarshalIndent(successList, "", "  ")
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal success data: %w", err)
// 	}

// 	if err := ioutil.WriteFile(consts.PATH_FILE_SUCCESS, successData, 0644); err != nil {
// 		return fmt.Errorf("failed to write success file: %w", err)
// 	}

// 	// Save error results
// 	errorData, err := json.MarshalIndent(errorList, "", "  ")
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal error data: %w", err)
// 	}

// 	if err := ioutil.WriteFile(consts.PATH_FILE_ERROR, errorData, 0644); err != nil {
// 		return fmt.Errorf("failed to write error file: %w", err)
// 	}

// 	return nil
// }

// // showDashboard displays current processing statistics
// func showDashboard(successCount, errorCount, totalVideos, processedCount int) {
// 	fmt.Println("\n" + strings.Repeat("=", 50))
// 	fmt.Println("PROCESSING DASHBOARD")
// 	fmt.Println(strings.Repeat("=", 50))
// 	fmt.Printf("Progress: %d/%d videos processed (%.1f%%)\n",
// 		processedCount, totalVideos, float64(processedCount)/float64(totalVideos)*100)
// 	fmt.Printf("Successful uploads: %d\n", successCount)
// 	fmt.Printf("Failed uploads: %d\n", errorCount)
// 	fmt.Printf("Remaining: %d\n", totalVideos-processedCount)
// 	fmt.Println(strings.Repeat("=", 50) + "\n")
// }

// // showFinalDashboard displays final processing statistics
// func showFinalDashboard(successCount, errorCount, totalVideos int) {
// 	fmt.Println("\n" + strings.Repeat("=", 60))
// 	fmt.Println("FINAL PROCESSING REPORT")
// 	fmt.Println(strings.Repeat("=", 60))
// 	fmt.Printf("Total videos processed: %d\n", totalVideos)
// 	fmt.Printf("Successful uploads: %d (%.1f%%)\n",
// 		successCount, float64(successCount)/float64(totalVideos)*100)
// 	fmt.Printf("Failed uploads: %d (%.1f%%)\n",
// 		errorCount, float64(errorCount)/float64(totalVideos)*100)
// 	fmt.Println()

// 	if successCount > 0 {
// 		fmt.Printf("Success! %d videos uploaded to YouTube\n", successCount)
// 	}
// 	if errorCount > 0 {
// 		fmt.Printf("%d videos failed to upload. Check error.json for details\n", errorCount)
// 	}

// 	fmt.Printf("Results saved to:\n")
// 	fmt.Printf("   - Success: %s\n", consts.PATH_FILE_SUCCESS)
// 	fmt.Printf("   - Errors: %s\n", consts.PATH_FILE_ERROR)
// 	fmt.Println(strings.Repeat("=", 60))
// }
