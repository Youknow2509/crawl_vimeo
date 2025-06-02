# Contact:
- **Mail**: *lytranvinh.work@gmail.com*
- **Github**: *https://github.com/Youknow2509*

# Project: *Crawl Video Vimeo*
# Requirements:
- **Golang**
- **ffmpeg**: 
    - *macOS*: `brew install ffmpeg`.
    - *Linux*: `sudo apt install ffmpeg`.
    - *Windows*: `winget install --id=Gyan.FFmpeg  -e` - Note path to `ffmpeg.exe` in `PATH` environment variable and ...
- **Enable Youtube Api**:
    - Get `client_secret.json` and save to `secret/client_secret.json`.
    - Scope set:
      - "https://www.googleapis.com/auth/youtube"
      - "https://www.googleapis.com/auth/youtube.force-ssl"
      - "https://www.googleapis.com/auth/youtube.readonly"
      - "https://www.googleapis.com/auth/youtube.upload"
    - Set `Redirect URI` to `http://localhost:8080/`
    - If have file token user authentication, `Youtube API Service` save it to `user_auth.json`.