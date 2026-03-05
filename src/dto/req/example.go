package req

type CreateExampleReq struct {
	Content string `json:"content" validate:"required,max=10000"`
}

type UpdateExampleReq struct {
	Content string `json:"content" validate:"required,max=10000"`
}
