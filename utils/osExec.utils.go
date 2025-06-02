package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/youknow2509/crawl_vimeo/consts"
)

// OSCommand represents a command for different operating systems
type OSCommand struct {
	Windows []string
	MacOS   []string
	Linux   []string
	Ubuntu  []string
}

// OSExecutor handles command execution across different operating systems
type OSExecutor struct {
	osType int
}

// NewOSExecutor creates a new OS executor
func NewOSExecutor(osType int) *OSExecutor {
	return &OSExecutor{osType: osType}
}

// GetCurrentOS detects the current operating system
func GetCurrentOS() int {
	switch runtime.GOOS {
	case "windows":
		return consts.OS_WINDOWS
	case "darwin":
		return consts.OS_MACOS
	case "linux":
		// Check if it's Ubuntu
		if isUbuntu() {
			return consts.OS_UBUNTU
		}
		return consts.OS_LINUX
	default:
		return consts.OS_LINUX // Default to Linux
	}
}

// isUbuntu checks if the current Linux distribution is Ubuntu
func isUbuntu() bool {
	if _, err := os.Stat("/etc/lsb-release"); err == nil {
		content, err := os.ReadFile("/etc/lsb-release")
		if err == nil {
			return strings.Contains(string(content), "Ubuntu")
		}
	}
	return false
}

// ExecuteCommand executes a command based on the operating system
func (oe *OSExecutor) ExecuteCommand(cmdMap OSCommand) ([]byte, error) {
	var cmd []string

	switch oe.osType {
	case consts.OS_WINDOWS:
		cmd = cmdMap.Windows
	case consts.OS_MACOS:
		cmd = cmdMap.MacOS
	case consts.OS_LINUX:
		cmd = cmdMap.Linux
	case consts.OS_UBUNTU:
		if len(cmdMap.Ubuntu) > 0 {
			cmd = cmdMap.Ubuntu
		} else {
			cmd = cmdMap.Linux // Fallback to Linux
		}
	default:
		return nil, fmt.Errorf("unsupported operating system: %d", oe.osType)
	}

	if len(cmd) == 0 {
		return nil, fmt.Errorf("no command defined for OS type: %d", oe.osType)
	}

	if len(cmd) == 1 {
		return exec.Command(cmd[0]).Output()
	}

	return exec.Command(cmd[0], cmd[1:]...).Output()
}

// CreateDirectory creates a directory path across different OS
func (oe *OSExecutor) CreateDirectory(dirPath string) error {
	cmdMap := OSCommand{
		Windows: []string{"cmd", "/C", "mkdir", dirPath},
		MacOS:   []string{"mkdir", "-p", dirPath},
		Linux:   []string{"mkdir", "-p", dirPath},
		Ubuntu:  []string{"mkdir", "-p", dirPath},
	}

	// First try using Go's built-in function
	if err := os.MkdirAll(dirPath, 0755); err == nil {
		return nil
	}

	// Fallback to OS-specific command
	_, err := oe.ExecuteCommand(cmdMap)
	return err
}

// RemoveDirectory removes a directory across different OS
func (oe *OSExecutor) RemoveDirectory(dirPath string) error {
	cmdMap := OSCommand{
		Windows: []string{"cmd", "/C", "rmdir", "/S", "/Q", dirPath},
		MacOS:   []string{"rm", "-rf", dirPath},
		Linux:   []string{"rm", "-rf", dirPath},
		Ubuntu:  []string{"rm", "-rf", dirPath},
	}

	_, err := oe.ExecuteCommand(cmdMap)
	return err
}

// CopyFile copies a file across different OS
func (oe *OSExecutor) CopyFile(source, destination string) error {
	cmdMap := OSCommand{
		Windows: []string{"cmd", "/C", "copy", source, destination},
		MacOS:   []string{"cp", source, destination},
		Linux:   []string{"cp", source, destination},
		Ubuntu:  []string{"cp", source, destination},
	}

	_, err := oe.ExecuteCommand(cmdMap)
	return err
}

// MoveFile moves a file across different OS
func (oe *OSExecutor) MoveFile(source, destination string) error {
	cmdMap := OSCommand{
		Windows: []string{"cmd", "/C", "move", source, destination},
		MacOS:   []string{"mv", source, destination},
		Linux:   []string{"mv", source, destination},
		Ubuntu:  []string{"mv", source, destination},
	}

	_, err := oe.ExecuteCommand(cmdMap)
	return err
}

// CheckFFmpeg checks if FFmpeg is available on the system
func (oe *OSExecutor) CheckFFmpeg() (bool, error) {
	cmdMap := OSCommand{
		Windows: []string{"where", "ffmpeg"},
		MacOS:   []string{"which", "ffmpeg"},
		Linux:   []string{"which", "ffmpeg"},
		Ubuntu:  []string{"which", "ffmpeg"},
	}

	_, err := oe.ExecuteCommand(cmdMap)
	return err == nil, err
}

// InstallFFmpeg attempts to install FFmpeg on the system
func (oe *OSExecutor) InstallFFmpeg() error {
	cmdMap := OSCommand{
		Windows: []string{"winget", "install", "FFmpeg"},
		MacOS:   []string{"brew", "install", "ffmpeg"},
		Linux:   []string{"sudo", "apt-get", "install", "-y", "ffmpeg"},
		Ubuntu:  []string{"sudo", "apt-get", "install", "-y", "ffmpeg"},
	}

	_, err := oe.ExecuteCommand(cmdMap)
	return err
}

