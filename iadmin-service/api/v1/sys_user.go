package v1

import (
	"ginProject/global"
	"ginProject/middleware"
	"ginProject/model"
	"ginProject/model/request"
	"ginProject/model/response"
	"ginProject/service"
	"ginProject/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"time"
)

// @Tags Base
// @Summary 用户登录
// @Produce  application/json
// @Param data body request.Login true "用户名, 密码, 验证码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登录成功"}"
// @Router /api/base/login [post]
func Login(c *gin.Context) {
	var L request.Login
	_ = c.ShouldBindJSON(&L)
	if store.Verify(L.CaptchaId, L.Captcha, true) {
		if err, user := service.Login(&L); err != nil {
			global.GvaLog.Error("登录失败！用户名不存在或者密码错误！", zap.Any("err", err))
			response.FailWithMessage("用户名不存在或者密码错误！", c)
		} else {
			tokenNext(c, *user)
		}
	} else {
		response.FailWithMessage("验证码错误！", c)
	}
}

// 登录以后签发jwt
func tokenNext(c *gin.Context, user model.SysUser) {
	j := &middleware.JWT{SigningKey: []byte(global.GvaConfig.JWT.SigningKey)} // 唯一签名
	claims := request.CustomClaims{
		ID:         user.ID,
		RealName:   user.RealName,
		Phone:      user.Phone,
		RoleId:     user.RoleId,
		BufferTime: global.GvaConfig.JWT.BufferTime, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,                             // 签名生效时间
			ExpiresAt: time.Now().Unix() + global.GvaConfig.JWT.ExpiresTime, // 过期时间 7天  配置文件
			Issuer:    global.GvaConfig.JWT.SigningKey,                      // 签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		global.GvaLog.Error("获取token失败！", zap.Any("err", err))
		response.FailWithMessage("获取token失败！", c)
		return
	}
	// 多点登录拦截
	if !global.GvaConfig.System.UseMultipoint {
		response.OkWithDetailed(response.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
		return
	}
	if err, jwtStr := service.GetRedisJWT(user.Phone); err == redis.Nil {
		if err := service.SetRedisJWT(token, user.Phone); err != nil {
			global.GvaLog.Error("设置登录状态失败！", zap.Any("err", err))
			response.FailWithMessage("设置登录状态失败！", c)
			return
		}
		response.OkWithDetailed(response.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
	} else if err != nil {
		global.GvaLog.Error("设置登录状态失败！", zap.Any("err", err))
		response.FailWithMessage("设置登录状态失败！", c)
	} else {
		var blackJWT model.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := service.JsonInBlacklist(blackJWT); err != nil {
			global.GvaLog.Error("jwt作废失败！", zap.Any("err", err))
			response.FailWithMessage("jwt作废失败！", c)
			return
		}
		if err := service.SetRedisJWT(token, user.Phone); err != nil {
			global.GvaLog.Error("设置登录状态失败！", zap.Any("err", err))
			response.FailWithMessage("设置登录状态失败！", c)
			return
		}
		response.OkWithDetailed(response.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
	}
}

// @Tags Base
// @Summary 用户注册
// @Produce  application/json
// @Param data body request.Register true "用户名, 密码, 真实姓名, 验证码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"注册成功"}"
// @Router /api/base/register [post]
func Register(c *gin.Context) {
	var R request.Register
	_ = c.ShouldBindJSON(&R)
	if store.Verify(R.CaptchaId, R.Captcha, true) {
		err, userReturn := service.Register(R)
		if err != nil {
			global.GvaLog.Error("注册失败！", zap.Any("err", err))
			response.FailWithDetailed(response.SysUserResponse{User: userReturn}, err.Error(), c)
		} else {
			response.OkWithDetailed(response.SysUserResponse{User: userReturn}, "注册成功", c)
		}
	} else {
		global.GvaLog.Error("验证码错误！")
		response.FailWithMessage("验证码错误！", c)
	}
}

// @Tags SysUser
// @Summary 分页获取用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysUserListSearch true "页码, 每页大小，搜索条件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api/sysUser/getSysUserList [POST]
func GetSysUserList(c *gin.Context) {
	var U request.SysUserListSearch
	_ = c.ShouldBindJSON(&U)
	err, total, list := service.GetSysUserList(&U)
	if err != nil {
		global.GvaLog.Error("获取用户列表失败！", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Current:  U.Pagination.Current,
			PageSize: U.Pagination.PageSize,
		}, "获取数据成功", c)
	}
}

// @Tags SysUser
// @Summary 更新用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysUser true "手机号，真实姓名"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /api/sysUser/updateSysUserInfo [PUT]
func UpdateSysUserInfo(c *gin.Context) {
	var user model.SysUser
	_ = c.ShouldBindJSON(&user)
	if err, ReqUser := service.UpdateSysUserInfo(user); err != nil {
		global.GvaLog.Error("更新失败！", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(gin.H{"userInfo": ReqUser}, "更新成功", c)
	}
}

// @Tags SysUser
// @Summary 新增用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysUser true "手机号，真实姓名"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"新增成功"}"
// @Router /api/sysUser/addSysUserInfo [POST]
func AddSysUserInfo(c *gin.Context) {
	var user model.SysUser
	_ = c.ShouldBindJSON(&user)
	if err, ReqUser := service.AddSysUserInfo(user); err != nil {
		global.GvaLog.Error("新增失败！", zap.Any("err", err))
		response.FailWithDetailed(response.SysUserResponse{User: ReqUser}, err.Error(), c)
	} else {
		response.OkWithDetailed(gin.H{"userInfo": ReqUser}, "新增成功", c)
	}
}

// @Tags SysUser
// @Summary 批量删除用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "id列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /api/sysUser/deleteBatchSysUser [DELETE]
func DeleteBatchSysUser(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteBatchSysUser(IDS); err != nil {
		global.GvaLog.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// @Tags SysUser
// @Summary 批量导出用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysUserExcelOut true "批量导出请求模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量导出成功"}"
// @Router /api/sysUser/excelOut [POST]
func ExcelOutSysUser(c *gin.Context) {
	var E request.SysUserExcelOut
	_ = c.ShouldBindJSON(&E)
	err, excelFilePath := service.ExcelOutSysUser(E)
	if err != nil {
		global.GvaLog.Error("批量导出失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		// 如果文件路径不为空，则页需要删除文件
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
	c.Writer.Header().Add("Content-Disposition", "userTable.xlsx")
	c.File(excelFilePath)
	// 删除Excel文件
	err = utils.DeleteLocalFile(excelFilePath)
	if err != nil {
		global.GvaLog.Error("删除文件失败!", zap.Any("err", err))
		return
	}
	return
}

// @Tags SysUser
// @Summary 批量导入用户预览
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param file formData file true "上传文件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量导入预览成功"}"
// @Router /api/sysUser/excelInPreview [POST]
func ExcelInPreviewSysUser(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GvaLog.Error("接收文件失败!", zap.Any("err", err))
		response.FailWithMessage("接收文件失败！", c)
		return
	}
	if header != nil {
		errPreview, newFileName, dataList, allDataCorrect := service.ExcelInPreviewSysUser(header)
		if errPreview != nil {
			global.GvaLog.Error("批量导入预览失败!", zap.Any("err", errPreview))
			response.FailWithMessage(errPreview.Error(), c)
			return
		}
		response.OkWithDetailed(gin.H{"saveFileName": newFileName, "dataList": dataList, "allDataCorrect": allDataCorrect}, "批量导入预览成功", c)
		return
	}
	response.FailWithMessage("批量导入预览失败！", c)
}

// @Tags SysUser
// @Summary 批量导入用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ExcelInRequest true "保存的文件名称"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量导入成功"}"
// @Router /api/sysUser/excelIn [POST]
func ExcelInSysUser(c *gin.Context) {
	var E request.ExcelInRequest
	_ = c.ShouldBindJSON(&E)
	err := service.ExcelInSysUser(E)
	if err != nil {
		global.GvaLog.Error("批量导入失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("批量导入成功", c)
}

// @Tags SysUser
// @Summary 重置用户密码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysUserID true "用户ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"重置密码成功"}"
// @Router /api/sysUser/resetPassword [PUT]
func ResetSysUserPassword(c *gin.Context) {
	var SysUserID request.SysUserID
	_ = c.ShouldBindJSON(&SysUserID)
	err := service.ResetSysUserPassword(SysUserID.ID)
	if err != nil {
		global.GvaLog.Error("重置密码失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("重置密码成功", c)
}

// @Tags SysUser
// @Summary 修改用户密码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.UpdatePasswordByToken true "修改密码模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改密码成功"}"
// @Router /api/sysUser/updatePasswordByToken [PUT]
func UpdatePasswordByToken(c *gin.Context) {
	var U request.UpdatePasswordByToken
	_ = c.ShouldBindJSON(&U)
	var userId uint
	if claims, ok := c.Get("claims"); ok {
		waitUse := claims.(*request.CustomClaims)
		userId = waitUse.ID
	} else {
		response.FailWithMessage("无法获取登录信息!", c)
		return
	}

	err := service.UpdatePasswordByToken(U, userId)
	if err != nil {
		global.GvaLog.Error("修改密码失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"login": true}, "修改密码成功", c)
}

// @Tags SysUser
// @Summary 修改用户基本信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysUserBasicInfo true "修改用户基本信息模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改基本信息成功"}"
// @Router /api/sysUser/updateBasicInfoByToken [PUT]
func UpdateBasicInfoByToken(c *gin.Context) {
	var U request.SysUserBasicInfo
	_ = c.ShouldBindJSON(&U)
	var userId uint
	if claims, ok := c.Get("claims"); ok {
		waitUse := claims.(*request.CustomClaims)
		userId = waitUse.ID
	} else {
		response.FailWithMessage("无法获取登录信息!", c)
		return
	}
	err, newUser := service.UpdateBasicInfoByToken(U, userId)
	if err != nil {
		global.GvaLog.Error("修改基本信息失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(newUser, "修改基本信息成功", c)
}
