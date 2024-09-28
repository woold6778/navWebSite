package mydb

import (
	"fmt"
	"nav-web-site/config"
	"nav-web-site/util"
)

// StructNews 定义 news 结构体
type StructNews struct {
	ID             int    `db:"id"`
	Admin_id       int    `db:"admin_id"`
	Class_id       int    `db:"class_id"`
	Title          string `db:"title"`
	Subtitle       string `db:"subtitle"`
	Url            string `db:"url"`
	Description    string `db:"description"`
	Icon           string `db:"icon"`
	Keywords       string `db:"keywords"`
	Sort           int    `db:"sort"`
	Is_show        bool   `db:"is_show"`
	Status         int    `db:"status"`
	Create_time    int64  `db:"create_time"`
	Author         string `db:"author"`
	Source         string `db:"source"`
	View_count     int    `db:"view_count"`
	Comment_count  int    `db:"comment_count"`
	Language       string `db:"language" default:"cn"` //文章使用的语言
	Is_hot         bool   `db:"is_hot"`                //是否热门
	Is_headline    bool   `db:"is_headline"`           //是否头条
	Is_recommended bool   `db:"is_recommended"`        //是否推荐
	Content        string `db:"content"`               // 文章内容
}

// StructNewsContent 定义 news 内容结构体
type StructNewsContent struct {
	ID      int    `db:"id"`
	NewsID  int    `db:"news_id"` // 关联的新闻ID
	Content string `db:"content"` // 文章内容
}

// GetTableName 获取表名
func (s *StructNews) GetTableName() string {
	return "news"
}

// GetRequiredFields 获取必填字段
func (s *StructNews) GetRequiredFields() []string {
	return []string{"Title", "Description"}
}

// 插入数据时查重的字段
func (s StructNews) GetUniqueFields() []string {
	return []string{"Title", "Source"}
}
func (s StructNewsContent) GetUniqueFields() []string {
	return []string{}
}

// find 方法查询 news 表的单条数据，返回的数据要包括news_content表里面的content字段
func (s *StructNews) Find(condition string) (StructNews, error) {
	// 构建带前后缀的表名
	//fullTableName := fmt.Sprintf("%s%s%s", config.Config.MySQL.TablePrefix, s.GetTableName(), "")
	var item StructNews

	// 使用StructNews.Select查询news表的数据
	params := QueryParams{
		Condition: condition,
		Page:      1,
		PageSize:  1,
	}
	newsList, _, err := s.Select(params)
	if err != nil {
		return item, util.WrapError(err, "查询失败:")
	}

	if len(newsList) > 0 {
		item = newsList[0]

		// 查询 news_content 表中的内容
		fullTableName_content := fmt.Sprintf("%s%s%s", config.Config.MySQL.TablePrefix, "news_content", "")
		contentQuery := fmt.Sprintf("SELECT content FROM %s WHERE news_id = %d", fullTableName_content, item.ID)
		contentResult, err := Db.Query(contentQuery)
		if err != nil {
			return item, util.WrapError(err, "查询内容失败:")
		}
		defer contentResult.Close()

		if contentResult.Next() {
			var content string
			err := contentResult.Scan(&content)
			if err != nil {
				return item, util.WrapError(err, "扫描内容失败:")
			}
			item.Content = content
		}
	} else {
		return item, util.WrapError(fmt.Errorf("未找到记录"), "")
	}

	return item, nil
}

