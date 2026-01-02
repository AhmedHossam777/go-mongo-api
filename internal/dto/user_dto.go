package dto

type CreateUserDto struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	Role     string `json:"role" validate:"omitempty"`
}

type UpdateUserDto struct {
	Name     *string `json:"name" validate:"omitempty,min=2,max=100"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Password *string `json:"password" validate:"omitempty,min=6,max=100"`
}
