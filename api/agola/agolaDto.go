package agola

type AgolaCreateTokenDto struct {
	Token string `json:"token"`
}

type AgolaCreateORGDto struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Visibility string `json:"visibility"`
}
