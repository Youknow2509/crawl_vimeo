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
type MacOsOsExecutorService struct{}

// ########################################################################################

// ConvertM3U8ToMP4 implements services.IOSExecutor.
func (m *MacOsOsExecutorService) ConvertM3U8ToMP4(inputUrl string, outputPath string) error {
    // Create directory if not exists
    if err := m.CreateDirIfNotExists(outputPath); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }

    // Check if FFmpeg is available
    if _, err := exec.LookPath("ffmpeg"); err != nil {
        return fmt.Errorf("FFmpeg is not available on macOS: %w", err)
    }

    // FFmpeg command for macOS
    cmd := exec.Command("ffmpeg", "-y", "-i", inputUrl, "-c", "copy", outputPath)
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("FFmpeg conversion failed on macOS: %w\nOutput: %s", err, string(output))
    }

    fmt.Printf("macOS: Conversion completed: %s -> %s\n", inputUrl, outputPath)
    return nil
}

// CreateDirIfNotExists implements services.IOSExecutor.
func (m *MacOsOsExecutorService) CreateDirIfNotExists(path string) error {
    dirPath := filepath.Dir(path)
    return m.CreateDirectory(dirPath)
}

// CreateDirectory implements services.IOSExecutor.
func (m *MacOsOsExecutorService) CreateDirectory(path string) error {
    // Try Go's built-in function first
    if err := os.MkdirAll(path, 0755); err == nil {
        return nil
    }

    // Fallback to macOS command
    cmd := exec.Command("mkdir", "-p", path)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("macOS mkdir failed: %w\nOutput: %s", err, string(output))
    }

    return nil
}

// DeleteFile implements services.IOSExecutor.
func (m *MacOsOsExecutorService) DeleteFile(path string) error {
    if path == "" {
        return fmt.Errorf("file path cannot be empty")
    }

    // Check if file exists
    if !m.FileExists(path) {
        fmt.Printf("macOS: File doesn't exist: %s\n", path)
        return nil
    }

    // Try Go's built-in function first
    if err := os.Remove(path); err == nil {
        fmt.Printf("macOS: File deleted successfully: %s\n", path)
        return nil
    }

    // Fallback to macOS command
    cmd := exec.Command("rm", "-f", path)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("macOS rm command failed: %w\nOutput: %s", err, string(output))
    }

    fmt.Printf("macOS: File deleted via command: %s\n", path)
    return nil
}

// DeleteFileForce implements services.IOSExecutor.
func (m *MacOsOsExecutorService) DeleteFileForce(path string) error {
    if path == "" {
        return fmt.Errorf("file path cannot be empty")
    }

    if !m.FileExists(path) {
        fmt.Printf("macOS: File doesn't exist: %s\n", path)
        return nil
    }

    // Change permissions first (for protected files)
    os.Chmod(path, 0666)

    // Try normal delete first
    if err := m.DeleteFile(path); err == nil {
        return nil
    }

    // Force delete with more aggressive approach
    cmd := exec.Command("rm", "-rf", path)
    output, err := cmd.CombinedOutput()
    if err != nil {
        // Try with sudo as last resort (will ask for password)
        sudoCmd := exec.Command("sudo", "rm", "-rf", path)
        sudoOutput, sudoErr := sudoCmd.CombinedOutput()
        if sudoErr != nil {
            return fmt.Errorf("macOS force delete failed: %w\nOutput: %s\nSudo Output: %s", 
                err, string(output), string(sudoOutput))
        }
    }

    fmt.Printf("macOS: File force deleted: %s\n", path)
    return nil
}

// DeleteFiles implements services.IOSExecutor.
func (m *MacOsOsExecutorService) DeleteFiles(paths []string) error {
    var errors []string

    for _, path := range paths {
        if err := m.DeleteFile(path); err != nil {
            errors = append(errors, fmt.Sprintf("failed to delete %s: %v", path, err))
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("macOS: some files could not be deleted:\n%s", strings.Join(errors, "\n"))
    }

    return nil
}

// ExecuteCommand implements services.IOSExecutor.
func (m *MacOsOsExecutorService) ExecuteCommand(command string, args ...string) ([]byte, error) {
    cmd := exec.Command(command, args...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return output, fmt.Errorf("macOS command failed: %w\nCommand: %s %v\nOutput: %s", 
            err, command, args, string(output))
    }
    return output, nil
}

// ExecuteShellCommand implements services.IOSExecutor.
func (m *MacOsOsExecutorService) ExecuteShellCommand(command string, args ...string) ([]byte, error) {
    // Combine command and args into a single command string
    fullCommand := command
    if len(args) > 0 {
        fullCommand += " " + strings.Join(args, " ")
    }

    cmd := exec.Command("sh", "-c", fullCommand)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return output, fmt.Errorf("macOS shell command failed: %w\nCommand: %s\nOutput: %s", 
            err, fullCommand, string(output))
    }
    return output, nil
}

// FileExists implements services.IOSExecutor.
func (m *MacOsOsExecutorService) FileExists(path string) bool {
    _, err := os.Stat(path)
    return !os.IsNotExist(err)
}

// GetFileInfo implements services.IOSExecutor.
func (m *MacOsOsExecutorService) GetFileInfo(path string) (os.FileInfo, error) {
    return os.Stat(path)
}

// RemoveDirectory implements services.IOSExecutor.
func (m *MacOsOsExecutorService) RemoveDirectory(path string) error {
    // Try Go's built-in function first
    if err := os.RemoveAll(path); err == nil {
        return nil
    }

    // Fallback to macOS command
    cmd := exec.Command("rm", "-rf", path)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("macOS rm -rf failed: %w\nOutput: %s", err, string(output))
    }

    fmt.Printf("macOS: Directory removed: %s\n", path)
    return nil
}

// ########################################################################################

// new service implementation IOSExecutor
func NewMacOsOsExecutorService() services.IOSExecutor {
    return &MacOsOsExecutorService{}
}