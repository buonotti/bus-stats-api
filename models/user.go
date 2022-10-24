package models

type BaseUser struct {
	Id    string `json:"id"`
	Email string `json:"email" `
}

type UserId string

type UserSelectResult struct {
	dbResponse
	Result []userSelectResultData `json:"result"`
}

type userSelectResultData struct {
	Id       string    `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Image    userImage `json:"image"`
}

type userImage struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type UserInsertResult struct {
	dbResponse
	Result []userInsertResultData `json:"result"`
}

type userInsertResultData struct {
	Id string `json:"id"`
}
