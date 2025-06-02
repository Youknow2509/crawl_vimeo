package utils

import (
	"context"
	"fmt"
	"net/http"
	"golang.org/x/oauth2"
)

// Helper method to validate token with Google API
func ValidateTokenWithGoogle(ctx context.Context, token *oauth2.Token) error {
    // Create a simple HTTP client with the token
    client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
    
    // Make a simple API call to verify the token
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        return fmt.Errorf("failed to validate token: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("token validation failed with status: %d", resp.StatusCode)
    }
    
    return nil
}