package entity

type Video struct {
	Title       string `json:"title" binding:"required" validate:"happy"`
	URL         string `json:"url" binding:"url"`
	Description string `json:"description"`
}
