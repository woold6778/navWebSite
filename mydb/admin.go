package mydb

import (
	"fmt"
	"nav-web-site/config"
	"nav-web-site/util"
)

// StructAdmin 定义 admin 结构体
type StructAdmin struct {
	ID            int    `db:"id"`
	Username      string `db:"username"`
	Password      string `db:"password"`
	Salt          string `db:"salt"`
	Email         string `db:"email"`
	PhoneNumber   string `db:"phone_number"`
	Status        int    `db:"status"`
	CreateTime    int64  `db:"create_time"`
	UpdateTime    int64  `db:"update_time"`
	LastLoginTime int64  `db:"last_login_time"`
	Role          string `db:"role"`
	Avatar        string `db:"avatar"` // 用户头像
}

// GetTableName 获取表名
func (s *StructAdmin) GetTableName() string {
	return "admin"
}

// GetRequiredFields 获取必填字段
func (s *StructAdmin) GetRequiredFields() []string {
	return []string{"Username", "Password"}
}

// 获取插入数据时查重的字段
func (s StructAdmin) GetUniqueFields() []string {
	return []string{"Username"}
}

// Find 方法查询 admin 表的第一条数据
func (s *StructAdmin) Find(params QueryParams) (StructAdmin, error) {
	var item StructAdmin
	params.Limit = 1 // 设置查询限制为1条
	results, err := GenericSelect(Db, s.GetTableName(), params, config.Config.MySQL.TablePrefix, "")
	if err != nil {
		return item, util.WrapError(err, "Query failed(find):")
	}

	if len(results) > 0 {
		util.InfoLogger.Println("Processing result:", results[0])
		item, err = s.mapResultToStructItem(results[0])
		if err != nil {
			return item, util.WrapError(err, "将结果映射到StructAdmin时出错:")
		}
	} else {
		return item, util.WrapError(fmt.Errorf("EmptyData"), "")
	}
	return item, nil
}

// Select 方法查询 admin 表的数据
func (s *StructAdmin) Select(params QueryParams) ([]StructAdmin, int, error) {
	var list []StructAdmin
	results, err := GenericSelect(Db, s.GetTableName(), params, config.Config.MySQL.TablePrefix, "")
	if err != nil {
		return list, 400, util.WrapError(err, "Query failed(select):")
	}

	if len(results) > 0 {
		for _, result := range results {
			util.InfoLogger.Println("Processing result:", result)

			item, err := s.mapResultToStructItem(result)
			if err != nil {
				util.ErrorLogger.Println("将结果映射到StructAdmin时出错:", err)
				continue
			}

			list = append(list, item)
			util.InfoLogger.Println("Appended to list:", list[len(list)-1])
		}
	} else {
		return list, 200, util.WrapError(fmt.Errorf("EmptyData"), "")
	}
	return list, 0, nil
}

// Insert 插入新记录到 admin 表
func (s *StructAdmin) Insert() (int64, error) {
	requiredFields := s.GetRequiredFields()
	insertedCount, insertedIDs, err := GenericInsert(s.GetTableName(), []StructAdmin{*s}, requiredFields, config.Config.MySQL.TablePrefix, "")
	if err != nil {
		return 0, util.WrapError(err, "插入记录失败:")
	}
	if insertedCount == 0 {
		return 0, util.WrapError(fmt.Errorf("没有记录被插入"), "")
	}
	return insertedIDs[0], nil
}

// Update 更新 admin 表的记录
func (s *StructAdmin) Update(condition string) (int64, error) {
	updatedCount, _, err := GenericUpdate(s.GetTableName(), []StructAdmin{*s}, condition, config.Config.MySQL.TablePrefix, "")
	if err != nil {
		return 0, util.WrapError(err, "更新记录失败:")
	}
	if updatedCount == 0 {
		return 0, util.WrapError(fmt.Errorf("没有记录被更新"), "")
	}
	return int64(updatedCount), nil
}

// Delete 删除 admin 表的记录
func (s *StructAdmin) Delete(condition string) (int64, error) {
	deletedCount, _, err := GenericDelete(s.GetTableName(), condition, config.Config.MySQL.TablePrefix, "")
	if err != nil {
		return 0, util.WrapError(err, "删除记录失败:")
	}
	if deletedCount == 0 {
		return 0, util.WrapError(fmt.Errorf("没有记录被删除"), "")
	}
	return int64(deletedCount), nil
}

// 将查询结果映射到 StructAdmin 结构体
func (s *StructAdmin) mapResultToStructItem(result map[string]interface{}) (StructAdmin, error) {
	var item StructAdmin
	var ok bool

	if id, ok := result["id"].(int64); ok {
		item.ID = int(id)
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将id转换为int64：%v", result["id"]), "")
	}
	if item.Username, ok = result["username"].(string); !ok {
		return item, fmt.Errorf("类型断言失败: username")
	}
	if item.Password, ok = result["password"].(string); !ok {
		return item, fmt.Errorf("类型断言失败: password")
	}
	if item.Email, ok = result["email"].(string); !ok {
		return item, fmt.Errorf("类型断言失败: email")
	}
	if item.PhoneNumber, ok = result["phone_number"].(string); !ok {
		return item, fmt.Errorf("类型断言失败: phone_number")
	}
	if status, ok := result["status"].(int64); ok {
		item.Status = int(status)
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将status转换为int64：%v", result["status"]), "")
	}
	if item.CreateTime, ok = result["create_time"].(int64); !ok {
		return item, fmt.Errorf("类型断言失败: create_time")
	}
	if item.UpdateTime, ok = result["update_time"].(int64); !ok {
		return item, fmt.Errorf("类型断言失败: update_time")
	}
	if item.LastLoginTime, ok = result["last_login_time"].(int64); !ok {
		return item, fmt.Errorf("类型断言失败: last_login_time")
	}
	if item.Role, ok = result["role"].(string); !ok {
		return item, fmt.Errorf("类型断言失败: role")
	}
	if item.Salt, ok = result["salt"].(string); !ok {
		return item, fmt.Errorf("类型断言失败: salt")
	}
	if item.Avatar, ok = result["avatar"].(string); !ok {
		return item, fmt.Errorf("类型断言失败: avatar")
	}

	return item, nil
}
