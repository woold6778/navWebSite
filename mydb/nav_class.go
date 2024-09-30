package mydb

import (
	"fmt"
	"nav-web-site/config"
	"nav-web-site/util"
	"nav-web-site/util/log"
)

type StructNavClass struct {
	ID           int    `db:"id"`           //id
	Admin_id     int    `db:"admin_id"`     //管理员id
	Name         string `db:"name"`         //导航分类名称
	Parent_id    int    `db:"parent_id"`    //父级分类id
	Sort         int    `db:"sort"`         //排序
	Icon         string `db:"icon"`         //图标
	Description  string `db:"description"`  //描述
	Is_show      bool   `db:"is_show"`      //是否显示
	Is_recommend bool   `db:"is_recommend"` //是否推荐
	Is_hot       bool   `db:"is_hot"`       //是否热门
	Status       int    `db:"status"`       //状态:0=禁用,1=启用
	Create_time  int64  `db:"create_time"`  //创建时间
	Update_time  int64  `db:"update_time"`  //更新时间
}

// DefaultData 是一个构造函数，用于创建带有默认值的 StructNavClass 实例
func (s *StructNavClass) DefaultData() StructNavClass {
	return StructNavClass{
		Admin_id:     0,
		Name:         "",
		Parent_id:    0,
		Sort:         50,
		Icon:         "",
		Description:  "",
		Is_show:      true,  // 默认值
		Is_recommend: false, // 默认值
		Is_hot:       false, // 默认值
		Status:       1,     // 默认值
		Create_time:  util.GetTimestamp(10),
		Update_time:  util.GetTimestamp(10),
	}
}

// 获取表名（不含前后缀）
func (s *StructNavClass) GetTableName() string {
	return "nav_class"
}

// 获取插入数据时的必填字段
func (s *StructNavClass) GetRequiredFields() []string {
	return []string{
		"Name",
		"Parent_id",
		"Sort",
		"Is_show",
		"Status",
		"Create_time",
	}
}

// 插入数据时查重的字段
func (s StructNavClass) GetUniqueFields() []string {
	return []string{"Parent_id", "Name"}
}

// Find 方法根据条件查询单个 nav_class 记录
func (s *StructNavClass) Find(params QueryParams) (StructNavClass, error) {
	var class StructNavClass
	results, err := GenericSelect(Db, s.GetTableName(), params, config.Config.MySQL.TablePrefix, "")
	if err != nil {
		return class, util.WrapError(err, "Query failed(find):")
	}

	if len(results) > 0 {
		class, err = s.mapResultToStructItem(results[0])
		if err != nil {
			return class, util.WrapError(err, "将结果映射到StructNavClass时出错:")
		}
	} else {
		// 如果没有查询到数据，返回一个错误
		return class, util.WrapError(fmt.Errorf("EmptyData"), "")
	}
	return class, nil
}

// Select 方法查询 nav_class 表的数据
func (s *StructNavClass) Select(params QueryParams) ([]StructNavClass, int, error) {
	var list []StructNavClass
	results, err := GenericSelect(Db, s.GetTableName(), params, config.Config.MySQL.TablePrefix, "")
	if err != nil {
		return list, 400, util.WrapError(err, "Query failed(select):")
	}

	if len(results) > 0 {
		for _, result := range results {
			log.InfoLogger.Println("Processing result:", result)

			item, err := s.mapResultToStructItem(result)
			if err != nil {
				log.ErrorLogger.Println("将结果映射到StructNavClass时出错:", err)
				continue
			}

			list = append(list, item)
			log.InfoLogger.Println("Appended to list:", list[len(list)-1])
		}
	} else {
		// 如果没有查询到数据，返回一个错误
		return list, 200, util.WrapError(fmt.Errorf("EmptyData"), "")
	}
	return list, 0, nil
}

// Insert 方法插入新的 nav_class 记录
func (s *StructNavClass) Insert(datas []StructNavClass) (int, []int64, error) {
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

// Update 方法更新 nav_class 记录
func (s *StructNavClass) Update(datas []StructNavClass, condition string) (int, []int64, error) {
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

// Delete 方法删除 nav_class 记录
func (s *StructNavClass) Delete(condition string) (int, []int64, error) {
	count, ids, err := GenericDelete(
		s.GetTableName(),
		condition,
		config.Config.MySQL.TablePrefix,
		"",
	)
	if err != nil {
		return 0, ids, err
	}
	return count, ids, nil
}

// mapResultToStructItem 将查询结果映射到结构体
func (s *StructNavClass) mapResultToStructItem(result map[string]interface{}) (StructNavClass, error) {
	var item StructNavClass
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

	if item.Name, ok = result["name"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将name转换为string：%v", result["name"]), "")
	}

	if parentID, ok := result["parent_id"].(int64); ok {
		item.Parent_id = int(parentID)
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将parent_id转换为int64：%v", result["parent_id"]), "")
	}

	if sort, ok := result["sort"].(int64); ok {
		item.Sort = int(sort)
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将sort转换为int64：%v", result["sort"]), "")
	}

	if item.Icon, ok = result["icon"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将icon转换为string：%v", result["icon"]), "")
	}

	if item.Description, ok = result["description"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将description转换为string：%v", result["description"]), "")
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

	if isHot, ok := result["is_hot"].(int64); ok {
		item.Is_hot = isHot == 1
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将is_hot转换为int64：%v", result["is_hot"]), "")
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
