package main

import (
    "fmt"
    "io"
    "net/http"
    "regexp"
    "strings"
)

func extractHMACFromHTML(html string) (string, error) {
    // Pattern để tìm signed URL trong playerConfig
    patterns := []string{
        // Pattern 1: Tìm trong avc_url hoặc url
        `"(?:avc_)?url":"[^"]*vod-adaptive-ak\.vimeocdn\.com/(exp=\d+~acl=[^~]+~hmac=[a-f0-9]{64})/`,
        
        // Pattern 2: Tìm trực tiếp exp=...~hmac=...
        `(exp=\d+~acl=[^~]+~hmac=[a-f0-9]{64})`,
        
        // Pattern 3: Tìm trong các URL khác
        `vimeocdn\.com/([^/]*exp=\d+~acl=[^~]+~hmac=[a-f0-9]{64})`,
    }

    for i, pattern := range patterns {
        re := regexp.MustCompile(pattern)
        matches := re.FindStringSubmatch(html)
        if len(matches) > 1 {
            fmt.Printf("✅ Found HMAC with pattern %d\n", i+1)
            return matches[1], nil
        }
    }

    return "", fmt.Errorf("HMAC not found in HTML")
}

func extractAllSignedURLs(html string) ([]string, error) {
    // Tìm tất cả signed URLs
    pattern := `"[^"]*vod-adaptive-ak\.vimeocdn\.com/(exp=\d+~acl=[^~]+~hmac=[a-f0-9]{64})/[^"]*"`
    re := regexp.MustCompile(pattern)
    
    allMatches := re.FindAllStringSubmatch(html, -1)
    var signedURLs []string
    
    for _, match := range allMatches {
        if len(match) > 1 {
            signedURLs = append(signedURLs, match[1])
        }
    }
    
    // Remove duplicates
    seen := make(map[string]bool)
    var unique []string
    for _, url := range signedURLs {
        if !seen[url] {
            seen[url] = true
            unique = append(unique, url)
        }
    }
    
    return unique, nil
}

func parseSignedURL(signedURL string) map[string]string {
    // Parse exp=...~acl=...~hmac=... thành map
    parts := strings.Split(signedURL, "~")
    result := make(map[string]string)
    
    for _, part := range parts {
        if strings.Contains(part, "=") {
            kv := strings.SplitN(part, "=", 2)
            if len(kv) == 2 {
                result[kv[0]] = kv[1]
            }
        }
    }
    
    return result
}

func getVimeoHMAC(videoID string) (map[string]string, error) {
    url := fmt.Sprintf("https://player.vimeo.com/video/%s", videoID)
    
    headers := map[string]string{
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        "Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
    }
    
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    for key, value := range headers {
        req.Header.Set(key, value)
    }
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    html := string(body)
    
    // Extract HMAC
    signedURL, err := extractHMACFromHTML(html)
    if err != nil {
        return nil, err
    }

	print("✅ Extracted signed URL: ", signedURL)
    
    // Parse thành map
    params := parseSignedURL(signedURL)
    
    return params, nil
}

func main() {
    videoID := "713179128"
    
    fmt.Println("=== Extracting HMAC from Vimeo Player ===")
    
    // Method 1: Get from live page
    params, err := getVimeoHMAC(videoID)
    if err != nil {
        fmt.Printf("❌ Error: %v\n", err)
        return
    }
    
    fmt.Println("\n📋 Extracted Parameters:")
    for key, value := range params {
        fmt.Printf("  %s: %s\n", key, value)
    }
    
}