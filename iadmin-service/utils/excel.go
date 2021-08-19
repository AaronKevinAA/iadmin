package utils

import (
	"errors"
	"ginProject/global"
	"ginProject/model"
	"github.com/xuri/excelize/v2"
	"strconv"
)

// ExcelOut 批量导出方法
func ExcelOut(hasTableName bool, tableHeadName []string, tableData [][]string) (err error, filePath string) {
	// 生成一个excel文件
	excelFile := excelize.NewFile()
	if hasTableName {
		// 按行赋值 从A1单元格开始赋值
		errSet := excelFile.SetSheetRow("Sheet1", "A1", &tableHeadName)
		if errSet != nil {
			return errSet, ""
		}
	}
	for index, userData := range tableData {
		// 循环写入excel表格
		var rowIndex int
		if hasTableName {
			rowIndex = index + 2
		} else {
			rowIndex = index + 1
		}
		aixs := "A" + strconv.Itoa(rowIndex)
		errSet := excelFile.SetSheetRow("Sheet1", aixs, &userData)
		if errSet != nil {
			return errSet, ""
		}
	}
	filePathDir := global.GvaConfig.Excel.ExcelStoreDir
	filePath = filePathDir + "/excelOutFileTemp.xlsx"
	err = excelFile.SaveAs(filePath)
	if err != nil {
		return err, ""
	}
	return nil, filePath
}

// GetExcelDataList 获得批量导入Excel表中的数据列表
func GetExcelDataList(excelFilePath string, readTableHead bool) (err error, dataList [][]string) {
	excelFile, errOpen := excelize.OpenFile(excelFilePath)
	if errOpen != nil {
		return errOpen, nil
	}

	rows, errSheet := excelFile.GetRows("Sheet1")
	cols, _ := excelFile.GetCols("Sheet1")
	colLength := len(cols)
	if errSheet != nil {
		return errSheet, nil
	}
	for rowIndex, row := range rows {
		// 若设置了不读第一行则跳过
		if rowIndex == 0 && !readTableHead {
			continue
		}
		var dataRow []string
		for _, colCell := range row {
			dataRow = append(dataRow, colCell)
		}
		// 为了解决一个Bug
		if len(dataRow) < colLength {
			needPlusCount := colLength - len(dataRow)
			for i := 0; i < needPlusCount; i++ {
				dataRow = append(dataRow, "")
			}
		}
		dataList = append(dataList, dataRow)
	}
	return err, dataList
}

// GenerateExcelInTemplate 生成批量导入的模板
func GenerateExcelInTemplate(databaseName string) (filePath string, err error) {
	// 新建一个Excel
	excelFile := excelize.NewFile()
	// 获得表头信息
	var tableHeadName []string
	if databaseName == "sys_user" {
		tableHeadName = model.SysUserExcelInTableHeadName()
	} else if databaseName == "sys_api" {
		tableHeadName = model.SysApiExcelInTableHeadName()
	} else {
		return "", errors.New("数据库名称错误!")
	}
	// 按行赋值 从A1单元格开始赋值
	err = excelFile.SetSheetRow("Sheet1", "A1", &tableHeadName)
	if err != nil {
		return "", err
	}
	// 保存Excel
	filePathDir := global.GvaConfig.Excel.ExcelStoreDir
	filePath = filePathDir + "/excelInTemplateTemp.xlsx"
	err = excelFile.SaveAs(filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}