// ConvertM3U8ToMP4 converts M3U8 to MP4 using FFmpeg
func (oe *OSExecutor) ConvertM3U8ToMP4(inputUrl, outputPath string) error {
	// Create directory if not exists
	dirPath := filepath.Dir(outputPath)
	if err := oe.CreateDirectory(dirPath); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if FFmpeg is available
	available, err := oe.CheckFFmpeg()
	if !available {
		return fmt.Errorf("FFmpeg is not available on this system: %w", err)
	}

	// FFmpeg command (same across all OS)
	cmd := []string{"ffmpeg", "-i", inputUrl, "-c", "copy", outputPath}

	cmdMap := OSCommand{
		Windows: cmd,
		MacOS:   cmd,
		Linux:   cmd,
		Ubuntu:  cmd,
	}

	fmt.Printf("Converting M3U8 to MP4: %s -> %s\n", inputUrl, outputPath)

	output, err := oe.ExecuteCommand(cmdMap)
	if err != nil {
		return fmt.Errorf("FFmpeg conversion failed: %w", err)
	}

	fmt.Printf("Conversion completed: %s\n", string(output))
	return nil
}

// DeleteFile deletes a file across different OS with improved error handling
func (oe *OSExecutor) DeleteFile(filePath string) error {
	// Validate input
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Check if file exists first
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File doesn't exist, consider it as success
		fmt.Printf("File doesn't exist: %s\n", filePath)
		return nil
	}

	// Try Go's built-in function first
	if err := os.Remove(filePath); err == nil {
		fmt.Printf("File deleted successfully: %s\n", filePath)
		return nil
	}

	// Fallback to OS-specific commands
	cmdMap := OSCommand{
		Windows: []string{"cmd", "/C", "del", "/F", "/Q", filePath},
		MacOS:   []string{"rm", "-f", filePath},
		Linux:   []string{"rm", "-f", filePath},
		Ubuntu:  []string{"rm", "-f", filePath},
	}

	output, err := oe.ExecuteCommand(cmdMap)
	if err != nil {
		return fmt.Errorf("failed to delete file '%s': %w\nOutput: %s", filePath, err, string(output))
	}

	fmt.Printf("File deleted via command: %s\n", filePath)
	return nil
}

// DeleteFileForce forcefully deletes a file (handles readonly, hidden files)
func (oe *OSExecutor) DeleteFileForce(filePath string) error {
	// Validate input
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("File doesn't exist: %s\n", filePath)
		return nil
	}

	// Try to change file permissions first (for readonly files)
	if err := os.Chmod(filePath, 0666); err != nil {
		fmt.Printf("Warning: couldn't change file permissions: %v\n", err)
	}

	// Try Go's built-in function first
	if err := os.Remove(filePath); err == nil {
		fmt.Printf("File deleted successfully: %s\n", filePath)
		return nil
	}

	// Fallback to more forceful OS-specific commands
	var cmdMap OSCommand
	switch oe.osType {
	case consts.OS_WINDOWS:
		// More aggressive Windows deletion
		cmdMap = OSCommand{
			Windows: []string{"cmd", "/C", "attrib", "-R", "-H", "-S", filePath, "&", "del", "/F", "/Q", filePath},
		}
	default:
		cmdMap = OSCommand{
			MacOS:  []string{"rm", "-rf", filePath},
			Linux:  []string{"rm", "-rf", filePath},
			Ubuntu: []string{"rm", "-rf", filePath},
		}
	}

	output, err := oe.ExecuteCommand(cmdMap)
	if err != nil {
		return fmt.Errorf("failed to force delete file '%s': %w\nOutput: %s", filePath, err, string(output))
	}

	fmt.Printf("File force deleted: %s\n", filePath)
	return nil
}

// DeleteFiles deletes multiple files
func (oe *OSExecutor) DeleteFiles(filePaths []string) error {
	var errors []string

	for _, filePath := range filePaths {
		if err := oe.DeleteFile(filePath); err != nil {
			errors = append(errors, fmt.Sprintf("failed to delete %s: %v", filePath, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("some files could not be deleted:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

// FileExists checks if a file exists
func (oe *OSExecutor) FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// GetFileInfo returns file information
func (oe *OSExecutor) GetFileInfo(filePath string) (os.FileInfo, error) {
	return os.Stat(filePath)
}

// create dir to file if not exists
func (oe *OSExecutor) CreateDirIfNotExists(pathFile string) error {
	dirPath := filepath.Dir(pathFile)
	if err := oe.CreateDirectory(dirPath); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	fmt.Printf("Directory created: %s\n", dirPath)
	return nil
}

// GetOSName returns the OS name as string
func (oe *OSExecutor) GetOSName() string {
	switch oe.osType {
	case consts.OS_WINDOWS:
		return "Windows"
	case consts.OS_MACOS:
		return "macOS"
	case consts.OS_LINUX:
		return "Linux"
	case consts.OS_UBUNTU:
		return "Ubuntu"
	default:
		return "Unknown"
	}
}

// ExecuteShellCommand executes a shell command with proper shell for each OS
func (oe *OSExecutor) ExecuteShellCommand(command string) ([]byte, error) {
	var cmd *exec.Cmd

	switch oe.osType {
	case consts.OS_WINDOWS:
		cmd = exec.Command("cmd", "/C", command)
	case consts.OS_MACOS, consts.OS_LINUX, consts.OS_UBUNTU:
		cmd = exec.Command("sh", "-c", command)
	default:
		return nil, fmt.Errorf("unsupported operating system: %d", oe.osType)
	}

	return cmd.Output()
}

// Global executor instance
var GlobalOSExecutor *OSExecutor

// InitOSExecutor initializes the global OS executor
func InitOSExecutor(osType int) {
	GlobalOSExecutor = NewOSExecutor(osType)
}

// InitOSExecutorAuto initializes the global OS executor with auto-detection
func InitOSExecutorAuto() {
	osType := GetCurrentOS()
	GlobalOSExecutor = NewOSExecutor(osType)
}
