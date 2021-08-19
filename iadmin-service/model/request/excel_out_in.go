package request

type ExcelOutRequest struct {
	HasTableHead bool `json:"hasTableHead"`
	HasAllData   bool `json:"hasAllData"`
}

type ExcelInRequest struct {
	SaveFileName string `json:"saveFileName"`
}
