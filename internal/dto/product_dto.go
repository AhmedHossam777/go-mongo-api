package dto

type CreateProduct struct {
	Title       string `json:"title" validate:"required,min=2,max=100"`
	Price       int    `json:"price" validate:"required,gt=0"`
	Description string `json:"description" validate:"required,min=5"`
}

type UpdateProduct struct {
	Title       string `json:"title" validate:"omitempty,min=2,max=100"`
	Price       int    `json:"price" validate:"omitempty,gt=0"`
	Description string `json:"description" validate:"omitempty,min=5"`
}
