package request
type SysApiListSearch struct {
	Pagination Pagination `gorm:"embedded"`
	Description string 	`json:"description"`
	ApiGroup string 	`json:"apiGroup"`
	Method string 	`json:"method"`
	CreatedAt []int64  `json:"createdAt"`
}
