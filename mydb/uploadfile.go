package mydb

import (
	"fmt"
	"nav-web-site/config"
	"nav-web-site/util"
)

type StructUploadFile struct {
	ID         int    `db:"id"`          // 文件ID
	FileName   string `db:"file_name"`   // 文件名
	FilePath   string `db:"file_path"`   // 文件路径
	FileSize   int64  `db:"file_size"`   // 文件大小
	Hash       string `db:"hash"`        // 哈希
	FileType   string `db:"file_type"`   // 文件类型
	Extension  string `db:"extension"`   // 扩展名
	UploadTime int64  `db:"upload_time"` // 上传时间
}

// DefaultData 是一个构造函数，用于创建带有默认值的 StructUploadFile 实例
func (s *StructUploadFile) DefaultData() StructUploadFile {
	return StructUploadFile{
		FileName:   "",
		FilePath:   "",
		FileSize:   0,
		FileType:   "",
		UploadTime: util.GetTimestamp(10),
	}
}

// 获取表名（不含前后缀）
func (s *StructUploadFile) GetTableName() string {
	return "upload_file"
}

// 获取插入数据时的必填字段
func (s *StructUploadFile) GetRequiredFields() []string {
	return []string{
		"FileName",
		"FilePath",
		"FileSize",
		"FileType",
		"Hash",
	}
}

// 获取插入数据时查重的字段
func (s StructUploadFile) GetUniqueFields() []string {
	return []string{"Hash"}
}

// Insert 方法插入新的 upload_file 记录
func (s *StructUploadFile) Insert(datas []StructUploadFile) (int, []int64, error) {
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

// Select 方法查询 upload_file 表的数据
func (s *StructUploadFile) Select(params QueryParams) ([]StructUploadFile, int, error) {
	var list []StructUploadFile
	results, err := GenericSelect(Db, s.GetTableName(), params, config.Config.MySQL.TablePrefix, "")
	if err != nil {
		return list, 400, util.WrapError(err, "Query failed(select):")
	}

	if len(results) > 0 {
		for _, result := range results {
			util.InfoLogger.Println("Processing result:", result)

			item, err := s.mapResultToStructItem(result)
			if err != nil {
				util.ErrorLogger.Println("将结果映射到StructUploadFile时出错:", err)
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

// Find 方法根据条件查询单个 upload_file 记录
func (s *StructUploadFile) Find(params QueryParams) (StructUploadFile, error) {
	var item StructUploadFile
	results, _, err := s.Select(params)
	if err != nil {
		return item, util.WrapError(err, "Query failed(select):")
	}

	if len(results) > 0 {
		item = results[0]
	} else {
		return item, util.WrapError(fmt.Errorf("未找到记录"), "")
	}

	return item, nil
}

// mapResultToStructItem 将查询结果映射到结构体
func (s *StructUploadFile) mapResultToStructItem(result map[string]interface{}) (StructUploadFile, error) {
	var item StructUploadFile
	var ok bool

	if id, ok := result["id"].(int64); ok {
		item.ID = int(id)
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将id转换为int64：%v", result["id"]), "")
	}

	if item.FileName, ok = result["file_name"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将file_name转换为string：%v", result["file_name"]), "")
	}

	if item.FilePath, ok = result["file_path"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将file_path转换为string：%v", result["file_path"]), "")
	}

	if fileSize, ok := result["file_size"].(int64); ok {
		item.FileSize = fileSize
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将file_size转换为int64：%v", result["file_size"]), "")
	}

	if item.Hash, ok = result["hash"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将hash转换为string：%v", result["hash"]), "")
	}

	if item.FileType, ok = result["file_type"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将file_type转换为string：%v", result["file_type"]), "")
	}

	if item.Extension, ok = result["extension"].(string); !ok {
		return item, util.WrapError(fmt.Errorf("错误：无法将extension转换为string：%v", result["extension"]), "")
	}

	if uploadTime, ok := result["upload_time"].(int64); ok {
		item.UploadTime = uploadTime
	} else {
		return item, util.WrapError(fmt.Errorf("错误：无法将upload_time转换为int64：%v", result["upload_time"]), "")
	}

	return item, nil
}
