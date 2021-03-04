package model

type User struct {
	Email              string   `json:"email"`
	AgolaUsersRef      []string `json:"agolaUsersRef"`
	AgolaUserToken     string   `json:"agolaUserToken"`
	AgolauserTokenName string   `json:"agolauserTokenName"`
	//Role  string `json:"role"`
}
