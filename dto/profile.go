package dto

type ProfileResponse struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
