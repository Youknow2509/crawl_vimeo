package initializes

import (
	"context"
	"fmt"
	"time"

	"github.com/youknow2509/crawl_vimeo/global"
	"github.com/youknow2509/crawl_vimeo/services"
)

// InitializeYtb initializes the YouTube client and sets it in the global variable.
func initializeYtb() {
	ctx := context.Background()
	ytService, token, err := services.CreateYtbClient(ctx)
	if err != nil {
		panic("Không thể tạo client YouTube: " + err.Error())
	}
	fmt.Println("AccessToken:", token.AccessToken)
	fmt.Println("Expiry:", token.Expiry.Format(time.RFC3339))
	
	// save to global
	global.YTB_CLIENT = ytService
	// 
	fmt.Println("YouTube client initialized successfully")
}

// Initialize user auth client youtobe
func initializeYtbUserAuth() {
	
}