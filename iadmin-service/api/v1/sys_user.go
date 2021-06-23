package v1

import (
	"fmt"
	"ginProject/global"
	"ginProject/middleware"
	"ginProject/model"
	"ginProject/model/request"
	"ginProject/model/response"
	"ginProject/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"log"
	"time"
)


// @Tags Base
// @Summary 用户登录
// @Produce  application/json
// @Param data body request.Login true "用户名, 密码, 验证码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登录成功"}"
// @Router /api/base/login [post]
func Login(c *gin.Context)  {
	var L request.Login
	_ = c.ShouldBindJSON(&L)
	if store.Verify(L.CaptchaId,L.Captcha,true){
		U := &model.SysUser{Phone: L.Phone,Password: L.Password}
		if err, user := service.Login(U);err != nil{
			log.Printf("登录失败！用户名不存在或者密码错误！error:%s\n", err)
			response.FailWithMessage("用户名不存在或者密码错误！",c)
		}else{
			tokenNext(c, *user)
		}
	}else {
		response.FailWithMessage("验证码错误！",c)
	}
}

// 登录以后签发jwt
func tokenNext(c *gin.Context, user model.SysUser) {
	j := &middleware.JWT{SigningKey: []byte(global.GvaConfig.JWT.SigningKey)} // 唯一签名
	claims := request.CustomClaims{
		ID:          user.ID,
		RealName:    user.RealName,
		Phone:    user.Phone,
		RoleId:    user.RoleId,
		BufferTime:  global.GvaConfig.JWT.BufferTime, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,                             // 签名生效时间
			ExpiresAt: time.Now().Unix() + global.GvaConfig.JWT.ExpiresTime, // 过期时间 7天  配置文件
			Issuer:    global.GvaConfig.JWT.SigningKey,                      // 签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		log.Printf("获取token失败!%s\n",err)
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
			log.Printf("设置登录状态失败！%s\n",err)
			response.FailWithMessage("设置登录状态失败！", c)
			return
		}
		response.OkWithDetailed(response.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
	} else if err != nil {
		log.Printf("设置登录状态失败！%s\n",err)
		response.FailWithMessage("设置登录状态失败", c)
	} else {
		var blackJWT model.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := service.JsonInBlacklist(blackJWT); err != nil {
			response.FailWithMessage("jwt作废失败", c)
			return
		}
		if err := service.SetRedisJWT(token, user.Phone); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
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
func Register(c *gin.Context)  {
	var R request.Register
	_ = c.ShouldBindJSON(&R)
	if store.Verify(R.CaptchaId,R.Captcha,true){
		user :=  &model.SysUser{Phone: R.Phone,Password: R.Password,RealName: R.RealName}
		err, userReturn := service.Register(*user)
		if err != nil{
			errMsg := fmt.Sprintf("注册失败！%s\n",err)
			log.Printf(errMsg)
			response.FailWithDetailed(response.SysUserResponse{User: userReturn},errMsg,c)
		}else{
			response.OkWithDetailed(response.SysUserResponse{User: userReturn}, "注册成功", c)
		}
	}else {
		response.FailWithMessage("验证码错误！",c)
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
func GetSysUserList(c *gin.Context)  {
	var U request.SysUserListSearch
	_ = c.ShouldBindJSON(&U)
	err, total,list := service.GetSysUserList(&U)
	if  err != nil{
		response.FailWithMessage("获取用户列表失败！",c)
	}else{
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Current:     U.Pagination.Current,
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
		response.FailWithMessage("更新失败！", c)
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
func AddSysUserInfo(c *gin.Context)  {
	var user model.SysUser
	_ = c.ShouldBindJSON(&user)
	if err, ReqUser := service.AddSysUserInfo(user); err != nil {
		global.GvaLog.Error("新增失败！", zap.Any("err", err))
		response.FailWithDetailed(response.SysUserResponse{User: ReqUser},"新增失败！",c)
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
func DeleteBatchSysUser(c *gin.Context){
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteBatchSysUser(IDS); err != nil {
		global.GvaLog.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}