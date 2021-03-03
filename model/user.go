package model

type User struct {
	Email         string   `json:"email"`
	AgolaUsersRef []string `json:"agolaUsersRef"`
	//Role  string `json:"role"`
}
