package v1

import (
	"ginProject/global"
	"ginProject/model"
	"ginProject/model/request"
	"ginProject/model/response"
	"ginProject/service"
	"ginProject/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Tags SysRole
// @Summary 分页获取角色列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.Pagination true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api/sysRole/getSysRoleList [GET]
func GetSysRoleList(c *gin.Context) {
	var U request.Pagination
	_ = c.ShouldBindJSON(&U)
	err, total, list := service.GetSysRoleList(&U)
	if err != nil {
		global.GvaLog.Error("获取角色列表失败！", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Current:  U.Current,
			PageSize: U.PageSize,
		}, "获取数据成功", c)
	}
}

// @Tags SysRole
// @Summary 更新角色信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysRole true "角色模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /api/sysRole/updateSysRoleInfo [PUT]
func UpdateSysRoleInfo(c *gin.Context) {
	var role model.SysRole
	_ = c.ShouldBindJSON(&role)
	if err, ret := service.UpdateSysRoleInfo(role); err != nil {
		global.GvaLog.Error("更新失败！", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(gin.H{"roleInfo": ret}, "更新成功", c)
	}
}

// @Tags SysRole
// @Summary 新增角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysRole true "角色模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"新增成功"}"
// @Router /api/sysRole/addSysRoleInfo [POST]
func AddSysRoleInfo(c *gin.Context) {
	var role model.SysRole
	_ = c.ShouldBindJSON(&role)
	if err, ret := service.AddSysRoleInfo(role); err != nil {
		global.GvaLog.Error("新增失败！", zap.Any("err", err))
		response.FailWithDetailed(response.SysRoleResponse{Role: ret}, err.Error(), c)
	} else {
		response.OkWithDetailed(gin.H{"roleInfo": ret}, "新增成功", c)
	}
}

// @Tags SysRole
// @Summary 批量删除角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "id列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /api/sysRole/deleteBatchSysRole [DELETE]
func DeleteBatchSysRole(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteBatchSysRole(IDS); err != nil {
		global.GvaLog.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// @Tags SysRole
// @Summary 配置角色菜单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysRoleMenuConfig true "role_id,id列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"配置角色菜单成功"}"
// @Router /api/sysRole/updateSysRoleMenuConfig [PUT]
func UpdateSysRoleMenuConfig(c *gin.Context) {
	var s request.SysRoleMenuConfig
	_ = c.ShouldBindJSON(&s)
	if err := service.UpdateSysRoleMenuConfig(s); err != nil {
		global.GvaLog.Error("配置角色菜单失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("配置角色菜单成功", c)
	}
}

// @Tags SysRole
// @Summary 配置角色菜单首页
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysRoleDefaultRouter true "role_id,default_router"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"配置角色接口成功"}"
// @Router /api/sysRole/setRoleDefaultRouter [PUT]
func SetRoleDefaultRouter(c *gin.Context) {
	var s request.SysRoleDefaultRouter
	_ = c.ShouldBindJSON(&s)
	if err := service.SetRoleDefaultRouter(s); err != nil {
		global.GvaLog.Error("配置角色菜单首页失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("配置角色菜单首页成功", c)
	}
}

// @Tags SysRole
// @Summary 批量导出角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysRoleExcelOut true "批量导出请求模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量导出成功"}"
// @Router /api/sysRole/excelOut [POST]
func ExcelOutSysRole(c *gin.Context) {
	var E request.SysRoleExcelOut
	_ = c.ShouldBindJSON(&E)
	err, excelFilePath := service.ExcelOutSysRole(E)
	if err != nil {
		global.GvaLog.Error("批量导出失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		// 如果文件路径不为空，也需要删除文件
		if excelFilePath != "" {
			err = utils.DeleteLocalFile(excelFilePath)
			if err != nil {
				global.GvaLog.Error("删除文件失败!", zap.Any("err", err))
				return
			}
		}
		return
	}
	// 要设置这一条，要不然前端获取不到Content-Disposition
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", "roleTable.xlsx")
	c.File(excelFilePath)
	// 删除Excel文件
	err = utils.DeleteLocalFile(excelFilePath)
	if err != nil {
		global.GvaLog.Error("删除文件失败!", zap.Any("err", err))
		return
	}
	return
}
