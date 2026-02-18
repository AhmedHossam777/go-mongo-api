package dto

type CreateProduct struct {
	Title       string   `json:"title" validate:"required,min=2,max=200"`
	Price       float64  `json:"price" validate:"required,gt=0"`
	Description string   `json:"description" validate:"required,min=10,max=2000"`
	Category    string   `json:"category" validate:"required,min=2,max=50"`
	Brand       string   `json:"brand" validate:"required,min=2,max=100"`
	Stock       int      `json:"stock" validate:"required,gte=0"`
	Images      []string `json:"images" validate:"omitempty,dive,url"`
	Tags        []string `json:"tags" validate:"omitempty,dive,min=2,max=30"`
	IsActive    *bool    `json:"is_active" validate:"omitempty"`
}

type UpdateProduct struct {
	Title       string   `json:"title" validate:"omitempty,min=2,max=200"`
	Price       float64  `json:"price" validate:"omitempty,gt=0"`
	Description string   `json:"description" validate:"omitempty,min=10,max=2000"`
	Category    string   `json:"category" validate:"omitempty,min=2,max=50"`
	Brand       string   `json:"brand" validate:"omitempty,min=2,max=100"`
	Stock       *int     `json:"stock" validate:"omitempty,gte=0"`
	Images      []string `json:"images" validate:"omitempty,dive,url"`
	Tags        []string `json:"tags" validate:"omitempty,dive,min=2,max=30"`
	IsActive    *bool    `json:"is_active" validate:"omitempty"`
}
