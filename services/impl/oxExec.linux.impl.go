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
type LinuxOsExecutorService struct{}

// ########################################################################################

// ConvertM3U8ToMP4 implements services.IOSExecutor.
func (l *LinuxOsExecutorService) ConvertM3U8ToMP4(inputUrl string, outputPath string) error {
    // Create directory if not exists
    if err := l.CreateDirIfNotExists(outputPath); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }

    // Check if FFmpeg is available
    if _, err := exec.LookPath("ffmpeg"); err != nil {
        return fmt.Errorf("FFmpeg is not available on Linux: %w", err)
    }

    // FFmpeg command for Linux
    cmd := exec.Command("ffmpeg", "-y", "-i", inputUrl, "-c", "copy", outputPath)
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("FFmpeg conversion failed on Linux: %w\nOutput: %s", err, string(output))
    }

    fmt.Printf("Linux: Conversion completed: %s -> %s\n", inputUrl, outputPath)
    return nil
}

// CreateDirIfNotExists implements services.IOSExecutor.
func (l *LinuxOsExecutorService) CreateDirIfNotExists(path string) error {
    dirPath := filepath.Dir(path)
    return l.CreateDirectory(dirPath)
}

// CreateDirectory implements services.IOSExecutor.
func (l *LinuxOsExecutorService) CreateDirectory(path string) error {
    // Try Go's built-in function first
    if err := os.MkdirAll(path, 0755); err == nil {
        return nil
    }

    // Fallback to Linux command
    cmd := exec.Command("mkdir", "-p", path)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("Linux mkdir failed: %w\nOutput: %s", err, string(output))
    }

    return nil
}

// DeleteFile implements services.IOSExecutor.
func (l *LinuxOsExecutorService) DeleteFile(path string) error {
    if path == "" {
        return fmt.Errorf("file path cannot be empty")
    }

    // Check if file exists
    if !l.FileExists(path) {
        fmt.Printf("Linux: File doesn't exist: %s\n", path)
        return nil
    }

    // Try Go's built-in function first
    if err := os.Remove(path); err == nil {
        fmt.Printf("Linux: File deleted successfully: %s\n", path)
        return nil
    }

    // Fallback to Linux command
    cmd := exec.Command("rm", "-f", path)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("Linux rm command failed: %w\nOutput: %s", err, string(output))
    }

    fmt.Printf("Linux: File deleted via command: %s\n", path)
    return nil
}

// DeleteFileForce implements services.IOSExecutor.
func (l *LinuxOsExecutorService) DeleteFileForce(path string) error {
    if path == "" {
        return fmt.Errorf("file path cannot be empty")
    }

    if !l.FileExists(path) {
        fmt.Printf("Linux: File doesn't exist: %s\n", path)
        return nil
    }

    // Change permissions first (for protected files)
    os.Chmod(path, 0666)

    // Try normal delete first
    if err := l.DeleteFile(path); err == nil {
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
            return fmt.Errorf("Linux force delete failed: %w\nOutput: %s\nSudo Output: %s", 
                err, string(output), string(sudoOutput))
        }
    }

    fmt.Printf("Linux: File force deleted: %s\n", path)
    return nil
}

// DeleteFiles implements services.IOSExecutor.
func (l *LinuxOsExecutorService) DeleteFiles(paths []string) error {
    var errors []string

    for _, path := range paths {
        if err := l.DeleteFile(path); err != nil {
            errors = append(errors, fmt.Sprintf("failed to delete %s: %v", path, err))
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("Linux: some files could not be deleted:\n%s", strings.Join(errors, "\n"))
    }

    return nil
}

// ExecuteCommand implements services.IOSExecutor.
func (l *LinuxOsExecutorService) ExecuteCommand(command string, args ...string) ([]byte, error) {
    cmd := exec.Command(command, args...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return output, fmt.Errorf("Linux command failed: %w\nCommand: %s %v\nOutput: %s", 
            err, command, args, string(output))
    }
    return output, nil
}

// ExecuteShellCommand implements services.IOSExecutor.
func (l *LinuxOsExecutorService) ExecuteShellCommand(command string, args ...string) ([]byte, error) {
    // Combine command and args into a single command string
    fullCommand := command
    if len(args) > 0 {
        fullCommand += " " + strings.Join(args, " ")
    }

    cmd := exec.Command("bash", "-c", fullCommand)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return output, fmt.Errorf("Linux shell command failed: %w\nCommand: %s\nOutput: %s", 
            err, fullCommand, string(output))
    }
    return output, nil
}

// FileExists implements services.IOSExecutor.
func (l *LinuxOsExecutorService) FileExists(path string) bool {
    _, err := os.Stat(path)
    return !os.IsNotExist(err)
}

// GetFileInfo implements services.IOSExecutor.
func (l *LinuxOsExecutorService) GetFileInfo(path string) (os.FileInfo, error) {
    return os.Stat(path)
}

