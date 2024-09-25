package mydb

import (
	"database/sql"
	"fmt"
	"nav-web-site/config"
	"nav-web-site/util"
	"reflect"
	"strings"
)

// 定义一个结构体来封装查询参数
type QueryParams struct {
	Condition string
	OrderBy   string
	Limit     int
	Page      int
	PageSize  int
}

// UpdateData 包含要更新的数据和对应的查询条件
type UpdateData[T any] struct {
	Data      T
	Condition string
}

// 定义一个接口，包含 GetUniqueFields 方法
type UniqueFieldGetter interface {
	GetUniqueFields() []string
}

// 通用查询函数
func GenericSelect(db *sql.DB, tableName string, params QueryParams, tablePrefix string, tableSuffix string) ([]map[string]interface{}, error) {
	// 设置默认值
	if tablePrefix == "" {
		tablePrefix = config.Config.MySQL.TablePrefix
	}
	if tableSuffix == "" {
		tableSuffix = ""
	}
	// 构建带前后缀的表名
	fullTableName := fmt.Sprintf("%s%s%s", tablePrefix, tableName, tableSuffix)
	// 构建 SQL 查询语句
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", fullTableName, params.Condition)

	// 如果有排序条件，添加 ORDER BY 子句
	if params.OrderBy != "" {
		query = fmt.Sprintf("%s ORDER BY %s", query, params.OrderBy)
	}

	// 如果 limit 有值且大于 0，添加 LIMIT 子句
	if params.Limit > 0 {
		// 如果 page 有值且大于 0，计算 OFFSET 并添加分页支持
		if params.Page > 0 {
			offset := (params.Page - 1) * params.PageSize
			query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, params.Limit, offset)
		} else {
			// 如果 page 没有值或等于 0，不进行分页，只添加 LIMIT 子句
			query = fmt.Sprintf("%s LIMIT %d", query, params.Limit)
		}
	}

	// 打印构建的查询语句
	fmt.Println("Constructed Query:", query)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//fmt.Printf("Query rows: %+v\n", rows)

	columns, _ := rows.Columns()
	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		result := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]

			// 检查 nil，即数据库中的 NULL 值
			if val == nil {
				result[col] = "" // 将 NULL 转换为空字符串
				continue
			}

			// 如果 val 是 []byte 类型，将其转换为字符串
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}

			result[col] = v
		}

		results = append(results, result)
	}

	fmt.Printf("Query results: %+v\n", results)
	/*
		if len(results) == 0 {
			return nil, errors.New("no results found")
		}
	*/

	return results, nil
}

