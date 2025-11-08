package models

type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name" validate:"required,min=2,max=100"`
	Bio  string `json:"bio,omitempty"`
}
