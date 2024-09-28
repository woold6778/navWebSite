package news

import (
	"net/http"
	"strconv"

	"nav-web-site/app/api/v1/admin"
	"nav-web-site/mydb"
	"nav-web-site/util"

	"github.com/gin-gonic/gin"
)

// AddNews 添加新闻
// @Summary 添加新闻
// @Description 添加一条新的新闻记录
// @Tags news
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param LoginToken header string true "认证Token"
// @Param class_id formData int true "分类ID"
// @Param title formData string true "标题"
// @Param subtitle formData string false "副标题"
// @Param url formData string false "链接"
// @Param description formData string true "描述"
// @Param icon formData string false "图标"
// @Param keywords formData string false "关键词"
// @Param sort formData int false "排序"
// @Param is_show formData bool false "是否显示"
// @Param status formData int false "状态"
// @Param author formData string false "作者"
// @Param source formData string false "来源"
// @Param language formData string false "语言"
// @Param is_hot formData bool false "是否热门"
// @Param is_headline formData bool false "是否头条"
// @Param is_recommended formData bool false "是否推荐"
// @Param content formData string false "内容"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=interface{}} "新闻添加成功"
// @Failure 400 {object} util.APIResponse{code=int,message=string,data=interface{}} "请求参数错误"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "添加新闻失败"
// @Router /news/add [post]
func AddNews(c *gin.Context) {
	var news mydb.StructNews

	loginToken := c.GetHeader("LoginToken")
	adminID, err := admin.GetAdminIDFromToken(loginToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.APIResponse{Code: http.StatusUnauthorized, Message: "无权限添加导航", Data: err.Error()})
		return
	}

	news.Admin_id = adminID

	news.Class_id, _ = strconv.Atoi(c.PostForm("class_id"))
	news.Title = c.PostForm("title")
	news.Subtitle = c.PostForm("subtitle")
	news.Url = c.PostForm("url")
	news.Description = c.PostForm("description")
	news.Icon = c.PostForm("icon")
	news.Keywords = c.PostForm("keywords")
	news.Sort, _ = strconv.Atoi(c.PostForm("sort"))
	news.Is_show, _ = strconv.ParseBool(c.PostForm("is_show"))
	news.Status, _ = strconv.Atoi(c.PostForm("status"))
	news.Author = c.PostForm("author")
	news.Source = c.PostForm("source")
	news.Language = c.PostForm("language")
	news.Is_hot, _ = strconv.ParseBool(c.PostForm("is_hot"))
	news.Is_headline, _ = strconv.ParseBool(c.PostForm("is_headline"))
	news.Is_recommended, _ = strconv.ParseBool(c.PostForm("is_recommended"))
	news.Content = c.PostForm("content")
	news.Create_time = util.GetTimestamp(10)

	id, rowsAffected, err := mydb.Tables.News.Insert([]mydb.StructNews{news})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "添加新闻失败", Data: err.Error()})
		return
	}
	util.InfoLogger.Printf("新闻添加成功,标题=“%s”,id=%d,添加记录数:%d", news.Title, id, rowsAffected)

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "新闻添加成功", Data: "ok"})
}

// UpdateNews 编辑新闻
// @Summary 编辑新闻
// @Description 根据新闻ID编辑新闻内容
// @Tags news
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param LoginToken header string true "认证Token"
// @Param id path int true "新闻ID"
// @Param class_id formData int true "分类ID"
// @Param title formData string true "标题"
// @Param subtitle formData string false "副标题"
// @Param url formData string false "链接地址"
// @Param description formData string false "描述"
// @Param icon formData string false "图标"
// @Param keywords formData string false "关键词"
// @Param sort formData int false "排序"
// @Param is_show formData bool false "是否显示"
// @Param status formData int false "状态"
// @Param author formData string false "作者"
// @Param source formData string false "来源"
// @Param language formData string false "语言"
// @Param is_hot formData bool false "是否热门"
// @Param is_headline formData bool false "是否头条"
// @Param is_recommended formData bool false "是否推荐"
// @Param content formData string false "内容"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=interface{}} "新闻编辑成功"
// @Failure 400 {object} util.APIResponse{code=int,message=string,data=interface{}} "请求参数错误"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "编辑新闻失败"
// @Router /news/update/{id} [put]
func UpdateNews(c *gin.Context) {
	var news mydb.StructNews
	dataID := c.Param("id")
	news, err := (&mydb.StructNews{}).Find("id=" + dataID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "获取信息失败", Data: err.Error()})
		return
	}

	loginToken := c.GetHeader("LoginToken")
	adminID, err := admin.GetAdminIDFromToken(loginToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "未授权操作", Data: err.Error()})
		return
	}

	if news.Admin_id != adminID {
		c.JSON(http.StatusUnauthorized, util.APIResponse{Code: http.StatusUnauthorized, Message: "无权限修改", Data: "null"})
		return
	}
	if c.PostForm("class_id") != "" {
		news.Class_id, _ = strconv.Atoi(c.PostForm("class_id"))
	}
	if c.PostForm("title") != "" {
		news.Title = c.PostForm("title")
	}
	if c.PostForm("subtitle") != "" {
		news.Subtitle = c.PostForm("subtitle")
	}
	if c.PostForm("url") != "" {
		news.Url = c.PostForm("url")
	}
	if c.PostForm("description") != "" {
		news.Description = c.PostForm("description")
	}
	if c.PostForm("icon") != "" {
		news.Icon = c.PostForm("icon")
	}
	if c.PostForm("keywords") != "" {
		news.Keywords = c.PostForm("keywords")
	}
	if c.PostForm("sort") != "" {
		news.Sort, _ = strconv.Atoi(c.PostForm("sort"))
	}
	if c.PostForm("is_show") != "" {
		news.Is_show, _ = strconv.ParseBool(c.PostForm("is_show"))
	}
	if c.PostForm("status") != "" {
		news.Status, _ = strconv.Atoi(c.PostForm("status"))
	}
	if c.PostForm("author") != "" {
		news.Author = c.PostForm("author")
	}
	if c.PostForm("source") != "" {
		news.Source = c.PostForm("source")
	}
	if c.PostForm("language") != "" {
		news.Language = c.PostForm("language")
	}
	if c.PostForm("is_hot") != "" {
		news.Is_hot, _ = strconv.ParseBool(c.PostForm("is_hot"))
	}
	if c.PostForm("is_headline") != "" {
		news.Is_headline, _ = strconv.ParseBool(c.PostForm("is_headline"))
	}
	if c.PostForm("is_recommended") != "" {
		news.Is_recommended, _ = strconv.ParseBool(c.PostForm("is_recommended"))
	}
	if c.PostForm("content") != "" {
		news.Content = c.PostForm("content")
	}
	_, _, err = news.Update([]mydb.StructNews{news}, "id="+strconv.Itoa(news.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "修改信息失败", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "新闻编辑成功", Data: "ok"})
}

// GetNewsList 获取新闻列表
// @Summary 获取新闻列表
// @Description 获取新闻列表，支持按分类ID筛选
// @Tags news
// @Produce application/json
// @Param class_id query string false "新闻分类ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=[]mydb.StructNews} "获取新闻列表成功"
// @Failure 500 {object} util.APIResponse{code=int,message=string,data=interface{}} "获取新闻列表失败"
// @Router /news/list [get]
func GetNewsList(c *gin.Context) {
	classID := c.Query("class_id")
	condition := ""
	if classID != "" {
		condition = "class_id=" + classID
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
		Condition: condition,
		Page:      page,
		PageSize:  pageSize,
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
	util.InfoLogger.Println("获取的新闻详情:", news)

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
