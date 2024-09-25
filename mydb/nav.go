package mydb

import (
	"fmt"
	"nav-web-site/config"
	"nav-web-site/util"
)

type StructNav struct {
	ID           int    `db:"id"`           // id
	Admin_id     int    `db:"admin_id"`     // 管理员id
	Class_id     int    `db:"class_id"`     // 分类id
	Title        string `db:"title"`        // 标题
	Subtitle     string `db:"subtitle"`     // 副标题
	Url          string `db:"url"`          // 链接地址
	Description  string `db:"description"`  // 描述
	Icon         string `db:"icon"`         // 图标
	Keywords     string `db:"keywords"`     // 关键词
	Sort         int    `db:"sort"`         // 排序
	Views        int    `db:"views"`        // 浏览量
	Is_show      bool   `db:"is_show"`      // 是否显示
	Is_recommend bool   `db:"is_recommend"` // 是否推荐
	Status       int    `db:"status"`       // 状态:0=禁用,1=启用
	Create_time  int64  `db:"create_time"`  // 创建时间
	Update_time  int64  `db:"update_time"`  // 更新时间
}

// DefaultData 是一个构造函数，用于创建带有默认值的 StructNav 实例
func (s *StructNav) DefaultData() StructNav {
	return StructNav{
		Admin_id:     0,
		Class_id:     0,
		Title:        "",
		Subtitle:     "",
		Url:          "",
		Description:  "",
		Icon:         "",
		Keywords:     "",
		Sort:         50,
		Views:        0,
		Is_show:      true,  // 默认值
		Is_recommend: false, // 默认值
		Status:       1,     // 默认值
		Create_time:  util.GetTimestamp(10),
		Update_time:  util.GetTimestamp(10),
	}
}

// 获取表名（不含前后缀）
func (s *StructNav) GetTableName() string {
	return "nav"
}

// 获取插入数据时的必填字段
func (s *StructNav) GetRequiredFields() []string {
	return []string{
		"Class_id",
		"Title",
		"Url",
		"Sort",
		"Is_show",
		"Status",
		"Create_time",
	}
}

// 获取插入数据时查重的字段
func (s StructNav) GetUniqueFields() []string {
	return []string{"Class_id", "Title", "Url"}
}

// Find 方法根据条件查询单个 nav 记录
func (s *StructNav) Find(params QueryParams) (StructNav, error) {
	var nav StructNav
	results, err := GenericSelect(Db, s.GetTableName(), params, config.Config.MySQL.TablePrefix, "")
	if err != nil {
		return nav, util.WrapError(err, "Query failed(find):")
	}

	if len(results) > 0 {
		nav, err = s.mapResultToStructItem(results[0])
		if err != nil {
			return nav, util.WrapError(err, "将结果映射到StructNav时出错:")
		}
	} else {
		// 如果没有查询到数据，返回一个错误
		return nav, util.WrapError(fmt.Errorf("EmptyData"), "")
	}
	return nav, nil
}

// Select 方法查询 nav 表的数据
func (s *StructNav) Select(params QueryParams) ([]StructNav, int, error) {
	var list []StructNav
	results, err := GenericSelect(Db, s.GetTableName(), params, config.Config.MySQL.TablePrefix, "")
	if err != nil {
		return list, 400, util.WrapError(err, "Query failed(select):")
	}

	if len(results) > 0 {
		for _, result := range results {
			util.InfoLogger.Println("Processing result:", result)

			item, err := s.mapResultToStructItem(result)
			if err != nil {
				util.ErrorLogger.Println("将结果映射到StructNav时出错:", err)
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

// Insert 插入数据
func (s *StructNav) Insert(datas []StructNav) (int, []int64, error) {
	count, ids, err := GenericInsert(
		s.GetTableName(),
		datas,
		s.GetRequiredFields(),
		config.Config.MySQL.TablePrefix,
		"",
	)
	if err != nil {
		return 0, ids, err
	}
	return count, ids, nil
}

func (s *StructNav) Update(datas []StructNav, condition string) (int, []int64, error) {
	count, ids, err := GenericUpdate(
		s.GetTableName(),
		datas,
		condition,
		config.Config.MySQL.TablePrefix,
		"",
	)
	if err != nil {
		return 0, ids, err
	}
	return count, ids, nil
}

// Del 删除数据
func (s *StructNav) Delete(condition string) (int, []int, error) {
	params := QueryParams{
		Condition: condition,
	}
	recordsToDelete, _, err := s.Select(params)
	if err != nil {
		return 0, nil, util.WrapError(err, "查询要删除的记录失败:")
	}

	var idsToDelete []int
	for _, record := range recordsToDelete {
		idsToDelete = append(idsToDelete, record.ID)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", s.GetTableName(), condition)

	result, err := Db.Exec(query)
	if err != nil {
		return 0, nil, util.WrapError(err, "执行删除操作失败:")
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, nil, util.WrapError(err, "获取受影响的行数失败:")
	}

	return int(affectedRows), idsToDelete, nil
}

// mapResultToStructItem 将查询结果映射到结构体
func (s *StructNav) mapResultToStructItem(result map[string]interface{}) (StructNav, error) {
	var item StructNav
	var ok bool

	if id, ok := result["id"].(int64); ok {
		item.ID = int(id)
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将id转换为int64：%v", result["id"]), "")
	}

	if adminID, ok := result["admin_id"].(int64); ok {
		item.Admin_id = int(adminID)
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将admin_id转换为int64：%v", result["admin_id"]), "")
	}

	if classID, ok := result["class_id"].(int64); ok {
		item.Class_id = int(classID)
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将class_id转换为int64：%v", result["class_id"]), "")
	}

	if item.Title, ok = result["title"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将title转换为string：%v", result["title"]), "")
	}

	if item.Subtitle, ok = result["subtitle"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将subtitle转换为string：%v", result["subtitle"]), "")
	}

	if item.Url, ok = result["url"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将url转换为string：%v", result["url"]), "")
	}

	if item.Description, ok = result["description"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将description转换为string：%v", result["description"]), "")
	}

	if item.Icon, ok = result["icon"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将icon转换为string：%v", result["icon"]), "")
	}

	if item.Keywords, ok = result["keywords"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将keywords转换为string：%v", result["keywords"]), "")
	}

	if sort, ok := result["sort"].(int64); ok {
		item.Sort = int(sort)
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将sort转换为int64：%v", result["sort"]), "")
	}

	if views, ok := result["views"].(int64); ok {
		item.Views = int(views)
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将views转换为int64：%v", result["views"]), "")
	}

	if isShow, ok := result["is_show"].(int64); ok {
		item.Is_show = isShow == 1
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将is_show转换为int64：%v", result["is_show"]), "")
	}

	if isRecommend, ok := result["is_recommend"].(int64); ok {
		item.Is_recommend = isRecommend == 1
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将is_recommend转换为int64：%v", result["is_recommend"]), "")
	}

	if status, ok := result["status"].(int64); ok {
		item.Status = int(status)
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将status转换为int64：%v", result["status"]), "")
	}

	if createTime, ok := result["create_time"].(int64); ok {
		item.Create_time = createTime
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将create_time转换为int64：%v", result["create_time"]), "")
	}

	if updateTime, ok := result["update_time"].(int64); ok {
		item.Update_time = updateTime
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将update_time转换为int64：%v", result["update_time"]), "")
	}

	return item, nil
}
