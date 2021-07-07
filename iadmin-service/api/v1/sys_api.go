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

// @Tags SysApi
// @Summary 分页获取接口列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysApiListSearch true "页码, 每页大小，搜索条件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取数据成功"}"
// @Router /api/sysApi/getSysApiList [POST]
func GetSysApiList(c *gin.Context) {
	var U request.SysApiListSearch
	_ = c.ShouldBindJSON(&U)
	err, total, list := service.GetSysApiList(&U)
	if err != nil {
		response.FailWithMessage("获取接口列表失败！", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Current:  U.Pagination.Current,
			PageSize: U.Pagination.PageSize,
		}, "获取数据成功", c)
	}
}

// @Tags SysApi
// @Summary 更新接口信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysApi true "接口模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /api/sysApi/updateSysApiInfo [PUT]
func UpdateSysApiInfo(c *gin.Context) {
	var Api model.SysApi
	_ = c.ShouldBindJSON(&Api)
	if err := service.UpdateSysApiInfo(Api); err != nil {
		global.GvaLog.Error("更新失败！", zap.Any("err", err))
		response.FailWithMessage("更新失败！", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags SysApi
// @Summary 新增接口
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysApi true "接口模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"新增成功"}"
// @Router /api/sysApi/addSysApiInfo [POST]
func AddSysApiInfo(c *gin.Context) {
	var Api model.SysApi
	_ = c.ShouldBindJSON(&Api)
	if err, ret := service.AddSysApiInfo(Api); err != nil {
		global.GvaLog.Error("新增失败！", zap.Any("err", err))
		response.FailWithDetailed(response.SysApiResponse{Api: ret}, "新增失败！", c)
	} else {
		response.OkWithDetailed(gin.H{"ApiInfo": ret}, "新增成功", c)
	}
}

// @Tags SysApi
// @Summary 批量删除接口
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "id列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /api/sysApi/deleteBatchSysApi [DELETE]
func DeleteBatchSysApi(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteBatchSysApi(IDS); err != nil {
		global.GvaLog.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// @Tags SysApi
// @Summary 获得接口树
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取数据成功"}"
// @Router /api/sysApi/getSysApiTree [GET]
func GetSysApiTree(c *gin.Context) {
	if err, data := service.GetSysApiTree(); err != nil {
		global.GvaLog.Error("获取数据失败!", zap.Any("err", err))
		response.FailWithMessage("获取数据失败", c)
	} else {
		response.OkWithDetailed(data, "获取数据成功", c)
	}
}