// RemoveDirectory implements services.IOSExecutor.
func (l *LinuxOsExecutorService) RemoveDirectory(path string) error {
    // Try Go's built-in function first
    if err := os.RemoveAll(path); err == nil {
        return nil
    }

    // Fallback to Linux command
    cmd := exec.Command("rm", "-rf", path)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("Linux rm -rf failed: %w\nOutput: %s", err, string(output))
    }

    fmt.Printf("Linux: Directory removed: %s\n", path)
    return nil
}

// ########################################################################################
// Linux-specific methods

// InstallPackage installs a package using the appropriate package manager
func (l *LinuxOsExecutorService) InstallPackage(packageName string) error {
    // Try different package managers
    packageManagers := [][]string{
        {"apt-get", "install", "-y", packageName},       // Debian/Ubuntu
        {"yum", "install", "-y", packageName},           // RHEL/CentOS (old)
        {"dnf", "install", "-y", packageName},           // RHEL/CentOS/Fedora (new)
        {"pacman", "-S", "--noconfirm", packageName},    // Arch Linux
        {"zypper", "install", "-y", packageName},        // openSUSE
        {"apk", "add", packageName},                     // Alpine Linux
    }

    for _, manager := range packageManagers {
        if _, err := exec.LookPath(manager[0]); err == nil {
            cmd := exec.Command("sudo", append(manager, packageName)...)
            output, err := cmd.CombinedOutput()
            if err != nil {
                fmt.Printf("Failed to install %s with %s: %v\nOutput: %s\n", 
                    packageName, manager[0], err, string(output))
                continue
            }
            fmt.Printf("Linux: Package %s installed successfully with %s\n", packageName, manager[0])
            return nil
        }
    }

    return fmt.Errorf("no supported package manager found to install %s", packageName)
}

// UpdatePackageList updates the package list
func (l *LinuxOsExecutorService) UpdatePackageList() error {
    updateCommands := [][]string{
        {"apt-get", "update"},                    // Debian/Ubuntu
        {"yum", "check-update"},                  // RHEL/CentOS (old)
        {"dnf", "check-update"},                  // RHEL/CentOS/Fedora (new)
        {"pacman", "-Sy"},                        // Arch Linux
        {"zypper", "refresh"},                    // openSUSE
        {"apk", "update"},                        // Alpine Linux
    }

    for _, updateCmd := range updateCommands {
        if _, err := exec.LookPath(updateCmd[0]); err == nil {
            cmd := exec.Command("sudo", updateCmd...)
            output, err := cmd.CombinedOutput()
            if err != nil {
                fmt.Printf("Failed to update package list with %s: %v\nOutput: %s\n", 
                    updateCmd[0], err, string(output))
                continue
            }
            fmt.Printf("Linux: Package list updated successfully with %s\n", updateCmd[0])
            return nil
        }
    }

    return fmt.Errorf("no supported package manager found for updating package list")
}

// InstallFFmpeg installs FFmpeg on Linux
func (l *LinuxOsExecutorService) InstallFFmpeg() error {
    // Update package list first
    l.UpdatePackageList()

    // Try to install FFmpeg
    return l.InstallPackage("ffmpeg")
}

// GetLinuxDistribution detects the Linux distribution
func (l *LinuxOsExecutorService) GetLinuxDistribution() (string, error) {
    // Try /etc/os-release first (modern standard)
    if l.FileExists("/etc/os-release") {
        cmd := exec.Command("grep", "^ID=", "/etc/os-release")
        output, err := cmd.CombinedOutput()
        if err == nil {
            distro := strings.TrimSpace(string(output))
            distro = strings.TrimPrefix(distro, "ID=")
            distro = strings.Trim(distro, "\"")
            return distro, nil
        }
    }

    // Fallback to other methods
    distroFiles := map[string]string{
        "/etc/debian_version": "debian",
        "/etc/redhat-release": "rhel",
        "/etc/arch-release":   "arch",
        "/etc/SuSE-release":   "suse",
        "/etc/alpine-release": "alpine",
    }

    for file, distro := range distroFiles {
        if l.FileExists(file) {
            return distro, nil
        }
    }

    return "unknown", fmt.Errorf("could not determine Linux distribution")
}

// SetFilePermissions sets file permissions (Linux-specific)
func (l *LinuxOsExecutorService) SetFilePermissions(path string, mode os.FileMode) error {
    return os.Chmod(path, mode)
}

// ChangeFileOwnership changes file ownership (Linux-specific)
func (l *LinuxOsExecutorService) ChangeFileOwnership(path, owner string) error {
    cmd := exec.Command("sudo", "chown", owner, path)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to change ownership: %w\nOutput: %s", err, string(output))
    }
    return nil
}

// ########################################################################################

// new service implementation IOSExecutor
func NewLinuxOsExecutorService() services.IOSExecutor {
    return &LinuxOsExecutorService{}
}