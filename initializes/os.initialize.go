package initializes

import (
	"fmt"
	"log"

	"github.com/youknow2509/crawl_vimeo/global"
	"github.com/youknow2509/crawl_vimeo/services/create"
	"github.com/youknow2509/crawl_vimeo/utils"
)

// initializeOsSystem initializes the OS system and sets it in global
func initializeOsSystem() {
    // Auto-detect OS type
    osType := utils.GetCurrentOS()
    
    // Set OS type in global
    global.OS_SYSTEM = osType
    
    // Initialize global OS executor and set global OS exec system
	global.OS_EXECUTOR = create.FactoryCreateOsExce(osType)

    // Log the detected OS
    log.Printf("OS System initialized: %s (Type: %d)", utils.GetOSName(osType), osType)
}

// ValidateOSInitialization validates that OS system is properly initialized
func ValidateOSInitialization() error {
    if global.OS_EXECUTOR == nil {
        return fmt.Errorf("global OS executor is not initialized")
    }
    
    if global.OS_SYSTEM < 0 {
        return fmt.Errorf("global OS_SYSTEM is not properly set")
    }
    
    log.Printf("OS System validation passed: %s", utils.GetOSName(global.OS_SYSTEM))
    return nil
}