package health

type DbInfoResponse struct {
	Time   string       `json:"time"`
	Status string       `json:"status"`
	Result DbInfoResult `json:"result"`
}

type DbInfoResult struct {
	Dl any `json:"dl"`
	Dt any `json:"dt"`
	Sc any `json:"sc"`
	Tb any `json:"tb"`
}
