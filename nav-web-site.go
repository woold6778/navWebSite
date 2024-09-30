package main

import (
	"nav-web-site/app/api/upload"
	"nav-web-site/app/api/v1/admin"
	"nav-web-site/app/api/v1/nav"
	"nav-web-site/app/api/v1/news"
	"nav-web-site/app/webcrawler"
	"nav-web-site/config"
	"nav-web-site/middleware"
	"nav-web-site/mydb"
	"nav-web-site/util/log"
	"net/http"
	"os"
	"time"

	_ "nav-web-site/docs" // 这里导入生成的docs文件

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Nav Web Site API
// @version 1.0
// @description This is a sample server for a nav web site.
// @host nav.fandoc.org
// @BasePath /api/v1
func main() {

	log.InitLoggers() // 初始化日志

	config.InitConfig() // 初始化配置

	// 使用 InfoLogger 和 ErrorLogger 来记录日志
	log.InfoLogger.Println("Starting the server...")
	log.ErrorLogger.Println("This is an error log message test.")

	runMode := gin.ReleaseMode
	if config.Config.Base.Debug {
		runMode = gin.DebugMode
	}
	gin.SetMode(runMode)

	//初始数据库和redis
	mydb.InitDB()
	defer mydb.Db.Close() // 确保在程序结束时关闭数据库连接

	// 启动定时任务检测 Goroutine
	go startScheduledTaskChecker()

	//定义路由
	r := gin.Default()

	// 添加请求日志中间件
	r.Use(RequestLoggerMiddleware())

	// 添加CORS中间件
	// 从配置文件中读取 AllowOrigins
	allowOrigins := config.Config.Base.AllowOrigins
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins, // 允许的前端域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	log.InfoLogger.Printf("CORS configuration: AllowOrigins: %v, AllowMethods: %v, AllowHeaders: %v, ExposeHeaders: %v, AllowCredentials: %v",
		allowOrigins, []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "Accept"}, []string{"Content-Length"}, true)
	// 处理OPTIONS预检请求
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With, Accept")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Status(http.StatusOK)
	})
	// Swagger 路由配置
	r.GET("/swagger/*any", SwaggerAuthMiddleware(), ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 图片获取模块组
	imageGroup := r.Group("/images")
	{
		imageGroup.GET("/:hash", upload.GetImageByHash)
	}

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
		navGroup.PUT("/updateClass/:id", nav.UpdateClass) // 更新导航分类

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
		navGroup.GET("/getDetail/:id", nav.GetDataDetail) // 获取信息详情

		// @Summary 更新导航数据
		// @Description 根据导航信息ID更新导航数据
		// @Tags nav
		// @Accept json
		// @Produce json
		// @Param body body nav.UpdateDataRequest true "更新导航数据请求"
		// @Success 200 {object} nav.UpdateDataResponse
		// @Router /nav/updateData [put]
		navGroup.PUT("/updateData/:id", nav.UpdateData) // 更新导航数据
	}

	//新闻模块路由组
	newsGroup := v1.Group("/news")
	{
		// @Summary 添加新闻分类
		// @Description 添加新闻分类
		// @Tags news
		// @Accept json
		// @Produce json
		// @Param body body news.AddClassRequest true "添加新闻分类请求"
		// @Success 200 {object} news.AddClassResponse
		// @Router /news/addClass [post]
		newsGroup.POST("/addClass", news.AddClass) // 添加新闻分类

		// @Summary 编辑新闻分类
		// @Description 根据新闻分类ID编辑新闻分类
		// @Tags news
		// @Accept json
		// @Produce json
		// @Param body body news.UpdateClassRequest true "编辑新闻分类请求"
		// @Success 200 {object} news.UpdateClassResponse
		// @Router /news/updateClass/{id} [put]
		newsGroup.PUT("/updateClass/:id", news.UpdateClass) // 编辑新闻分类

		// @Summary 删除新闻分类
		// @Description 根据新闻分类ID删除新闻分类
		// @Tags news
		// @Produce json
		// @Param id path string true "新闻分类ID"
		// @Success 200 {object} gin.H{"message": string}
		// @Router /news/deleteClass/{id} [delete]
		newsGroup.DELETE("/deleteClass/:id", news.DeleteClass) // 删除新闻分类

		// @Summary 获取新闻分类列表
		// @Description 获取所有新闻分类的列表
		// @Tags news
		// @Produce json
		// @Success 200 {object} []news.Class
		// @Router /news/getClassList [get]
		newsGroup.GET("/getClassList", news.GetClassList) // 获取新闻分类列表

		// @Summary 获取新闻分类详情
		// @Description 根据新闻分类ID获取新闻分类详情
		// @Tags news
		// @Produce json
		// @Param id path string true "新闻分类ID"
		// @Success 200 {object} news.GetClassDetail
		// @Router /news/getClassDetail/{id} [get]
		newsGroup.GET("/getClassDetail/:id", news.GetClassDetail) // 获取新闻分类详情
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
		newsGroup.GET("/detail/:id", news.GetNewsDetail) // 获取新闻详情

		// @Summary 更新新闻
		// @Description 根据新闻ID更新新闻内容
		// @Tags news
		// @Accept json
		// @Produce json
		// @Param body body news.UpdateNewsRequest true "更新新闻请求"
		// @Success 200 {object} news.UpdateNewsResponse
		// @Router /news/update [put]
		newsGroup.PUT("/update/:id", news.UpdateNews) // 更新新闻

		// @Summary 删除新闻
		// @Description 根据新闻ID删除新闻
		// @Tags news
		// @Produce json
		// @Param id path string true "新闻ID"
		// @Success 200 {object} gin.H{"message": string}
		// @Router /news/delete/{id} [delete]
		newsGroup.DELETE("/delete/:id", news.DeleteNews) // 删除新闻
	}
	// 如果上面的路由都没匹配到，就到指定目录（如：/www/wwwroot/nav/）的对应url路径下查找文件，如果有就返回文件内容，否则就报404
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		if path == "/" {
			path = "/index.html"
		}

		filePath := config.Config.Base.BackendFilePath + path
		log.InfoLogger.Printf("Requested file path: %s", filePath)
		if _, err := os.Stat(filePath); err == nil {
			http.ServeFile(c.Writer, c.Request, filePath)

		} else {

			c.JSON(404, gin.H{"message": "Page not found."})

		}
	})

	// 默认路由处理
	/*
		r.NoRoute(func(c *gin.Context) {
			c.JSON(404, gin.H{"message": "Page not found"})
		})
	*/

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

