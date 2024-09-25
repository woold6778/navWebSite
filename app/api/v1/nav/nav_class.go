package nav

import (
	"nav-web-site/mydb"
	"nav-web-site/util"
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
// @Param name formData string true "导航分类名称"
// @Param description formData string false "导航分类描述"
// @Success 200 {object} util.APIResponse{code=int,message=string}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=string}
// @Router /nav/addClass [post]
func AddClass(c *gin.Context) {
	var class mydb.StructNavClass
	class.Name = c.PostForm("name")
	class.Description = c.PostForm("description")

	id, rowsAffected, err := class.Insert([]mydb.StructNavClass{class})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "添加导航分类失败", Data: err.Error()})
		return
	}
	util.InfoLogger.Printf("导航分类添加成功,id=%d,添加记录数:%d", id, rowsAffected)

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "导航分类添加成功"})
}

// UpdateClass 修改导航分类
// @Summary 修改导航分类
// @Description 根据导航分类ID修改导航分类信息
// @Tags nav
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param id path string true "导航分类ID"
// @Param name formData string true "导航分类名称"
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

	class.Name = c.PostForm("name")
	class.Description = c.PostForm("description")

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
