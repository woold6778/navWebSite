package nav

import (
	"nav-web-site/app/api/v1/admin"
	"nav-web-site/mydb"
	"nav-web-site/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddData 添加导航信息数据
// @Summary 添加导航信息数据
// @Description 添加新的导航信息数据
// @Tags nav
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param LoginToken header string true "认证Token"
// @Param class_id formData int true "分类id"
// @Param title formData string true "标题"
// @Param subtitle formData string false "副标题"
// @Param url formData string true "链接地址"
// @Param description formData string false "描述"
// @Param icon formData string false "图标"
// @Param keywords formData string false "关键词"
// @Param sort formData int false "排序"
// @Param is_show formData bool false "是否显示"
// @Param is_recommend formData bool false "是否推荐"
// @Param status formData int false "状态:0=禁用,1=启用"
// @Success 200 {object} util.APIResponse{code=int,message=string}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=string}
// @Router /nav/addData [post]
func AddData(c *gin.Context) {
	var data mydb.StructNav

	loginToken := c.GetHeader("LoginToken")
	adminID, err := admin.GetAdminIDFromToken(loginToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.APIResponse{Code: http.StatusUnauthorized, Message: "无权限添加导航", Data: err.Error()})
		return
	}

	data.Admin_id = adminID

	data.Class_id, _ = strconv.Atoi(c.PostForm("class_id"))
	data.Title = c.PostForm("title")
	data.Subtitle = c.PostForm("subtitle")
	data.Url = c.PostForm("url")
	data.Description = c.PostForm("description")
	data.Icon = c.PostForm("icon")
	data.Keywords = c.PostForm("keywords")
	data.Sort, _ = strconv.Atoi(c.PostForm("sort"))
	data.Is_show, _ = strconv.ParseBool(c.PostForm("is_show"))
	data.Is_recommend, _ = strconv.ParseBool(c.PostForm("is_recommend"))
	data.Status, _ = strconv.Atoi(c.PostForm("status"))
	data.Create_time = util.GetTimestamp(10)
	data.Update_time = util.GetTimestamp(10)

	id, rowsAffected, err := data.Insert([]mydb.StructNav{data})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "添加导航信息失败", Data: err.Error()})
		return
	}
	util.InfoLogger.Printf("导航信息添加成功,id=%d,添加记录数:%d", id, rowsAffected)

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "导航信息添加成功"})
}

// UpdateData 修改导航信息数据
// @Summary 修改导航信息数据
// @Description 根据导航信息ID修改导航信息数据
// @Tags nav
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param LoginToken header string true "认证Token"
// @Param id path string true "导航信息ID"
// @Param class_id formData int true "分类id"
// @Param title formData string true "标题"
// @Param subtitle formData string false "副标题"
// @Param url formData string true "链接地址"
// @Param description formData string false "描述"
// @Param icon formData string false "图标"
// @Param keywords formData string false "关键词"
// @Param sort formData int false "排序"
// @Param is_show formData bool false "是否显示"
// @Param is_recommend formData bool false "是否推荐"
// @Param status formData int false "状态:0=禁用,1=启用"
// @Success 200 {object} util.APIResponse{code=int,message=string}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=string}
// @Router /nav/updateData/{id} [put]
func UpdateData(c *gin.Context) {
	dataID := c.Param("id")
	data, err := (&mydb.StructNav{}).Find(mydb.QueryParams{Condition: "id=" + dataID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取导航信息失败", Data: err.Error()})
		return
	}

	loginToken := c.GetHeader("LoginToken")
	adminID, err := admin.GetAdminIDFromToken(loginToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "未授权操作", Data: err.Error()})
		return
	}

	if data.Admin_id != adminID {
		c.JSON(http.StatusUnauthorized, util.APIResponse{Code: http.StatusUnauthorized, Message: "无权限修改该导航", Data: "null"})
		return
	}

	if c.PostForm("class_id") != "" {
		data.Class_id, _ = strconv.Atoi(c.PostForm("class_id"))
	}
	if c.PostForm("title") != "" {
		data.Title = c.PostForm("title")
	}
	if c.PostForm("subtitle") != "" {
		data.Subtitle = c.PostForm("subtitle")
	}
	if c.PostForm("url") != "" {
		data.Url = c.PostForm("url")
	}
	if c.PostForm("description") != "" {
		data.Description = c.PostForm("description")
	}
	if c.PostForm("icon") != "" {
		data.Icon = c.PostForm("icon")
	}
	if c.PostForm("keywords") != "" {
		data.Keywords = c.PostForm("keywords")
	}
	if c.PostForm("sort") != "" {
		data.Sort, _ = strconv.Atoi(c.PostForm("sort"))
	}
	if c.PostForm("is_show") != "" {
		data.Is_show, _ = strconv.ParseBool(c.PostForm("is_show"))
	}
	if c.PostForm("is_recommend") != "" {
		data.Is_recommend, _ = strconv.ParseBool(c.PostForm("is_recommend"))
	}
	if c.PostForm("status") != "" {
		data.Status, _ = strconv.Atoi(c.PostForm("status"))
	}
	data.Update_time = util.GetTimestamp(10)

	_, _, err = data.Update([]mydb.StructNav{data}, "id="+strconv.Itoa(data.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "修改导航信息失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "导航信息修改成功"})
}

// DeleteData 删除导航信息数据
// @Summary 删除导航信息数据
// @Description 根据导航信息ID删除导航信息数据
// @Tags nav
// @Produce application/json
// @Param id path string true "导航信息ID"
// @Success 200 {object} util.APIResponse{code=int,message=string}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=string}
// @Router /nav/deleteData/{id} [delete]
func DeleteData(c *gin.Context) {
	dataID := c.Param("id")
	data, err := (&mydb.StructNav{}).Find(mydb.QueryParams{Condition: "id=" + dataID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取导航信息失败", Data: err.Error()})
		return
	}

	_, _, err = data.Delete("id=" + strconv.Itoa(data.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "删除导航信息失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "导航信息删除成功"})
}

// GetDataList 获取导航信息列表
// @Summary 获取导航信息列表
// @Description 获取所有导航信息的列表
// @Tags nav
// @Produce application/json
// @Success 200 {object} util.APIResponse{code=int,message=string,data=[]mydb.StructNav}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=string}
// @Router /nav/getList [get]
func GetDataList(c *gin.Context) {
	params := mydb.QueryParams{}
	dataList, _, err := (&mydb.StructNav{}).Select(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取导航信息列表失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "获取导航信息列表成功", Data: dataList})
}

// GetDataDetail 获取导航信息详情
// @Summary 获取导航信息详情
// @Description 根据导航信息ID获取导航信息详情
// @Tags nav
// @Produce application/json
// @Param id path string true "导航信息ID"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=mydb.StructNav}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=string}
// @Router /nav/getDetail/{id} [get]
func GetDataDetail(c *gin.Context) {
	dataID := c.Param("id")
	data, err := (&mydb.StructNav{}).Find(mydb.QueryParams{Condition: "id=" + dataID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取导航信息详情失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "获取导航信息详情成功", Data: data})
}
