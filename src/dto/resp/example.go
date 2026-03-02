package resp

type GetExampleResp struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
}

type CreateExampleResp struct {
	ID uint `json:"id"`
}
