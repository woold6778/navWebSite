package news

import (
	"nav-web-site/app/api/v1/admin"
	"nav-web-site/mydb"
	"nav-web-site/util"
	"nav-web-site/util/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddNewsClass 添加新闻分类
// @Summary 添加新闻分类
// @Description 添加新的新闻分类
// @Tags news
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param LoginToken header string true "认证Token"
// @Param name formData string true "新闻分类名称"
// @Param parent_id formData int false "父级分类id"
// @Param sort formData int false "排序"
// @Param icon formData string false "图标"
// @Param description formData string false "描述"
// @Param keywords formData string false "关键词"
// @Param is_show formData bool false "是否显示"
// @Param is_recommend formData bool false "是否推荐"
// @Param is_hot formData bool false "是否热门"
// @Param status formData int false "状态:0=禁用,1=启用"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=interface{}} "新闻分类添加成功"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "添加新闻分类失败"
// @Router /news/addClass [post]
func AddClass(c *gin.Context) {
	var class mydb.StructNewsClass

	loginToken := c.GetHeader("LoginToken")
	adminID, err := admin.GetAdminIDFromToken(loginToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.APIResponse{Code: http.StatusUnauthorized, Message: "无权限添加新闻分类", Data: err.Error()})
		return
	}

	class.Admin_id = adminID

	class.Name = c.PostForm("name")
	class.Parent_id, _ = strconv.Atoi(c.PostForm("parent_id"))
	class.Sort, _ = strconv.Atoi(c.PostForm("sort"))
	class.Icon = c.PostForm("icon")
	class.Description = c.PostForm("description")
	class.Keywords = c.PostForm("keywords")
	class.Is_show, _ = strconv.ParseBool(c.PostForm("is_show"))
	class.Is_recommend, _ = strconv.ParseBool(c.PostForm("is_recommend"))
	class.Is_hot, _ = strconv.ParseBool(c.PostForm("is_hot"))
	class.Status, _ = strconv.Atoi(c.PostForm("status"))
	class.Create_time = util.GetTimestamp(10)
	class.Update_time = util.GetTimestamp(10)

	id, rowsAffected, err := class.Insert([]mydb.StructNewsClass{class})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "添加新闻分类失败", Data: err.Error()})
		return
	}
	log.InfoLogger.Printf("新闻分类添加成功,id=%d,添加记录数:%d", id, rowsAffected)

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "新闻分类添加成功"})
}

// UpdateNewsClass 修改新闻分类
// @Summary 修改新闻分类
// @Description 根据新闻分类ID修改新闻分类信息
// @Tags news
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param LoginToken header string true "认证Token"
// @Param id path string true "新闻分类ID"
// @Param name formData string true "新闻分类名称"
// @Param parent_id formData int false "父级分类id"
// @Param sort formData int false "排序"
// @Param icon formData string false "图标"
// @Param description formData string false "描述"
// @Param keywords formData string false "关键词"
// @Param is_show formData bool false "是否显示"
// @Param is_recommend formData bool false "是否推荐"
// @Param is_hot formData bool false "是否热门"
// @Param status formData int false "状态:0=禁用,1=启用"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=interface{}} "新闻分类修改成功"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "修改新闻分类失败"
// @Router /news/updateClass/{id} [put]
func UpdateClass(c *gin.Context) {
	classID := c.Param("id")
	class, err := mydb.Tables.NewsClass.Find(mydb.QueryParams{Condition: "id=" + classID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取新闻分类信息失败", Data: err.Error()})
		return
	}

	loginToken := c.GetHeader("LoginToken")
	adminID, err := admin.GetAdminIDFromToken(loginToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "未授权操作", Data: err.Error()})
		return
	}

	if class.Admin_id != adminID {
		c.JSON(http.StatusUnauthorized, util.APIResponse{Code: http.StatusUnauthorized, Message: "无权限修改该分类", Data: "null"})
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
	if description := c.PostForm("description"); description != "" {
		class.Description = description
	}
	if keywords := c.PostForm("keywords"); keywords != "" {
		class.Keywords = keywords
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
	class.Update_time = util.GetTimestamp(10)

	_, _, err = class.Update([]mydb.StructNewsClass{class}, "id="+strconv.Itoa(class.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "修改新闻分类失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "新闻分类修改成功"})
}

// GetNewsClassList 获取新闻分类列表
// @Summary 获取新闻分类列表
// @Description 获取所有新闻分类的列表
// @Tags news
// @Produce application/json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=[]mydb.StructNewsClass} "获取新闻分类列表成功"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "获取新闻分类列表失败"
// @Router /news/getClassList [get]
func GetClassList(c *gin.Context) {
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

	classes, _, err := mydb.Tables.NewsClass.Select(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取新闻分类列表失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "获取新闻分类列表成功", Data: classes})
}

// GetNewsClassDetail 获取新闻分类详情
// @Summary 获取新闻分类详情
// @Description 根据新闻分类ID获取新闻分类详情
// @Tags news
// @Produce application/json
// @Param id path string true "新闻分类ID"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=mydb.StructNewsClass} "获取新闻分类详情成功"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "获取新闻分类详情失败"
// @Router /news/getClassDetail/{id} [get]
func GetClassDetail(c *gin.Context) {
	classID := c.Param("id")
	class, err := mydb.Tables.NewsClass.Find(mydb.QueryParams{Condition: "id=" + classID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取新闻分类详情失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "获取新闻分类详情成功", Data: class})
}

// DeleteNewsClass 删除新闻分类
// @Summary 删除新闻分类
// @Description 根据新闻分类ID删除新闻分类
// @Tags news
// @Produce application/json
// @Param id path string true "新闻分类ID"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=interface{}} "新闻分类删除成功"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "删除新闻分类失败"
// @Router /news/deleteClass/{id} [delete]
func DeleteClass(c *gin.Context) {
	classID := c.Param("id")
	class, err := mydb.Tables.NewsClass.Find(mydb.QueryParams{Condition: "id=" + classID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取新闻分类信息失败", Data: err.Error()})
		return
	}

	_, _, err = class.Delete("id=" + strconv.Itoa(class.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "删除新闻分类失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "新闻分类删除成功"})
}
