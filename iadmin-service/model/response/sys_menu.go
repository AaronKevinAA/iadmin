package response

import "ginProject/model"

type SysMenuResponse struct {
	Menu model.SysMenu `json:"menu"`
}

