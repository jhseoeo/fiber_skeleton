package req

type PaginationReq struct {
	Page  int `json:"page"  query:"page"  validate:"min=1,max=10000"`
	Limit int `json:"limit" query:"limit" validate:"min=1,max=100"`
}
