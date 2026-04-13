package dto

type CreateBlogDto struct {
	Title   string `json:"title" validate:"required,min=2,max=200"`
	Content string `json:"content" validate:"required,min=10"`
}

type UpdateBlogDto struct {
	Title   *string `json:"title" validate:"omitempty,min=2,max=200"`
	Content *string `json:"content" validate:"omitempty,min=10"`
}
