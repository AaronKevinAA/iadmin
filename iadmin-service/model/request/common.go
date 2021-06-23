package request


type Pagination struct {
	Current  int `json:"current" form:"page"`      // 页码
	PageSize int `json:"pageSize" form:"pageSize"` // 每页大小
}

type IdsReq struct {
	Ids []int `json:"ids" form:"ids"`
}