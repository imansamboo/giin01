package entity

type Video struct {
	Title       string `json:"title" binding:"required,happy" validate:"happy"`
	URL         string `json:"url" binding:"url"`
	Description string `json:"description"`
	Author      Person `json:"author" binding:"required"`
}

type Person struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
	Age  int    `json:"age" binding:"required,min=18,max=100"`
}
