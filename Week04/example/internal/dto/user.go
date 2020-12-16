package dto

type UserRequest struct {
	UserName string `json:"username"`
	Email string `json:"email"`
}

type UserResponse struct {
	UserName string `json:"username"`
	Email string `json:"email"`
}