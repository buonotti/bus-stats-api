package v1

type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginModel struct {
	Username string `json:"username"`
	Token string `json:"token"`
	Image string `json:"image"`
}