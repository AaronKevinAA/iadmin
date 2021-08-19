package request

type SysOperationRecordSearch struct {
	Method    string  `json:"method"` // 请求方法
	Path      string  `json:"path"`   // 请求路径
	Status    string  `json:"status"` // 请求状态
	CreatedAt []int64 `json:"createdAt"`
	Pagination
}

type SysOperationRecordExcelOut struct {
	SysOperationRecordSearch SysOperationRecordSearch `json:"sysOperationRecordSearch"`
	ExcelOutConfig           ExcelOutRequest          `json:"excelOutConfig"`
}
