package admin

import (
	"encoding/json"
	"fmt"
	"nav-web-site/config"
	"nav-web-site/mydb"
	"nav-web-site/util"
	"nav-web-site/util/log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginSuccessData struct {
	Username string
	Token    string
}

// Login 管理员登录
// @Summary 管理员登录
// @Description 管理员登录
// @Tags admin
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param username formData string true "用户名"
// @Param password formData string true "登录密码(进行过1次MD5的密码)"
// @Param expiration formData int false "过期时间(分钟), 默认120分钟"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=LoginSuccessData} "登录成功"
// @Failure 400 {object} util.APIResponse{code=int,message=string,data=object} "无效的过期时间参数"
// @Failure 401 {object} util.APIResponse{code=int,message=string,data=object} "用户名或密码错误"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=object} "查询管理员信息失败"
// @Router /admin/login [post]
func Login(c *gin.Context) {
	// 记录接收到的请求
	log.InfoLogger.Printf("Received login request - Username: %s, ClientIP: %s, User-Agent: %s", c.PostForm("username"), c.ClientIP(), c.Request.UserAgent())

	//接收用户名和密码，从mydb.TAbleS.Admind的Find函数进行查询(用户名作为查询条件)
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, util.APIResponse{Code: http.StatusBadRequest, Message: "用户名和密码是必须的", Data: "null"})
		return
	}

	// 创建查询参数
	params := mydb.QueryParams{
		Condition: "username='" + username + "'",
	}
	log.InfoLogger.Printf("Received login request - Username: %s, ClientIP: %s", username, c.ClientIP())

	// 查询管理员信息
	admin, err := mydb.Tables.Admin.Find(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "查询管理员信息失败", Data: err.Error()})
		return
	}

	// 校验密码,提交的密码是MD5过1次的，我们这里对提交的密码拼接admin.salt，然后再次MD5
	// 对提交的密码拼接admin.Salt，然后再次MD5
	submittedPassword := util.MD5Hash(password, admin.Salt)

	// 校验密码
	if admin.Password != submittedPassword {
		c.JSON(http.StatusUnauthorized, util.APIResponse{Code: http.StatusUnauthorized, Message: "用户名或密码错误", Data: "null"})
		return
	}

	// 生成登录凭证
	timestamp := util.GetTimestamp(10)
	randomString := util.GenerateRandomString(6, 1)
	tokenStr := username + fmt.Sprintf("%d", timestamp) + randomString
	login_token := util.MD5Hash(tokenStr, admin.Salt)
	// 将token缓存起来，key为"admin_login_token_" + login_token，内容为admin结构体，缓存时间为120分钟
	cacheKey := "admin_login_token_" + login_token
	// 获取客户端IP和浏览器信息
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	deviceFingerprint := util.GenerateDeviceFingerprint(c.Request)
	// 将admin结构体、客户端IP和浏览器信息缓存起来
	cacheValue := map[string]interface{}{
		"admin":              admin,
		"client_ip":          clientIP,
		"user_agent":         userAgent,
		"device_fingerprint": deviceFingerprint,
	}

	// 获取前端提交的过期时间参数
	expiration := c.PostForm("expiration")
	var cacheDuration time.Duration
	if expiration != "" {
		expirationInt, err := strconv.Atoi(expiration)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.APIResponse{Code: http.StatusBadRequest, Message: "无效的过期时间参数", Data: err.Error()})
			return
		}
		cacheDuration = time.Duration(expirationInt) * time.Minute
	} else {
		cacheDuration = 120 * time.Minute
	}

	util.C.Set(cacheKey, cacheValue, cacheDuration)

	c.SetCookie("session_id", login_token, int(cacheDuration.Seconds()), "/", config.Config.Base.SiteDomain, false, true)
	// 返回登录凭证
	data := LoginSuccessData{Username: admin.Username, Token: login_token}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "登录成功", Data: data})
}

// Register 注册新管理员
// @Summary 注册新管理员
// @Description 通过接收前端传递的参数，注册一个新的管理员账户
// @Tags admin
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Param email formData string false "邮箱地址"
// @Param phone_number formData string false "手机号"
// @Param avatar formData string false "用户头像"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=object} "注册成功"
// @Failure 400 {object} util.APIResponse{code=int,message=string,data=object} "无效的参数"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=object} "注册失败"
// @Router /admin/register [post]
func Register(c *gin.Context) {
	var admin mydb.StructAdmin

	admin.Username = c.PostForm("username")
	salt := util.GenerateRandomString(8, 1)
	admin.Password = util.MD5Hash(c.PostForm("password"), salt)
	admin.Email = c.PostForm("email")
	if admin.Email == "" {
		admin.Email = ""
	} else if !util.IsValidEmail(admin.Email) {
		c.JSON(http.StatusBadRequest, util.APIResponse{Code: http.StatusBadRequest, Message: "无效的邮箱地址", Data: "null"})
		return
	}

	admin.PhoneNumber = c.PostForm("phone_number")
	if admin.PhoneNumber == "" {
		admin.PhoneNumber = ""
	} else if !util.IsValidPhoneNumber(admin.PhoneNumber, "CN") { // 假设国家代码为 "CN"
		c.JSON(http.StatusBadRequest, util.APIResponse{Code: http.StatusBadRequest, Message: "无效的手机号", Data: "null"})
		return
	}

	admin.Status = 1
	admin.CreateTime = util.GetTimestamp(10)
	admin.UpdateTime = util.GetTimestamp(10)
	admin.LastLoginTime = util.GetTimestamp(10)
	admin.Role = ""
	admin.Salt = salt
	admin.Avatar = c.PostForm("avatar")
	if admin.Avatar == "" {
		admin.Avatar = ""
	}

	// 获取必填字段
	requiredFields := admin.GetRequiredFields()
	for _, field := range requiredFields {
		if reflect.ValueOf(admin).FieldByName(field).String() == "" {
			c.JSON(http.StatusBadRequest, util.APIResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("%s 是必填字段", field), Data: "null"})
			return
		}
	}

	_, err := admin.Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "注册失败", Data: err.Error()})
		return
	}

	dataJSON, err := json.Marshal(map[string]interface{}{"username": admin.Username})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "转换数据失败", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "注册成功", Data: string(dataJSON)})
}
