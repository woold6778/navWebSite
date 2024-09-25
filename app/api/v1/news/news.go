package news

import (
	"encoding/json"
	"net/http"
	"strconv"

	"nav-web-site/mydb"
	"nav-web-site/util"

	"github.com/gin-gonic/gin"
)

// AddNews 添加新闻
// @Summary 添加新闻
// @Description 添加一条新的新闻记录
// @Tags news
// @Accept application/json
// @Produce application/json
// @Param news body mydb.StructNews true "新闻内容"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=interface{}} "新闻添加成功"
// @Failure 400 {object} util.APIResponse{code=int,message=string,data=interface{}} "请求参数错误"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "添加新闻失败"
// @Router /news/add [post]
func AddNews(c *gin.Context) {
	var news mydb.StructNews
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse{Code: http.StatusBadRequest, Message: "请求参数错误", Data: err.Error()})
		return
	}

	_, ids, err := mydb.Tables.News.Insert([]mydb.StructNews{news})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "添加新闻失败", Data: err.Error()})
		return
	}

	idsJSON, err := json.Marshal(ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "转换新闻ID列表失败", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "新闻添加成功", Data: string(idsJSON)})
}

// UpdateNews 编辑新闻
// @Summary 编辑新闻
// @Description 根据新闻ID编辑新闻内容
// @Tags news
// @Accept application/json
// @Produce application/json
// @Param news body mydb.StructNews true "新闻内容"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=interface{}} "新闻编辑成功"
// @Failure 400 {object} util.APIResponse{code=int,message=string,data=interface{}} "请求参数错误"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "编辑新闻失败"
// @Router /news/update [put]
func UpdateNews(c *gin.Context) {
	var news mydb.StructNews
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse{Code: http.StatusBadRequest, Message: "请求参数错误", Data: err.Error()})
		return
	}

	condition := "id=" + strconv.Itoa(news.ID)
	_, ids, err := mydb.Tables.News.Update([]mydb.StructNews{news}, condition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "编辑新闻失败", Data: err.Error()})
		return
	}

	idsJSON, err := json.Marshal(ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "转换新闻ID列表失败", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "新闻编辑成功", Data: string(idsJSON)})
}

// GetNewsList 获取新闻列表
// @Summary 获取新闻列表
// @Description 获取新闻列表，支持按分类ID筛选
// @Tags news
// @Produce application/json
// @Param class_id query string false "新闻分类ID"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=[]mydb.StructNews} "获取新闻列表成功"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "获取新闻列表失败"
// @Router /news/list [get]
func GetNewsList(c *gin.Context) {
	classID := c.Query("class_id")
	condition := ""
	if classID != "" {
		condition = "class_id=" + classID
	}

	params := mydb.QueryParams{
		Condition: condition,
		Limit:     100,
		Page:      0,
		PageSize:  10,
	}
	newsList, _, err := mydb.Tables.News.Select(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取新闻列表失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "获取新闻列表成功", Data: newsList})
}

// GetNewsDetail 获取新闻详情
// @Summary 获取新闻详情
// @Description 根据新闻ID获取新闻详情
// @Tags news
// @Produce application/json
// @Param id path string true "新闻ID"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=mydb.StructNews} "获取新闻详情成功"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "获取新闻详情失败"
// @Router /news/detail/{id} [get]
func GetNewsDetail(c *gin.Context) {
	newsID := c.Param("id")
	news, err := mydb.Tables.News.Find("id=" + newsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取新闻详情失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "获取新闻详情成功", Data: news})
}

// DeleteNews 删除新闻
// @Summary 删除新闻
// @Description 根据新闻ID删除新闻
// @Tags news
// @Produce application/json
// @Param id path string true "新闻ID"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=interface{}} "新闻删除成功"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "删除新闻失败"
// @Router /news/delete/{id} [delete]
func DeleteNews(c *gin.Context) {
	newsID := c.Param("id")
	news, err := mydb.Tables.News.Find("id=" + newsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取新闻信息失败", Data: err.Error()})
		return
	}

	_, _, err = news.Delete("id=" + strconv.Itoa(news.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "删除新闻失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "新闻删除成功"})
}
