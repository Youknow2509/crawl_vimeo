package services

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/youknow2509/crawl_vimeo/consts"
)

/**
 *  Handle conversion from M3U8 to MP4
 *	@param inputUrl string - The URL of the M3U8 file to convert
 *	@param outPath string - The output path for the converted MP4 file
 *	@return error - Returns an error if the conversion fails
 */
func M3U8ToMP4(inputUrl string, outPath string, osSystem int) error {
	fmt.Printf("Converting M3U8 to MP4: %s -> %s\n", inputUrl, outPath)
	// create directory if not exists
	err := createDirIfNotExists(outPath, osSystem)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// ffmpeg command to convert M3U8 to MP4
	out, err := exec.Command("ffmpeg", "-i", inputUrl, "-c", "copy", outPath).Output()
    if err != nil {
        fmt.Println("Error:", err)
        return err
    }
    fmt.Println(string(out))
	return nil
}

// create dir to file if not exists
func createDirIfNotExists(pathFile string, osSystem int) error {
	var cmd *exec.Cmd
	dirPath := filepath.Dir(pathFile)
	if osSystem == consts.OS_WINDOWS {
		cmd = exec.Command("cmd", "/C", "mkdir", dirPath)
	} else {
		cmd = exec.Command("mkdir", "-p", dirPath)
	}
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create directory path %s\n%w", dirPath, err)
	}	
	fmt.Printf("Directory created: %s\n", dirPath)
	return nil
}
