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

type SectionVideo2 struct {
	CourseID    int       `json:"course_id"`
	SectionID   int       `json:"section_id"`
	SectionPath string    `json:"section_path"`
	M3u8Link    string    `json:"m3u8_link"`
	YtbVideoID  string    `json:"ytb_video_id"`
	Path        string    `json:"path"`
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
		ID        int    `json:"id"`
		Title     string `json:"title"`
		CourseID  int    `json:"course_id"`
		UnitTitle string `json:"unit_title"`
		Content   string `json:"content"`
		Element   struct {
			ID             int       `json:"id"`
			Name           string    `json:"name"`
			TimestampVideo []Chapter `json:"timestamp_video"`
			VideoSource    struct {
				WebUrl    string `json:"web_url"`
				Url       string `json:"url"`
				Extension string `json:"extension"`
			} `json:"video_source"`
		} `json:"element"`
		SubElement struct {
			MediaPath string `json:"media_path"`
		} `json:"sub_element"`
	} `json:"data"`
}
