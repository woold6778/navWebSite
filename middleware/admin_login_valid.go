package middleware

import (
	"nav-web-site/app/api/v1/admin"
	"nav-web-site/util/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 验证 logintoken 的中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求路径
		path := c.Request.URL.Path

		// 如果请求路径是 /admin/register 或 /admin/login，则跳过验证
		if path == "/api/v1/admin/register" || path == "/api/v1/admin/login" {
			c.Next()
			return
		}

		// 获取请求头中的 logintoken
		token := c.GetHeader("LoginToken")

		// 检查 logintoken 是否存在并且有效
		if token == "" {
			log.InfoLogger.Printf("Missing logintoken: %s %s, ClientIP: %s",
				c.Request.Method, c.Request.URL.Path, c.ClientIP())
			c.AbortWithStatus(http.StatusForbidden)
			return
			/*
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing logintoken"})
				c.Abort() // 终止请求
				return
			*/
		}

		isValid, errMsg := isValidToken(c, token)
		if !isValid {
			log.InfoLogger.Printf("Invalid logintoken:%s %s %s, ClientIP: %s", errMsg,
				c.Request.Method, c.Request.URL.Path, c.ClientIP())
			c.AbortWithStatus(http.StatusForbidden)
			return
			/*
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid logintoken: " + errMsg})
				c.Abort() // 终止请求
				return
			*/
		}

		// 如果 token 有效，继续处理请求
		c.Next()
	}
}

// 用于检查 token 的有效性
func isValidToken(c *gin.Context, token string) (bool, string) {
	tokenContent, err := admin.GetTokenContent(token)
	if err != nil {
		return false, err.Error()
	}
	log.InfoLogger.Println(tokenContent)

	// 获取请求的客户端IP、User-Agent和设备指纹
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	//deviceFingerprint := util.GenerateDeviceFingerprint(c.Request)

	// 验证缓存中的客户端信息
	if tokenContent["client_ip"] != clientIP {
		return false, "Client IP mismatch"
	}
	if tokenContent["user_agent"] != userAgent {
		return false, "User-Agent mismatch"
	}
	/*
		if tokenContent["device_fingerprint"] != deviceFingerprint {
			return false, "Device fingerprint mismatch"
		}
	*/

	return true, "ok"
}
