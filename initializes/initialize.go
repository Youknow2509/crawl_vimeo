package initializes

import "github.com/youknow2509/crawl_vimeo/services"

func Initialize() {
	initializeOsSystem()	
	initializeYtb()

	url := services.GetPathM3U8("713179128")
	if url == "" {
		panic("Failed to get M3U8 path")
	}
	println("M3U8 Path:", url)
	// err := services.M3U8ToMP4(url, "~/data/1/output.mp4")
	// if err != nil {
	// 	panic(err)
	// }
}