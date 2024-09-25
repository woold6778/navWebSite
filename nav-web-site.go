package main

import (
	"nav-web-site/app/api/upload"
	"nav-web-site/app/api/v1/admin"
	"nav-web-site/app/api/v1/nav"
	"nav-web-site/app/api/v1/news"
	"nav-web-site/config"
	"nav-web-site/middleware"
	"nav-web-site/mydb"
	"nav-web-site/util"

	_ "nav-web-site/docs" // 这里导入生成的docs文件

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Nav Web Site API
// @version 1.0
// @description This is a sample server for a nav web site.
// @host nav.fandoc.org
// @BasePath /api/v1
func main() {

	config.InitConfig() // 初始化配置

	util.InitLoggers() // 初始化日志

	// 使用 InfoLogger 和 ErrorLogger 来记录日志
	util.InfoLogger.Println("Starting the server...")
	util.ErrorLogger.Println("This is an error log message test.")

	runMode := gin.ReleaseMode
	if config.Config.Base.Debug {
		runMode = gin.DebugMode
	}
	gin.SetMode(runMode)

	//初始数据库和redis
	mydb.InitDB()
	defer mydb.Db.Close() // 确保在程序结束时关闭数据库连接

	//定义路由
	r := gin.Default()

	// Swagger 路由配置
	r.GET("/swagger/*any", SwaggerAuthMiddleware(), ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// 图片上传模块组
	uploadGroup := v1.Group("/upload")
	{
		// @Summary 上传图片
		// @Description 上传图片文件
		// @Tags upload
		// @Accept multipart/form-data
		// @Produce application/json
		// @Param file formData file true "图片文件"
		// @Success 200 {object} gin.H{"message": string, "file_path": string}
		// @Failure 400 {object} gin.H{"message": string}
		// @Router /upload/image [post]
		uploadGroup.POST("/image", upload.UploadImage)
	}

	// 管理员用户模块组
	adminGroup := v1.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware())
	{
		// @Summary 登录
		// @Description 管理员登录
		// @Tags admin
		// @Accept json
		// @Produce json
		// @Param body body admin.LoginRequest true "登录请求"
		// @Success 200 {object} admin.LoginResponse
		// @Router /admin/login [post]
		adminGroup.POST("/login", admin.Login)

		// @Summary 注册
		// @Description 管理员注册
		// @Tags admin
		// @Accept json
		// @Produce json
		// @Param body body admin.RegisterRequest true "注册请求"
		// @Success 200 {object} admin.RegisterResponse
		// @Router /admin/register [post]
		adminGroup.POST("/register", admin.Register)

		// @Summary 获取管理员列表
		// @Description 获取所有管理员的列表
		// @Tags admin
		// @Produce json
		// @Success 200 {object} []admin.User
		// @Router /admin/list [get]
		adminGroup.GET("/list", admin.GetUserList)

		// @Summary 获取管理员详情
		// @Description 根据管理员ID获取管理员详情
		// @Tags admin
		// @Produce json
		// @Param id path string true "管理员ID"
		// @Success 200 {object} admin.User
		// @Router /admin/detail/{id} [get]
		adminGroup.GET("/detail/:id", admin.GetUserDetail)

		// @Summary 修改管理员密码
		// @Description 根据管理员ID修改管理员密码
		// @Tags admin
		// @Accept json
		// @Produce json
		// @Param id path string true "管理员ID"
		// @Param body body admin.UpdatePasswordRequest true "修改密码请求"
		// @Success 200 {object} gin.H{"message": string}
		// @Router /admin/updatePassword/{id} [put]
		adminGroup.PUT("/updatePassword/:id", admin.UpdateUserPassword)

		// @Summary 编辑管理员信息
		// @Description 根据管理员ID编辑管理员信息
		// @Tags admin
		// @Accept json
		// @Produce json
		// @Param id path string true "管理员ID"
		// @Param body body admin.EditProfileRequest true "编辑信息请求"
		// @Success 200 {object} gin.H{"message": string}
		// @Router /admin/editProfile/{id} [put]
		adminGroup.PUT("/editProfile/:id", admin.EditUserProfile)

		// @Summary 删除管理员用户
		// @Description 根据管理员ID删除管理员用户
		// @Tags admin
		// @Produce json
		// @Param id path string true "管理员ID"
		// @Success 200 {object} gin.H{"message": string}
		// @Router /admin/delete/{id} [delete]
		adminGroup.DELETE("/delete/:id", admin.DeleteUser)
	}

	//导航模块路由组
	navGroup := v1.Group("/nav")
	{
		// @Summary 添加导航分类
		// @Description 添加导航分类
		// @Tags nav
		// @Accept json
		// @Produce json
		// @Param body body nav.AddClassRequest true "添加导航分类请求"
		// @Success 200 {object} nav.AddClassResponse
		// @Router /nav/addClass [post]
		navGroup.POST("/addClass", nav.AddClass) // 添加导航分类

		// @Summary 获取导航分类列表
		// @Description 获取所有导航分类的列表
		// @Tags nav
		// @Produce json
		// @Success 200 {object} []nav.Class
		// @Router /nav/getClassList [get]
		navGroup.GET("/getClassList", nav.GetClassList) // 获取导航分类列表

		// @Summary 更新导航分类
		// @Description 根据导航分类ID更新导航分类
		// @Tags nav
		// @Accept json
		// @Produce json
		// @Param body body nav.UpdateClassRequest true "更新导航分类请求"
		// @Success 200 {object} nav.UpdateClassResponse
		// @Router /nav/updateClass [put]
		navGroup.PUT("/updateClass", nav.UpdateClass) // 更新导航分类

		// @Summary 添加导航信息数据
		// @Description 添加导航信息数据
		// @Tags nav
		// @Accept json
		// @Produce json
		// @Param body body nav.AddDataRequest true "添加导航信息请求"
		// @Success 200 {object} nav.AddDataResponse
		// @Router /nav/addData [post]
		navGroup.POST("/addData", nav.AddData) // 添加导航信息数据

		// @Summary 获取导航列表
		// @Description 获取所有导航信息的列表
		// @Tags nav
		// @Produce json
		// @Success 200 {object} []nav.Data
		// @Router /nav/getList [get]
		navGroup.GET("/getList", nav.GetDataList) // 获取导航列表

		// @Summary 获取信息详情
		// @Description 根据导航信息ID获取导航信息详情
		// @Tags nav
		// @Produce json
		// @Param id path string true "导航信息ID"
		// @Success 200 {object} nav.Data
		// @Router /nav/getDetail/{id} [get]
		navGroup.GET("/getDetail", nav.GetDataDetail) // 获取信息详情

		// @Summary 更新导航数据
		// @Description 根据导航信息ID更新导航数据
		// @Tags nav
		// @Accept json
		// @Produce json
		// @Param body body nav.UpdateDataRequest true "更新导航数据请求"
		// @Success 200 {object} nav.UpdateDataResponse
		// @Router /nav/updateData [put]
		navGroup.PUT("/updateData", nav.UpdateData) // 更新导航数据
	}

	//新闻模块路由组
	newsGroup := v1.Group("/news")
	{
		// @Summary 添加新闻
		// @Description 添加新闻
		// @Tags news
		// @Accept json
		// @Produce json
		// @Param body body news.AddNewsRequest true "添加新闻请求"
		// @Success 200 {object} news.AddNewsResponse
		// @Router /news/add [post]
		newsGroup.POST("/add", news.AddNews) // 添加新闻

		// @Summary 获取新闻列表
		// @Description 获取所有新闻的列表
		// @Tags news
		// @Produce json
		// @Success 200 {object} []news.News
		// @Router /news/list [get]
		newsGroup.GET("/list", news.GetNewsList) // 获取新闻列表

		// @Summary 获取新闻详情
		// @Description 根据新闻ID获取新闻详情
		// @Tags news
		// @Produce json
		// @Param id path string true "新闻ID"
		// @Success 200 {object} news.News
		// @Router /news/detail/{id} [get]
		newsGroup.GET("/detail", news.GetNewsDetail) // 获取新闻详情

		// @Summary 更新新闻
		// @Description 根据新闻ID更新新闻内容
		// @Tags news
		// @Accept json
		// @Produce json
		// @Param body body news.UpdateNewsRequest true "更新新闻请求"
		// @Success 200 {object} news.UpdateNewsResponse
		// @Router /news/update [put]
		newsGroup.PUT("/update", news.UpdateNews) // 更新新闻

		// @Summary 删除新闻
		// @Description 根据新闻ID删除新闻
		// @Tags news
		// @Produce json
		// @Param id path string true "新闻ID"
		// @Success 200 {object} gin.H{"message": string}
		// @Router /news/delete/{id} [delete]
		newsGroup.DELETE("/delete", news.DeleteNews) // 删除新闻
	}

	// 默认路由处理
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Page not found"})
	})

	r.Run(":8080")
}

// SwaggerAuthMiddleware 是一个简单的中间件，用于保护 Swagger 文档
func SwaggerAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()
		if hasAuth && username == "admin" && password == "wolfa1" {
			c.Next()
		} else {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatus(401)
		}
	}
}
