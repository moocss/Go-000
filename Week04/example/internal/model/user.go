package model

type User struct {
	Id int64 `json:"id"`
	UserName string `json:"username"`
	Email string `json:"email"`
}