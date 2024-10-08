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
    go get -u github.com/PuerkitoBio/goquery
    go get github.com/robfig/cron/v3
    go get github.com/gin-contrib/cors

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

查找占用8080端口的程序
    sudo netstat -tulpn | grep :8080    

创建navwebsite.service文件
    sudo nano /etc/systemd/system/navwebsite.service
    sudo vim /etc/systemd/system/navwebsite.service

    文件内容如下：
[Unit]
Description=navwebsite Service
After=network.target    

[Service]
User=navwebsiteuser
Group=navwebsiteuser
ExecStart=/www/wwwroot/navwebsite/navwebsite
WorkingDirectory=/www/wwwroot/navwebsite
Restart=always
RestartSec=5
StandardOutput=syslog  

创建用户navwebsiteuser和组navwebsiteuser，并添加相应的权限
    sudo useradd navwebsiteuser
    sudo groupadd navwebsiteuser
    sudo usermod -a -G navwebsiteuser navwebsiteuser
    sudo chown -R navwebsiteuser:navwebsiteuser /www/wwwroot/go/navwebsite/ 

启动和启用服务    
    sudo systemctl start navwebsite.service
    sudo systemctl enable navwebsite.service
检查服务状态
        sudo systemctl status navwebsite.service
查看日志
        sudo journalctl -u navwebsite.service
重新加载 systemd 配置并重启服务
        sudo systemctl daemon-reload
        sudo systemctl restart navwebsite.service

