package services

import "os"

// interface for executing OS commands
type (
	IOSExecutor interface {
		ExecuteCommand(command string, args ...string) ([]byte, error)
		// directory management
		CreateDirIfNotExists(path string) error
		CreateDirectory(path string) error
		RemoveDirectory(path string) error
		// ...

		// file management
		DeleteFile(path string) error
		DeleteFileForce(path string) error
		DeleteFiles(paths []string) error
		FileExists(path string) bool
		GetFileInfo(path string) (os.FileInfo, error)
		// ...
		
		// FFmpeg commands
		ConvertM3U8ToMP4(inputUrl, outputPath string) error
		// ...

		// other commands
		ExecuteShellCommand(command string, args ...string) ([]byte, error)
	}	
)

// variable to hold the OS executor instance
var (
	vIOSExecutorInstance IOSExecutor
)

// initializes the OS executor instance
func InitOSExecutorService(executor IOSExecutor) {
	vIOSExecutorInstance = executor
}

// get instance of the OS executor
func GetOSExecutorService() IOSExecutor {
	if vIOSExecutorInstance == nil {
		panic("OS Executor service is not initialized")
	}
	return vIOSExecutorInstance
}


