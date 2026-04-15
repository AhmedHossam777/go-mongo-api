package dto

type CreateBlogDto struct {
	Title    string `validate:"required,min=2,max=200"`
	Content  string `validate:"required,min=10"`
	ImageURL string
}

type UpdateBlogDto struct {
	Title   *string `json:"title" validate:"omitempty,min=2,max=200"`
	Content *string `json:"content" validate:"omitempty,min=10"`
}

type AddCommentDto struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}
