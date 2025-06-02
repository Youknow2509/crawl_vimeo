package exec

func Execute() {
	/**
	 * - Get list id videos m3u8 from file json, chapters in video 
	 * -> {
	 * 		"course_id": "12345",
	 * 		"section_id": "67890",
	 * 		"m3u8_link": "https://example.com/video.m3u8",
	 * 		"ytb_video_id": "",
	 * 		"chapters": [models.Chapter]
	 * 		"title": "Video Title",
	 * 		"description": "Video Description"
	 * }
	 * -> Count number, create file error.json, create file success.json
	 * Loop:
	 * 	- Get video from m3u8 link
	 * 	- Upload video to YouTube
	 *  - Save video id ytb, section id, course id to file json
	 *  - if success, save to file success.json
	 *  - if error, save to file error.json
	 *  - time sleep random 5-10s
	 *
	 * - Show dashboard with number of success, error, total
	 */
}