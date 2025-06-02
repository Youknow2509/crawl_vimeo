package helpers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/youknow2509/crawl_vimeo/global"
	"github.com/youknow2509/crawl_vimeo/models"
	"github.com/youknow2509/crawl_vimeo/utils"
)

// processVideo downloads M3U8 and uploads to YouTube
func ProcessVideo(ctx context.Context, sectionVideo models.SectionVideo) (*models.SectionVideo, error) {
	// Generate safe filename
	safeTitle := utils.GenerateSafeFilename(sectionVideo.Title)
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