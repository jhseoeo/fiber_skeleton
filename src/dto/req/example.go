package req

type CreateExampleReq struct {
	Content string `json:"content" validate:"required"`
}

type UpdateExampleReq struct {
	Content string `json:"content" validate:"required"`
}
