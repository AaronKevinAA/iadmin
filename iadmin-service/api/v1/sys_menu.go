package v1

import (
	"ginProject/global"
	"ginProject/model"
	"ginProject/model/request"
	"ginProject/model/response"
	"ginProject/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Tags SysMenu
// @Summary 分页获取路由菜单列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api/sysMenu/getSysRouteList [POST]
func GetSysRouteList(c *gin.Context)  {
	err,tree := service.GetSysRouteList()
	if  err != nil{
		response.FailWithMessage("获取路由菜单列表失败！",c)
	}else{
		response.OkWithDetailed(response.PageResult{
			List:     tree,
			Total: int64(len(tree)),
			Current:  0,
			PageSize: 0,
		}, "获取数据成功", c)
	}
}

// @Tags SysMenu
// @Summary 更新菜单信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysMenu true "菜单模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /api/sysMenu/updateSysMenuInfo [PUT]
func UpdateSysMenuInfo(c *gin.Context) {
	var menu model.SysMenu
	_ = c.ShouldBindJSON(&menu)
	if err, ret := service.UpdateSysMenuInfo(menu); err != nil {
		global.GvaLog.Error("更新失败！", zap.Any("err", err))
		response.FailWithMessage("更新失败！", c)
	} else {
		response.OkWithDetailed(gin.H{"menuInfo": ret}, "更新成功", c)
	}
}

// @Tags SysMenu
// @Summary 新增菜单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysMenu true "菜单模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"新增成功"}"
// @Router /api/sysMenu/addSysMenuInfo [POST]
func AddSysMenuInfo(c *gin.Context)  {
	var menu model.SysMenu
	_ = c.ShouldBindJSON(&menu)
	if err, ret := service.AddSysMenuInfo(menu); err != nil {
		global.GvaLog.Error("新增失败！", zap.Any("err", err))
		response.FailWithDetailed(response.SysMenuResponse{Menu: ret},"新增失败！",c)
	} else {
		response.OkWithDetailed(gin.H{"menuInfo": ret}, "新增成功", c)
	}
}

// @Tags SysMenu
// @Summary 批量删除菜单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "id列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /api/sysMenu/deleteBatchSysMenu [DELETE]
func DeleteBatchSysMenu(c *gin.Context){
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteBatchSysMenu(IDS); err != nil {
		global.GvaLog.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败！", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}
// @Tags SysMenu
// @Summary 获得某人的菜单树
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api/sysMenu/getSysMenuByToken [GET]
func GetSysMenuByToken(c *gin.Context)  {
	var roleId uint
	if claims, ok := c.Get("claims"); ok {
		waitUse := claims.(*request.CustomClaims)
		roleId = waitUse.RoleId
	}
	err,routes,menus := service.GetSysMenuByToken(roleId)
	if err != nil{
		global.GvaLog.Error("获得菜单失败!", zap.Any("err", err))
		response.FailWithMessage("获得菜单失败！", c)
	}else{
		response.OkWithDetailed(gin.H{"routes":routes,"menus":menus},"获得菜单成功", c)
	}
}