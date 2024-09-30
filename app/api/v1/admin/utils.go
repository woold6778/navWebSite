package admin

import (
	"fmt"
	"nav-web-site/mydb"
	"nav-web-site/util"
	"nav-web-site/util/log"
)

// 通过login_token获取管理员ID
func GetAdminIDFromToken(loginToken string) (int, error) {
	cacheKey := "admin_login_token_" + loginToken
	log.InfoLogger.Println("缓存Key:", cacheKey)
	cacheValue, found := util.C.Get(cacheKey)
	if !found {
		return 0, fmt.Errorf("认证失败1")
	}

	adminData, ok := cacheValue.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("获取管理员信息失败")
	}

	admin, ok := adminData["admin"].(mydb.StructAdmin)
	if !ok {
		return 0, fmt.Errorf("解析管理员信息失败")
	}

	return admin.ID, nil
}

// GetTokenContent 获取token的内容
func GetTokenContent(loginToken string) (map[string]interface{}, error) {
	cacheKey := "admin_login_token_" + loginToken
	log.InfoLogger.Println("缓存Key:", cacheKey)
	cacheValue, found := util.C.Get(cacheKey)
	if !found {
		return nil, fmt.Errorf("认证失败2")
	}

	tokenContent, ok := cacheValue.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("获取token内容失败")
	}

	return tokenContent, nil
}
