package model

import "ginProject/global"

type SysApi struct {
	global.GvaModel
	Path        string `json:"path"`
	Description string `json:"description"`
	ApiGroup    string `json:"apiGroup" gorm:"column:apiGroup;"`
	Method      string `json:"method"`
}

func SysApiExcelOutTableHeadName() (tableHeadName []string) {
	tableHeadName = []string{"ID", "接口描述", "请求路径", "请求方法", "接口分组", "创建时间", "最后更新时间"}
	return tableHeadName
}

func SysApiExcelOutTableData(apiList []SysApi) (tableData [][]string) {
	for _, api := range apiList {
		menuInfo := []string{global.Uint2String(api.ID), api.Description, api.Path, api.Method, api.ApiGroup, global.Timestamp2DateTime(api.CreatedAt), global.Timestamp2DateTime(api.UpdatedAt)}
		tableData = append(tableData, menuInfo)
	}
	return tableData
}

// SysApiExcelInTableHeadName 批量导入模板的表头
func SysApiExcelInTableHeadName() (tableHeadName []string) {
	tableHeadName = []string{"接口描述", "请求路径", "请求方法", "接口分组"}
	return tableHeadName
}
