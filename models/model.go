package models

//
type SectionVideo struct {
	CourseID    int       `json:"course_id"`
	SectionID   int       `json:"section_id"`
	SectionPath string    `json:"section_path"`
	M3u8Link    string    `json:"m3u8_link"`
	YtbVideoID  string    `json:"ytb_video_id"`
	Chapters    []Chapter `json:"chapters"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

// data error
type DataError struct {
	CourseID    int    `json:"course_id"`
	SectionID   int    `json:"section_id"`
	SectionPath string `json:"section_path"`
	M3u8Link    string `json:"m3u8_link"`
}

// structure for JSON data files
type DataFile struct {
	Data struct {
		ID         int    `json:"id"`
		SkillID    int    `json:"skill_id"`
		Title      string `json:"title"`
		CourseID   int    `json:"course_id"`
		Compulsory bool   `json:"compulsory"`
		UnitTitle  string `json:"unit_title"`
		Content    string `json:"content"`
		Element    struct {
			ID              int           `json:"id"`
			Name            string        `json:"name"`
			Type            int           `json:"type"`
			SlowSound       bool          `json:"slow_sound"`
			DefaultProvider string        `json:"default_provider"`
			Content         string        `json:"content"`
			MuseVideoID     *int          `json:"muse_video_id"`
			YoutubeURL      *string       `json:"youtube_url"`
			MediaPath       string        `json:"media_path"`
			TimestampVideo  []Chapter     `json:"timestamp_video"`
			PlaytimeString  *string       `json:"playtime_string"`
			Photos          []interface{} `json:"photos"`
			VideoSource     struct {
				WebURL    string `json:"web_url"`
				URL       string `json:"url"`
				Extension string `json:"extension"`
			} `json:"video_source"`
			Flashcards []interface{} `json:"flashcards"`
		} `json:"element"`
	} `json:"data"`
}
