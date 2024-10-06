package types

type Student struct {
	Id    int    `json:"id" validate:"required"`
	Name  string `json:"fullName" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}
