package impl

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"

    "github.com/youknow2509/crawl_vimeo/services"
)

// struct
type WindowsOsExecutorService struct{}

// ########################################################################################

// ConvertM3U8ToMP4 implements services.IOSExecutor.
func (w *WindowsOsExecutorService) ConvertM3U8ToMP4(inputUrl string, outputPath string) error {
    // Create directory if not exists
    if err := w.CreateDirIfNotExists(outputPath); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }

    // Check if FFmpeg is available
    if _, err := exec.LookPath("ffmpeg"); err != nil {
        return fmt.Errorf("FFmpeg is not available on Windows: %w", err)
    }

    // FFmpeg command for Windows
    cmd := exec.Command("ffmpeg", "-y", "-i", inputUrl, "-c", "copy", outputPath)
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("FFmpeg conversion failed on Windows: %w\nOutput: %s", err, string(output))
    }

    fmt.Printf("Windows: Conversion completed: %s -> %s\n", inputUrl, outputPath)
    return nil
}

// CreateDirIfNotExists implements services.IOSExecutor.
func (w *WindowsOsExecutorService) CreateDirIfNotExists(path string) error {
    dirPath := filepath.Dir(path)
    return w.CreateDirectory(dirPath)
}

// CreateDirectory implements services.IOSExecutor.
func (w *WindowsOsExecutorService) CreateDirectory(path string) error {
    // Try Go's built-in function first
    if err := os.MkdirAll(path, 0755); err == nil {
        return nil
    }

    // Fallback to Windows command
    // Use mkdir with quotes to handle spaces in path
    cmd := exec.Command("cmd", "/C", "mkdir", fmt.Sprintf(`"%s"`, path))
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("Windows mkdir failed: %w\nOutput: %s", err, string(output))
    }

    return nil
}

// DeleteFile implements services.IOSExecutor.
func (w *WindowsOsExecutorService) DeleteFile(path string) error {
    if path == "" {
        return fmt.Errorf("file path cannot be empty")
    }

    // Check if file exists
    if !w.FileExists(path) {
        fmt.Printf("Windows: File doesn't exist: %s\n", path)
        return nil
    }

    // Try Go's built-in function first
    if err := os.Remove(path); err == nil {
        fmt.Printf("Windows: File deleted successfully: %s\n", path)
        return nil
    }

    // Fallback to Windows command
    cmd := exec.Command("cmd", "/C", "del", "/F", "/Q", fmt.Sprintf(`"%s"`, path))
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("Windows del command failed: %w\nOutput: %s", err, string(output))
    }

    fmt.Printf("Windows: File deleted via command: %s\n", path)
    return nil
}

// DeleteFileForce implements services.IOSExecutor.
func (w *WindowsOsExecutorService) DeleteFileForce(path string) error {
    if path == "" {
        return fmt.Errorf("file path cannot be empty")
    }

    if !w.FileExists(path) {
        fmt.Printf("Windows: File doesn't exist: %s\n", path)
        return nil
    }

    // Remove file attributes (readonly, hidden, system)
    cmd1 := exec.Command("cmd", "/C", "attrib", "-R", "-H", "-S", fmt.Sprintf(`"%s"`, path))
    cmd1.CombinedOutput() // Ignore errors, just try

    // Try normal delete first
    if err := w.DeleteFile(path); err == nil {
        return nil
    }

    // Force delete with more aggressive approach
    cmd := exec.Command("cmd", "/C", "del", "/F", "/A", "/Q", fmt.Sprintf(`"%s"`, path))
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("Windows force delete failed: %w\nOutput: %s", err, string(output))
    }

    fmt.Printf("Windows: File force deleted: %s\n", path)
    return nil
}

// DeleteFiles implements services.IOSExecutor.
func (w *WindowsOsExecutorService) DeleteFiles(paths []string) error {
    var errors []string

    for _, path := range paths {
        if err := w.DeleteFile(path); err != nil {
            errors = append(errors, fmt.Sprintf("failed to delete %s: %v", path, err))
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("Windows: some files could not be deleted:\n%s", strings.Join(errors, "\n"))
    }

    return nil
}

// ExecuteCommand implements services.IOSExecutor.
func (w *WindowsOsExecutorService) ExecuteCommand(command string, args ...string) ([]byte, error) {
    cmd := exec.Command(command, args...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return output, fmt.Errorf("Windows command failed: %w\nCommand: %s %v\nOutput: %s", 
            err, command, args, string(output))
    }
    return output, nil
}

// ExecuteShellCommand implements services.IOSExecutor.
func (w *WindowsOsExecutorService) ExecuteShellCommand(command string, args ...string) ([]byte, error) {
    // Combine command and args into a single command string
    fullCommand := command
    if len(args) > 0 {
        fullCommand += " " + strings.Join(args, " ")
    }

    cmd := exec.Command("cmd", "/C", fullCommand)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return output, fmt.Errorf("Windows shell command failed: %w\nCommand: %s\nOutput: %s", 
            err, fullCommand, string(output))
    }
    return output, nil
}

// FileExists implements services.IOSExecutor.
func (w *WindowsOsExecutorService) FileExists(path string) bool {
    _, err := os.Stat(path)
    return !os.IsNotExist(err)
}

// GetFileInfo implements services.IOSExecutor.
func (w *WindowsOsExecutorService) GetFileInfo(path string) (os.FileInfo, error) {
    return os.Stat(path)
}

// RemoveDirectory implements services.IOSExecutor.
func (w *WindowsOsExecutorService) RemoveDirectory(path string) error {
    // Try Go's built-in function first
    if err := os.RemoveAll(path); err == nil {
        return nil
    }

    // Fallback to Windows command
    cmd := exec.Command("cmd", "/C", "rmdir", "/S", "/Q", fmt.Sprintf(`"%s"`, path))
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("Windows rmdir failed: %w\nOutput: %s", err, string(output))
    }

    fmt.Printf("Windows: Directory removed: %s\n", path)
    return nil
}

// ########################################################################################

// new service implementation IOSExecutor
func NewWindowsOsExecutorService() services.IOSExecutor {
    return &WindowsOsExecutorService{}
}