// GenericInsert 通用批量数据插入操作
func GenericInsert[T any](tableName string, datas []T, requiredFields []string, tablePrefix string, tableSuffix string) (int, []int64, error) {
	if len(datas) == 0 {
		return 0, nil, util.WrapError(fmt.Errorf("数据数组为空"), "")
	}
	// 设置默认值
	if tablePrefix == "" {
		tablePrefix = config.Config.MySQL.TablePrefix
	}
	if tableSuffix == "" {
		tableSuffix = ""
	}
	// 构建带前后缀的表名
	fullTableName := fmt.Sprintf("%s%s%s", tablePrefix, tableName, tableSuffix)

	var insertedIDs []int64
	var nextID = int64(0)
	var err error
	// 将datas转换为[]interface{}
	interfaceDatas := make([]interface{}, len(datas))
	for i, v := range datas {
		//如果必须字段（在变量requiredFields里面）没有有效值就报错
		for _, field := range requiredFields {
			value := reflect.ValueOf(v).FieldByName(field)
			if !value.IsValid() || value.IsZero() {
				return 0, nil, util.WrapError(fmt.Errorf("必填字段 %s 缺失或为空值", field), "")
			}
		}

		// 检查是否已存在相同的记录
		uniqueFieldGetter, ok := interface{}(v).(UniqueFieldGetter)
		if !ok {
			return 0, nil, util.WrapError(fmt.Errorf("类型 %T 没有实现 UniqueFieldGetter 接口", v), "")
		}
		uniqueFields := uniqueFieldGetter.GetUniqueFields()
		exists, err := CheckExistingRecord(Db, v, uniqueFields, tableName, tablePrefix, tableSuffix)
		if err != nil {
			return 0, nil, util.WrapError(err, "检查记录是否存在时发生错误:")
		}
		if exists {
			continue // 如果记录已存在，跳过此条数据
		}
		// 检测T里面的ID的数据类型
		idField := reflect.ValueOf(&v).Elem().FieldByName("ID")
		if !idField.IsValid() {
			return 0, nil, util.WrapError(fmt.Errorf("类型 %T 没有 ID 字段", v), "")
		}
		switch idField.Kind() {
		case reflect.Int, reflect.Int64:
			//获取一个新的ID
			nextID, err = GetNextID(fullTableName)
			if err != nil {
				return 0, nil, util.WrapError(err, "获取下一个ID失败:")
			}
			// 使用反射设置 ID 字段(即v.ID)
			idField.SetInt(nextID)
		default:
			return 0, nil, util.WrapError(fmt.Errorf("ID 字段的类型 %s 不受支持", idField.Kind()), "")
		}
		interfaceDatas[i] = v
		insertedIDs = append(insertedIDs, nextID) //返回的ID列表
	}

	// 移除interfaceDatas中的nil值
	validInterfaceDatas := make([]interface{}, 0)
	for _, data := range interfaceDatas {
		if data != nil {
			validInterfaceDatas = append(validInterfaceDatas, data)
		}
	}
	interfaceDatas = validInterfaceDatas
	sql, valueArgs, err := GenerateInsertSQL(fullTableName, interfaceDatas)
	if err != nil {
		return 0, nil, util.WrapError(err, "生成SQL失败:")
	}

	util.InfoLogger.Println("生成的SQL语句:", sql)

	// 执行SQL语句
	result, err := Db.Exec(sql, valueArgs...)
	if err != nil {
		return 0, nil, util.WrapError(err, "执行SQL失败:")
	}

	// 获取插入的行数
	insertedCount, err := result.RowsAffected()
	if err != nil {
		return 0, nil, util.WrapError(err, "获取插入行数失败:")
	}

	// 如果没有插入任何行，返回错误
	if insertedCount == 0 {
		return 0, nil, util.WrapError(fmt.Errorf("没有插入任何数据"), "")
	}

	return int(insertedCount), insertedIDs, nil
}

// GenericUpdate 批量更新数据的通用函数
func GenericUpdate[T any](tableName string, datas []T, condition string, tablePrefix string, tableSuffix string) (int, []int64, error) {
	// 设置默认值
	if tablePrefix == "" {
		tablePrefix = config.Config.MySQL.TablePrefix
	}
	if tableSuffix == "" {
		tableSuffix = ""
	}
	// 构建带前后缀的表名
	fullTableName := fmt.Sprintf("%s%s%s", tablePrefix, tableName, tableSuffix)

	// 将数据转换为 []interface{}
	interfaceDatas := interfaceSlice(datas)

	// 获取字段名
	firstItem := interfaceDatas[0]
	val := reflect.ValueOf(firstItem)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typeOfS := val.Type()

	var columns []string
	for i := 0; i < val.NumField(); i++ {
		field := typeOfS.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			columns = append(columns, dbTag)
		} else {
			columns = append(columns, typeOfS.Field(i).Name)
		}
	}

	// 生成更新SQL
	var setClauses []string
	var valueArgs []interface{}
	for _, data := range interfaceDatas {
		val := reflect.ValueOf(data)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		var setClause []string
		for i := 0; i < val.NumField(); i++ {
			setClause = append(setClause, fmt.Sprintf("%s = ?", columns[i]))
			valueArgs = append(valueArgs, val.Field(i).Interface())
		}
		setClauses = append(setClauses, strings.Join(setClause, ", "))
	}

	updateSQL := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		fullTableName,
		strings.Join(setClauses, ", "),
		condition)

	// 打印生成的SQL语句
	util.InfoLogger.Println("生成的Update SQL语句:", updateSQL)

	// 执行SQL语句
	result, err := Db.Exec(updateSQL, valueArgs...)
	if err != nil {
		return 0, nil, util.WrapError(err, "执行SQL失败:"+updateSQL)
	}

	// 获取受影响的行数
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, nil, util.WrapError(err, "获取受影响的行数失败:")
	}

	// 获取更新的记录ID列表
	var updatedIDs []int64
	for _, data := range interfaceDatas {
		val := reflect.ValueOf(data)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		idField := val.FieldByName("ID")
		if idField.IsValid() && (idField.Kind() == reflect.Int64 || idField.Kind() == reflect.Int) {
			updatedIDs = append(updatedIDs, idField.Int())
		}
	}

	return int(affectedRows), updatedIDs, nil
}

