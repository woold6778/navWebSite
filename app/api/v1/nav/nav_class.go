package nav

import (
	"nav-web-site/app/api/v1/admin"
	"nav-web-site/mydb"
	"nav-web-site/util"
	"nav-web-site/util/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddClass 添加导航分类
// @Summary 添加导航分类
// @Description 添加新的导航分类
// @Tags nav
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param LoginToken header string true "认证Token"
// @Param name formData string true "导航分类名称"
// @Param parent_id formData int true "父级分类id"
// @Param sort formData string false "排序"
// @Param icon formData string false "图标"
// @Param is_show formData bool false "是否显示"
// @Param is_recommend formData bool false "是否推荐"
// @Param is_hot formData bool false "是否热门"
// @Param status formData int false "状态:0=禁用,1=启用"
// @Param description formData string false "导航分类描述"
// @Success 200 {object} util.APIResponse{code=int,message=string}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=string}
// @Router /nav/addClass [post]
func AddClass(c *gin.Context) {
	var class mydb.StructNavClass

	loginToken := c.GetHeader("LoginToken")
	adminID, err := admin.GetAdminIDFromToken(loginToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.APIResponse{Code: http.StatusUnauthorized, Message: "无权限添加导航分类", Data: err.Error()})
		return
	}

	class.Admin_id = adminID

	class.Name = c.PostForm("name")
	class.Parent_id, _ = strconv.Atoi(c.PostForm("parent_id"))
	class.Sort, _ = strconv.Atoi(c.PostForm("sort"))
	class.Icon = c.PostForm("icon")
	class.Is_show, _ = strconv.ParseBool(c.PostForm("is_show"))
	class.Is_recommend, _ = strconv.ParseBool(c.PostForm("is_recommend"))
	class.Is_hot, _ = strconv.ParseBool(c.PostForm("is_hot"))
	class.Status, _ = strconv.Atoi(c.PostForm("status"))
	class.Description = c.PostForm("description")
	class.Create_time = util.GetTimestamp(10)
	class.Update_time = util.GetTimestamp(10)

	id, rowsAffected, err := class.Insert([]mydb.StructNavClass{class})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "添加导航分类失败", Data: err.Error()})
		return
	}
	log.InfoLogger.Printf("导航分类添加成功,id=%d,添加记录数:%d", id, rowsAffected)

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "导航分类添加成功"})
}

// UpdateClass 修改导航分类
// @Summary 修改导航分类
// @Description 根据导航分类ID修改导航分类信息
// @Tags nav
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param id path string true "导航分类ID"
// @Param LoginToken header string true "认证Token"
// @Param name formData string false "导航分类名称"
// @Param parent_id formData int false "父级分类id"
// @Param sort formData string false "排序"
// @Param icon formData string false "图标"
// @Param is_show formData bool false "是否显示"
// @Param is_recommend formData bool false "是否推荐"
// @Param is_hot formData bool false "是否热门"
// @Param status formData int false "状态:0=禁用,1=启用"
// @Param description formData string false "导航分类描述"
// @Success 200 {object} util.APIResponse{code=int,message=string}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=string}
// @Router /nav/updateClass/{id} [put]
func UpdateClass(c *gin.Context) {
	classID := c.Param("id")
	class, err := mydb.Tables.NavClass.Find(mydb.QueryParams{Condition: "id=" + classID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取导航分类信息失败", Data: err.Error()})
		return
	}

	loginToken := c.GetHeader("LoginToken")
	adminID, err := admin.GetAdminIDFromToken(loginToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "未授权操作", Data: err.Error()})
		return
	}

	if class.Admin_id != adminID {
		c.JSON(http.StatusUnauthorized, util.APIResponse{Code: http.StatusUnauthorized, Message: "无权限修改该导航分类", Data: "null"})
		return
	}

	if name := c.PostForm("name"); name != "" {
		class.Name = name
	}
	if parentID := c.PostForm("parent_id"); parentID != "" {
		class.Parent_id, _ = strconv.Atoi(parentID)
	}
	if sort := c.PostForm("sort"); sort != "" {
		class.Sort, _ = strconv.Atoi(sort)
	}
	if icon := c.PostForm("icon"); icon != "" {
		class.Icon = icon
	}
	if isShow := c.PostForm("is_show"); isShow != "" {
		class.Is_show, _ = strconv.ParseBool(isShow)
	}
	if isRecommend := c.PostForm("is_recommend"); isRecommend != "" {
		class.Is_recommend, _ = strconv.ParseBool(isRecommend)
	}
	if isHot := c.PostForm("is_hot"); isHot != "" {
		class.Is_hot, _ = strconv.ParseBool(isHot)
	}
	if status := c.PostForm("status"); status != "" {
		class.Status, _ = strconv.Atoi(status)
	}
	if description := c.PostForm("description"); description != "" {
		class.Description = description
	}
	class.Update_time = util.GetTimestamp(10)

	_, _, err = class.Update([]mydb.StructNavClass{class}, "id="+strconv.Itoa(class.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "修改导航分类失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "导航分类修改成功"})
}

// DeleteClass 删除导航分类
// @Summary 删除导航分类
// @Description 根据导航分类ID删除导航分类
// @Tags nav
// @Produce application/json
// @Param id path string true "导航分类ID"
// @Success 200 {object} util.APIResponse{code=int,message=string}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=string}
// @Router /nav/deleteClass/{id} [delete]
func DeleteClass(c *gin.Context) {
	classID := c.Param("id")
	class, err := mydb.Tables.NavClass.Find(mydb.QueryParams{Condition: "id=" + classID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取导航分类信息失败", Data: err.Error()})
		return
	}

	_, _, err = class.Delete("id=" + strconv.Itoa(class.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "删除导航分类失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "导航分类删除成功"})
}

// GetClassList 获取导航分类列表
// @Summary 获取导航分类列表
// @Description 获取所有导航分类的列表
// @Tags nav
// @Produce application/json
// @Success 200 {object} util.APIResponse{code=int,message=string,data=[]mydb.StructNavClass}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=string}
// @Router /nav/getClassList [get]
func GetClassList(c *gin.Context) {
	params := mydb.QueryParams{}
	classes, _, err := mydb.Tables.NavClass.Select(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取导航分类列表失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "获取导航分类列表成功", Data: classes})
}
