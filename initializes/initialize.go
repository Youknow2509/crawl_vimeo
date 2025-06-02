package initializes

import (
	"log"

	"github.com/youknow2509/crawl_vimeo/global"
)

func Initialize() {
	// Initialize OS system first
	initializeOsSystem()

	// Validate OS initialization
	if err := ValidateOSInitialization(); err != nil {
		log.Fatalf("OS system validation failed: %v", err)
	}

	// Initialize Google services
	// if err := initializeGoogleServices(); err != nil {
	// 	log.Fatalf("Failed to initialize Google services: %v", err)
	// }

	// Initialize M3U8 service
	initializeM3u8()

	// Test M3U8 functionality
	url, err := global.M3U8_SERVICE.GetPathM3U8("713179128")
	if err != nil {
		log.Printf("Failed to get M3U8 path: %v", err)
		return
	}
	log.Println("M3U8 Path:", url)

	// Uncomment to test M3U8 to MP4 conversion using global OS executor
	// err = global.M3U8_SERVICE.M3U8ToMP4(url, "videos/output.mp4")
	// if err != nil {
	// 	log.Printf("Failed to convert M3U8 to MP4: %v", err)
	// }

	// test upload video
	// ctx := context.Background()
	// videoPath := "videos/video_test.mp4"
	// title := "test upload api"
	// description := "description test upload api"
	// privacyStatus := consts.VIDEO_STATUS_PRIVATE
	// response, err := global.YTB_SERVICE_ACTION.UploadVideoBase(ctx, title, description, privacyStatus, videoPath)
	// if err != nil {
	//     log.Fatalf("Failed to upload video: %v", err)
	// }
	// fmt.Printf("Đã upload video thành công! Video ID: %s\n", response.Id)
	// fmt.Printf("Link: https://www.youtube.com/watch?v=%s\n", response.Id)

	// err = global.OS_EXECUTOR.DeleteFile("videos/a.txt")
	// if err != nil {
	// 	log.Fatalln(err)
	// 	return
	// }

	log.Println("Initialization completed successfully")
}