// Select 方法查询 news 表的数据
func (s *StructNews) Select(params QueryParams) ([]StructNews, int, error) {
	var list []StructNews
	results, err := GenericSelect(Db, s.GetTableName(), params, config.Config.MySQL.TablePrefix, "")
	if err != nil {
		return list, 400, util.WrapError(err, "Query failed(select):")
	}

	if len(results) > 0 {
		for _, result := range results {
			util.InfoLogger.Println("Processing result:", result)

			item, err := s.mapResultToStructItem(result)
			if err != nil {
				util.ErrorLogger.Println("将结果映射到StructNews时出错:", err)
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

// Insert 方法插入新的 news 记录,如果有content就能插入news_content表，否在就把content拿掉
func (s *StructNews) Insert(datas []StructNews) (int, []int64, error) {
	var ids []int64
	for i, data := range datas {
		// 插入 news 表
		_, tempIds, err := GenericInsert(
			s.GetTableName(),
			[]StructNews{data},
			s.GetRequiredFields(),
			config.Config.MySQL.TablePrefix,
			"",
		)
		if err != nil {
			return 0, ids, err
		}
		ids = append(ids, tempIds...)

		// 如果有 content 字段，插入 news_content 表
		if data.Content != "" {
			contentData := StructNewsContent{
				NewsID:  int(tempIds[0]),
				Content: data.Content,
			}
			_, _, err := GenericInsert(
				"news_content",
				[]StructNewsContent{contentData},
				[]string{"NewsID", "Content"},
				config.Config.MySQL.TablePrefix,
				"",
			)
			if err != nil {
				return 0, ids, err
			}
		} else {
			// 如果没有 content 字段，移除 content
			datas[i].Content = ""
		}
	}
	return len(ids), ids, nil
}

// Update 方法更新 news 记录
func (s *StructNews) Update(datas []StructNews, condition string) (int, []int64, error) {
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

	for _, data := range datas {
		if data.Content != "" {
			contentData := StructNewsContent{
				NewsID:  data.ID,
				Content: data.Content,
			}
			// 先判断news_id对应的content是否存在
			existingContent, err := GenericSelect(
				Db,
				"news_content",
				QueryParams{
					Condition: fmt.Sprintf("news_id = %d", data.ID),
				},
				config.Config.MySQL.TablePrefix,
				"",
			)
			if err != nil {
				return 0, ids, err
			}

			if len(existingContent) > 0 {
				// 如果存在就更新
				_, _, err := GenericUpdate(
					"news_content",
					[]StructNewsContent{contentData},
					fmt.Sprintf("news_id = %d", data.ID),
					config.Config.MySQL.TablePrefix,
					"",
				)
				if err != nil {
					return 0, ids, err
				}
			} else {
				// 否则就用GenericInsert新增
				_, _, err := GenericInsert(
					"news_content",
					[]StructNewsContent{contentData},
					[]string{"NewsID", "Content"},
					config.Config.MySQL.TablePrefix,
					"",
				)
				if err != nil {
					return 0, ids, err
				}
			}
		}
	}
	return count, ids, nil
}

// Delete 方法删除 news 记录
func (s *StructNews) Delete(condition string) (int, []int64, error) {
	count, ids, err := GenericDelete(
		s.GetTableName(),
		condition,
		config.Config.MySQL.TablePrefix,
		"",
	)
	if err != nil {
		return 0, ids, err
	}

	// 同步删除 news_content 表中的对应数据
	for _, id := range ids {
		_, _, err := GenericDelete(
			"news_content",
			fmt.Sprintf("news_id = %d", id),
			config.Config.MySQL.TablePrefix,
			"",
		)
		if err != nil {
			return 0, ids, err
		}
	}

	return count, ids, nil
}

// mapResultToStructItem 将查询结果映射到结构体
func (s *StructNews) mapResultToStructItem(result map[string]interface{}) (StructNews, error) {
	var item StructNews
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

	if isShow, ok := result["is_show"].(int64); ok {
		item.Is_show = isShow == 1
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将is_show转换为int64：%v", result["is_show"]), "")
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

	if item.Author, ok = result["author"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将author转换为string：%v", result["author"]), "")
	}

	if item.Source, ok = result["source"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将source转换为string：%v", result["source"]), "")
	}

	if viewCount, ok := result["view_count"].(int64); ok {
		item.View_count = int(viewCount)
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将view_count转换为int64：%v", result["view_count"]), "")
	}

	if commentCount, ok := result["comment_count"].(int64); ok {
		item.Comment_count = int(commentCount)
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将comment_count转换为int64：%v", result["comment_count"]), "")
	}

	if item.Language, ok = result["language"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将language转换为string：%v", result["language"]), "")
	}

	if isHot, ok := result["is_hot"].(int64); ok {
		item.Is_hot = isHot == 1
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将is_hot转换为int64：%v", result["is_hot"]), "")
	}

	if isHeadline, ok := result["is_headline"].(int64); ok {
		item.Is_headline = isHeadline == 1
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将is_headline转换为int64：%v", result["is_headline"]), "")
	}

	if isRecommended, ok := result["is_recommended"].(int64); ok {
		item.Is_recommended = isRecommended == 1
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将is_recommended转换为int64：%v", result["is_recommended"]), "")
	}

	return item, nil
}
