package middleware

import (
	"nav-web-site/util"
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
		token := c.GetHeader("logintoken")

		// 检查 logintoken 是否存在并且有效
		if token == "" || !isValidToken(c, token) {
			// 如果 token 无效，返回 401 未授权
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing logintoken"})
			c.Abort() // 终止请求
			return
		}

		// 如果 token 有效，继续处理请求
		c.Next()
	}
}

// 用于检查 token 的有效性
func isValidToken(c *gin.Context, token string) bool {
	//在缓存在查找是否有key="admin_login_token_" + token的内容
	cacheKey := "admin_login_token_" + token

	// 获取缓存的值
	cacheValue, found := util.C.Get(cacheKey)
	if !found {
		return false
	}

	// 断言缓存值的类型
	cacheData, ok := cacheValue.(map[string]interface{})
	if !ok {
		return false
	}

	// 获取请求的客户端IP、User-Agent和设备指纹
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	deviceFingerprint := util.GenerateDeviceFingerprint(c.Request)

	// 验证缓存中的客户端信息
	if cacheData["client_ip"] != clientIP || cacheData["user_agent"] != userAgent || cacheData["device_fingerprint"] != deviceFingerprint {
		return false
	}

	return found
}
