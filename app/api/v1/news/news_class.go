package news

import (
	"nav-web-site/mydb"
	"nav-web-site/util"
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
// @Param name formData string true "新闻分类名称"
// @Param description formData string false "新闻分类描述"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=interface{}} "新闻分类添加成功"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "添加新闻分类失败"
// @Router /news/addClass [post]
func AddNewsClass(c *gin.Context) {
	var class mydb.StructNewsClass
	class.Name = c.PostForm("name")
	class.Description = c.PostForm("description")

	id, rowsAffected, err := class.Insert([]mydb.StructNewsClass{class})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "添加新闻分类失败", Data: err.Error()})
		return
	}
	util.InfoLogger.Printf("新闻分类添加成功,id=%d,添加记录数:%d", id, rowsAffected)

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "新闻分类添加成功"})
}

// UpdateNewsClass 修改新闻分类
// @Summary 修改新闻分类
// @Description 根据新闻分类ID修改新闻分类信息
// @Tags news
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param id path string true "新闻分类ID"
// @Param name formData string true "新闻分类名称"
// @Param description formData string false "新闻分类描述"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=interface{}} "新闻分类修改成功"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "修改新闻分类失败"
// @Router /news/updateClass/{id} [put]
func UpdateNewsClass(c *gin.Context) {
	classID := c.Param("id")
	class, err := mydb.Tables.NewsClass.Find(mydb.QueryParams{Condition: "id=" + classID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取新闻分类信息失败", Data: err.Error()})
		return
	}

	class.Name = c.PostForm("name")
	class.Description = c.PostForm("description")

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
// @Success 200 {object} util.APIResponse{code=int,message=string,data=[]mydb.StructNewsClass} "获取新闻分类列表成功"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "获取新闻分类列表失败"
// @Router /news/getClassList [get]
func GetNewsClassList(c *gin.Context) {
	classes, _, err := mydb.Tables.NewsClass.Select(mydb.QueryParams{})
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
func GetNewsClassDetail(c *gin.Context) {
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
func DeleteNewsClass(c *gin.Context) {
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
