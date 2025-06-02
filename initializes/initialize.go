package initializes

import (
	"context"
	"fmt"
	"log"

	"github.com/youknow2509/crawl_vimeo/consts"
	"github.com/youknow2509/crawl_vimeo/global"
	"github.com/youknow2509/crawl_vimeo/services"
)

func Initialize() {
	initializeOsSystem()
	// Initialize Google services
	if err := initializeGoogleServices(); err != nil {
		log.Fatalf("Failed to initialize Google services: %v", err)
	}

	// Test M3U8 functionality
	url := services.GetPathM3U8("713179128")
	if url == "" {
		panic("Failed to get M3U8 path")
	}
	println("M3U8 Path:", url)
	// err := services.M3U8ToMP4(url, "~/data/1/output.mp4", global.OS_SYSTEM)
	// if err != nil {
	// 	panic(err)
	// }

	// test upload video
	ctx := context.Background()
	videoPath := "videos/video_test.mp4"
	title := "test upload api"
	description := "description test upload api"
	privacyStatus := consts.VIDEO_STATUS_PRIVATE
	response, err := global.YTB_SERVICE_ACTION.UploadVideoBase(ctx, title, description, privacyStatus, videoPath)
	if err != nil {
		log.Fatalf("Failed to upload video: %v", err)
	}
	fmt.Printf("Đã upload video thành công! Video ID: %s\n", response.Id)
	fmt.Printf("Link: https://www.youtube.com/watch?v=%s\n", response.Id)

	log.Println("Initialization completed successfully")
}
