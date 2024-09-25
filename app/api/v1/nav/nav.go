package nav

import (
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
// @Param name formData string true "导航名称"
// @Param url formData string true "导航链接"
// @Param description formData string false "导航描述"
// @Success 200 {object} util.APIResponse{code=int,message=string}
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=string}
// @Router /nav/addData [post]
func AddData(c *gin.Context) {
	var data mydb.StructNav
	data.Title = c.PostForm("name")
	data.Url = c.PostForm("url")
	data.Description = c.PostForm("description")

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
// @Param id path string true "导航信息ID"
// @Param name formData string true "导航名称"
// @Param url formData string true "导航链接"
// @Param description formData string false "导航描述"
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

	data.Title = c.PostForm("name")
	data.Url = c.PostForm("url")
	data.Description = c.PostForm("description")

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
