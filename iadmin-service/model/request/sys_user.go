package request

type Login struct {
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
	RedisKey  string `json:"redisKey"`
}

type Register struct {
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	RealName  string `json:"realName"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
	RedisKey  string `json:"redisKey"`
}

type SysUserListSearch struct {
	Pagination Pagination `gorm:"embedded"`
	Phone      string     `json:"phone"`
	RealName   string     `json:"realName"`
	CreatedAt  []int64    `json:"createdAt"`
}

type SysUserExcelOut struct {
	SysUserListSearch SysUserListSearch `json:"sysUserListSearch"`
	ExcelOutConfig    ExcelOutRequest   `json:"excelOutConfig"`
}

type SysUserID struct {
	ID uint `json:"userId"`
}

type UpdatePasswordByToken struct {
	OldPassword     string `json:"oldPassword"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
	RedisKey        string `json:"redisKey"`
}

type SysUserBasicInfo struct {
	Phone    string `json:"phone"`
	RealName string `json:"realName"`
}
