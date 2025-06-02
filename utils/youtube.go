package utils

import (
	"fmt"

	"github.com/youknow2509/crawl_vimeo/models"
)

// Tạo description có chapters
func CreateDescriptionWithChapters(chapters []models.Chapter) string {
	desc := "Nội dung video:\n"
	for _, ch := range chapters {
		desc += fmt.Sprintf("%s %s\n", ch.Time, ch.Content)
	}
	desc += "\n\n"
	desc += "Cảm ơn các bạn đã xem video! Nếu thấy hay hãy like và đăng ký kênh để ủng hộ mình nhé!\n"
	return desc
}
