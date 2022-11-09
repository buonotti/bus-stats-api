package v1

type RefreshRequest struct {
	Token string `json:"token" binding:"required,jwt"`
	Id    string `json:"id" binding:"required"`
}

type RefreshResponse struct {
	Result string `json:"result"`
	Token  string `json:"token"`
}
