package initializes

import (
	"log"

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

	log.Println("Initialization completed successfully")
}
