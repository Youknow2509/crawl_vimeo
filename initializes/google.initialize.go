package initializes

import (
	"context"
	"log"

	"github.com/youknow2509/crawl_vimeo/consts"
	"github.com/youknow2509/crawl_vimeo/global"
	"github.com/youknow2509/crawl_vimeo/services"
	"github.com/youknow2509/crawl_vimeo/services/impl"
	"google.golang.org/api/option"
)

// InitializeGoogleServices initializes all Google-related services
func initializeGoogleServices() error {
    // Initialize Google service
    if err := initializeGoogleService(); err != nil {
        return err
    }

    // Initialize Google User service
    if err := initializeGoogleUserService(); err != nil {
        return err
    }

    // Initialize YouTube service
    if err := initializeYtbService(); err != nil {
        return err
    }

    return nil
}

// Initialize Google service
func initializeGoogleService() error {
    googleService := impl.NewGoogleService()
    services.InitGoogleService(googleService)
    log.Println("Google service initialized successfully")
    return nil
}

// Initialize Google User service
func initializeGoogleUserService() error {
    googleUserService := impl.NewGoogleUserService()
    services.InitGoogleUserService(googleUserService)
    log.Println("Google User service initialized successfully")
    return nil
}

// InitializeYtbService initializes the YouTube client and sets it in the global variable.
func initializeYtbService() error {
    ctx := context.Background()

    // Get services
    googleService, err := services.GetGoogleService()
    if err != nil {
        return err
    }

    googleUserService, err := services.GetGoogleUserService()
    if err != nil {
        return err
    }

    // Get config
    config, err := googleService.GetConfigServer(consts.YTB_CLIENT_SECRET_PATH)
    if err != nil {
        return err
    }

    // Get or create token
    token, err := googleUserService.GetTokenFromFile(consts.YTB_USER_AUTH_FILE)
    if err != nil {
        log.Printf("Could not load token from file: %v", err)
        token, err = googleUserService.GetTokenFromWeb(ctx, config, consts.YTB_USER_AUTH_FILE)
        if err != nil {
            return err
        }
    }

    // Validate and refresh token if needed
    status, isValid := googleUserService.ValidateToken(ctx, token)
    if !isValid {
        if status == consts.YTB_USER_TOKEN_EXPIRED {
            log.Println("Token expired, attempting to refresh...")
            newToken, err := googleUserService.RefreshToken(ctx, config, token)
            if err != nil {
                log.Printf("Failed to refresh token: %v", err)
                token, err = googleUserService.GetTokenFromWeb(ctx, config, consts.YTB_USER_AUTH_FILE)
                if err != nil {
                    return err
                }
            } else {
                token = newToken
                googleUserService.SaveToken(consts.YTB_USER_AUTH_FILE, token)
            }
        } else {
            log.Printf("Token invalid (status: %d), getting new token...", status)
            token, err = googleUserService.GetTokenFromWeb(ctx, config, consts.YTB_USER_AUTH_FILE)
            if err != nil {
                return err
            }
        }
    }

    // Create YouTube service
    client := config.Client(ctx, token)
    youtubeService, err := googleUserService.GetYoutubeService(ctx, option.WithHTTPClient(client))
    if err != nil {
        return err
    }
	// Set the YouTube service in the global variable
	global.YTB_SERVICE = youtubeService

    // Initialize YouTube service
    ytbService := impl.NewYoutubeService(youtubeService)
    services.InitYoutubeService(ytbService)

	// Set the YouTube service in the global variable
	global.YTB_SERVICE_ACTION = ytbService
    
    log.Println("YouTube service initialized successfully")
    return nil
}