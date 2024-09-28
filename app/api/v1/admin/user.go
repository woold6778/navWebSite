package admin

import (
	"nav-web-site/mydb"
	"nav-web-site/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 获取所有用户的列表
// @Tags admin
// @Produce application/json
// @Param LoginToken header string true "认证Token"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=[]mydb.StructAdmin{password=string,salt=string}}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}}
// @Router /admin/list [get]
func GetUserList(c *gin.Context) {
	loginToken := c.GetHeader("LoginToken")
	_, err := GetAdminIDFromToken(loginToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.APIResponse{Code: http.StatusUnauthorized, Message: "无权限进行此操作", Data: err.Error()})
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 20
	}
	params := mydb.QueryParams{
		Page:     page,
		PageSize: pageSize,
	}
	users, _, err := mydb.Tables.Admin.Select(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取用户列表失败", Data: err.Error()})
		return
	}

	// 移除password和salt字段
	for i := range users {
		users[i].Password = ""
		users[i].Salt = ""
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: -1, Message: "获取用户列表成功", Data: users})
}

// GetUserDetail 获取用户详情
// @Summary 获取用户详情
// @Description 根据用户ID获取用户的详细信息
// @Tags admin
// @Produce application/json
// @Param LoginToken header string true "认证Token"
// @Param id path string true "用户ID"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=mydb.StructAdmin}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}}
// @Router /admin/detail/{id} [get]
func GetUserDetail(c *gin.Context) {
	loginToken := c.GetHeader("LoginToken")
	_, err := GetAdminIDFromToken(loginToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.APIResponse{Code: http.StatusUnauthorized, Message: "无权限进行此操作", Data: err.Error()})
		return
	}

	userID := c.Param("id")
	params := mydb.QueryParams{
		Condition: "id=" + userID,
	}
	user, err := mydb.Tables.Admin.Find(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取用户详情失败", Data: err.Error()})
		return
	}
	user.Password = ""
	user.Salt = ""
	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "获取用户详情成功", Data: user})
}

// UpdateUserPassword 修改用户密码
// @Summary 修改用户密码
// @Description 根据用户ID修改用户的密码
// @Tags admin
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param LoginToken header string true "认证Token"
// @Param id path string true "用户ID"
// @Param password formData string true "新密码"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=interface{}}
// @Failure 400 {object} util.APIResponse{code=int,message=string,data=interface{}}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}}
// @Router /admin/updatePassword/{id} [put]
func UpdateUserPassword(c *gin.Context) {
	loginToken := c.GetHeader("LoginToken")
	_, err := GetAdminIDFromToken(loginToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.APIResponse{Code: http.StatusUnauthorized, Message: "无权限进行此操作", Data: err.Error()})
		return
	}

	userID := c.Param("id")
	newPassword := c.PostForm("password")
	if newPassword == "" {
		c.JSON(http.StatusBadRequest, util.APIResponse{Code: http.StatusBadRequest, Message: "密码不能为空"})
		return
	}

	user, err := mydb.Tables.Admin.Find(mydb.QueryParams{Condition: "id=" + userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取用户信息失败", Data: err.Error()})
		return
	}

	user.Password = util.MD5Hash(newPassword, user.Salt)
	_, err = user.Update("id=" + userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "修改密码失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "密码修改成功"})
}

// EditUserProfile 编辑用户资料
// @Summary 编辑用户资料
// @Description 根据用户ID编辑用户的资料
// @Tags admin
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param LoginToken header string true "认证Token"
// @Param id path string true "用户ID"
// @Param email formData string false "邮箱地址"
// @Param phone_number formData string false "手机号"
// @Param avatar formData string false "用户头像"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=interface{}}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}}
// @Router /admin/editProfile/{id} [put]
func EditUserProfile(c *gin.Context) {
	loginToken := c.GetHeader("LoginToken")
	_, err := GetAdminIDFromToken(loginToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.APIResponse{Code: http.StatusUnauthorized, Message: "无权限进行此操作", Data: err.Error()})
		return
	}

	userID := c.Param("id")
	user, err := mydb.Tables.Admin.Find(mydb.QueryParams{Condition: "id=" + userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取用户信息失败", Data: err.Error()})
		return
	}

	user.Email = c.PostForm("email")
	user.PhoneNumber = c.PostForm("phone_number")
	user.Avatar = c.PostForm("avatar")

	_, err = user.Update("id=" + userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "编辑用户资料失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "用户资料编辑成功"})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 根据用户ID删除用户
// @Tags admin
// @Produce application/json
// @Param LoginToken header string true "认证Token"
// @Param id path string true "用户ID"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=interface{}}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}}
// @Router /admin/delete/{id} [delete]
func DeleteUser(c *gin.Context) {
	loginToken := c.GetHeader("LoginToken")
	_, err := GetAdminIDFromToken(loginToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.APIResponse{Code: http.StatusUnauthorized, Message: "无权限进行此操作", Data: err.Error()})
		return
	}

	userID := c.Param("id")
	user, err := mydb.Tables.Admin.Find(mydb.QueryParams{Condition: "id=" + userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取用户信息失败", Data: err.Error()})
		return
	}

	_, err = user.Delete("id=" + userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "删除用户失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "用户删除成功"})
}