// RequestLoggerMiddleware 是一个中间件，用于记录每一次用户的请求
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 记录请求日志
		duration := time.Since(startTime)
		userAgent := c.Request.UserAgent()
		log.InfoLogger.Printf("Request: %s %s, Status: %d, Duration: %v, ClientIP: %s, User-Agent: %s",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration, c.ClientIP(), userAgent)
	}
}

// 开始计划任务
func startScheduledTaskChecker() {
	c := cron.New()
	for _, task := range config.Config.Tasks {
		task := task                                                                           // 避免闭包问题
		log.InfoLogger.Printf("Scheduling task %s with schedule %s", task.Type, task.Schedule) // 添加日志
		_, err := c.AddFunc(task.Schedule, func() {
			log.InfoLogger.Printf("Executing scheduled task %s", task.Type) // 添加日志
			executeTask(task.Type)
		})
		if err != nil {
			log.ErrorLogger.Printf("Error scheduling task %s: %v", task.Type, err)
		}
	}
	c.Start()
	log.InfoLogger.Println("Scheduled tasks started") // 添加日志
}

// 检测redis是否有任务
func checkAndExecuteTasks() {
	// 从redis查看有没有待执行的任务
	keys, err := mydb.RedisClient.Keys(mydb.Ctx, "scheduled_task:*").Result()
	if err != nil {
		log.ErrorLogger.Println("Error fetching scheduled tasks:", err)
		return
	}

	for _, key := range keys {
		go executeTask(key)
	}
}

// 按任务类型执行任务
func executeTask(taskType string) {
	log.InfoLogger.Printf("Starting task execution for %s", taskType)
	defer func() {
		if r := recover(); r != nil {
			log.ErrorLogger.Printf("Task %s panicked: %v", taskType, r)
		}
	}()
	switch taskType {
	case "news163":
		// 获取163的头条新闻
		webcrawler.FetchAndStoreNews163()
		log.InfoLogger.Println("Executing task news163")
		// ... 任务代码 ...
	case "type2":
		// 执行类型2的任务
		log.InfoLogger.Println("Executing task type2")
		// ... 任务代码 ...
	default:
		log.ErrorLogger.Println("Unknown task type:", taskType)
	}
}