// GenericDelete 通用批量数据删除操作
func GenericDelete(tableName string, condition string, tablePrefix string, tableSuffix string) (int, []int64, error) {
	// 设置默认值
	if tablePrefix == "" {
		tablePrefix = config.Config.MySQL.TablePrefix
	}
	if tableSuffix == "" {
		tableSuffix = ""
	}
	// 构建带前后缀的表名
	fullTableName := fmt.Sprintf("%s%s%s", tablePrefix, tableName, tableSuffix)

	// 构建删除SQL语句
	deleteSQL := fmt.Sprintf("DELETE FROM %s WHERE %s", fullTableName, condition)

	// 打印生成的SQL语句
	util.InfoLogger.Println("生成的Delete SQL语句:", deleteSQL)

	// 执行SQL语句
	result, err := Db.Exec(deleteSQL)
	if err != nil {
		return 0, nil, util.WrapError(err, "执行SQL失败:"+deleteSQL)
	}

	// 获取受影响的行数
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, nil, util.WrapError(err, "获取受影响的行数失败:")
	}

	// 获取删除的记录ID列表
	var deletedIDs []int64
	rows, err := Db.Query(fmt.Sprintf("SELECT ID FROM %s WHERE %s", fullTableName, condition))
	if err != nil {
		return 0, nil, util.WrapError(err, "查询删除的记录ID失败:")
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return 0, nil, util.WrapError(err, "扫描记录ID失败:")
		}
		deletedIDs = append(deletedIDs, id)
	}

	return int(affectedRows), deletedIDs, nil
}

// CheckExistingRecord 通用检查数据是否存在的函数
func CheckExistingRecord[T any](db *sql.DB, data T, uniqueFields []string, tableName string, tablePrefix string, tableSuffix string) (bool, error) {
	condition := ""
	for i, field := range uniqueFields {
		if i > 0 {
			condition += " AND "
		}
		value := reflect.ValueOf(data).FieldByName(field).Interface()
		if str, ok := value.(string); ok {
			condition += fmt.Sprintf("%s = '%s'", field, str)
		} else {
			condition += fmt.Sprintf("%s = %v", field, value)
		}
	}
	params := QueryParams{
		Condition: condition,
		Limit:     1,
	}

	results, err := GenericSelect(db, tableName, params, tablePrefix, tableSuffix)
	if err != nil {
		return false, util.WrapError(err, "查询记录是否存在时发生错误:")
	}

	return len(results) > 0, nil
}

// interfaceSlice 将任意类型的切片转换为 []interface{}
func interfaceSlice[T any](slice []T) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

// GenerateInsertSQL 生成插入多条数据的SQL语句
func GenerateInsertSQL(tableName string, datas []interface{}) (string, []interface{}, error) {
	if len(datas) == 0 {
		return "", nil, util.WrapError(fmt.Errorf("数据数组为空"), "")
	}

	// 获取字段名
	firstItem := datas[0]
	val := reflect.ValueOf(firstItem)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typeOfS := val.Type()

	var columns []string
	for i := 0; i < val.NumField(); i++ {
		field := typeOfS.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			columns = append(columns, dbTag)
		} else {
			columns = append(columns, typeOfS.Field(i).Name)
		}

	}

	// 生成插入SQL
	var valueStrings []string
	var valueArgs []interface{}
	for _, data := range datas {
		val := reflect.ValueOf(data)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		var values []string
		for i := 0; i < val.NumField(); i++ {
			values = append(values, "?")
			valueArgs = append(valueArgs, val.Field(i).Interface())
		}
		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(values, ", ")))
	}

	insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(valueStrings, ", "))

	return insertSQL, valueArgs, nil
}
