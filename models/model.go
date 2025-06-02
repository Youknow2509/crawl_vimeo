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
			Content         string        `json:"content"`
			YoutubeURL      *string       `json:"youtube_url"`
			MediaPath       string        `json:"media_path"`
			TimestampVideo  []Chapter     `json:"timestamp_video"`
		} `json:"element"`
	} `json:"data"`
}
