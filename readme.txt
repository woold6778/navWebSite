初始化项目
    go mod init nav-web-site    
安装依赖包
    go get github.com/gin-gonic/gin
    go get github.com/spf13/viper
    go get github.com/go-redis/redis/v8
    go get -u github.com/go-sql-driver/mysql
    go get github.com/patrickmn/go-cache
    go get -u github.com/swaggo/swag/cmd/swag
    go get -u github.com/swaggo/gin-swagger
    go get -u github.com/swaggo/files
    go get -u github.com/swaggo/swag

自动生成api文档
    安装 swag 工具
         go install github.com/swaggo/swag/cmd/swag@latest
    确保 GOPATH/bin 已经添加到你的系统 PATH 环境变量中。你可以通过以下命令来添加（以 Windows 为例）：
        set PATH=%PATH%;%GOPATH%\bin
    生成 docs 包：
         swag init -g nav-web-site.go

    自动生成api文档    
        swag init -g nav-web-site.go
        启动项目后，访问 http://localhost:8080/swagger/index.html 即可查看自动生成的 API 文档。

检查 Go 的环境配置
    go env GOOS
    go env GOARCH
临时设置环境为windows
    set GOOS=windows
    set GOARCH=amd64
临时设置环境为linux
    set CGO_ENABLED=0         
    set GOOS=linux
    set GOARCH=amd64

