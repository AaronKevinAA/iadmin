package v1

import (
	"ginProject/global"
	"ginProject/model/request"
	"ginProject/model/response"
	"ginProject/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Tags Casbin
// @Summary 更新角色api权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CasbinInReceive true "权限id, 权限模型列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /api/casbin/updateCasbin [post]
func UpdateCasbin(c *gin.Context) {
	var cmr request.CasbinInReceive
	_ = c.ShouldBindJSON(&cmr)
	if err := service.UpdateCasbin(cmr.RoleId, cmr.CasbinInfos); err != nil {
		global.GvaLog.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Casbin
// @Summary 获取权限列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CasbinInReceive true "权限id, 权限模型列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api/casbin/getPolicyPathByRoleId [post]
func GetPolicyPathByRoleId(c *gin.Context) {
	var casbin request.CasbinInReceive
	_ = c.ShouldBindJSON(&casbin)
	paths := service.GetPolicyPathByRoleId(casbin.RoleId)
	response.OkWithDetailed(response.PolicyPathResponse{Paths: paths}, "获取成功", c)
}
