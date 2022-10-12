package models

type UserSelectResult struct {
	dbResponse
	Result []userSelectResultData `json:"result"`
}

type userSelectResultData struct {
	Id string `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserInsertResult struct {
	dbResponse
	Result []userInsertResultData `json:"result"`
}

type userInsertResultData struct {
	Id string `json:"id"`
}