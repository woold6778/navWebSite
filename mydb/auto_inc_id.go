package mydb

import (
	"encoding/json"
	"fmt"
	"nav-web-site/util"
	"os"
	"path/filepath"
	"sync"
)

var (
	idMap   = make(map[string]int64)
	idMutex sync.Mutex
	idFile  = "id_store.json"
)

// 初始化时从文件加载ID
func init() {
	loadIDs()
}

// 获取下一个ID
func GetNextID(tableName string) (int64, error) {
	idMutex.Lock()
	defer idMutex.Unlock()

	// 获取当前ID
	currentID, exists := idMap[tableName]
	if !exists {
		currentID = 0
	}

	// 生成下一个ID
	nextID := currentID + 1
	idMap[tableName] = nextID

	// 将ID持久化到文件
	err := saveIDs()
	if err != nil {
		return 0, util.WrapError(err, "ID持久化失败")
	}

	return nextID, nil
}

// 将ID持久化到文件
func saveIDs() error {
	file, err := os.Create(idFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(idMap)
}

// 从文件加载ID
func loadIDs() {
	absPath, err := filepath.Abs(idFile)
	if err != nil {
		fmt.Println("获取文件绝对路径失败:", err)
		return
	}
	fmt.Println("ID文件的绝对路径:", absPath)

	file, err := os.Open(idFile)
	if err != nil {
		if os.IsNotExist(err) {
			return // 文件不存在，跳过加载
		}
		fmt.Println("加载ID失败:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&idMap)
	if err != nil {
		fmt.Println("解析ID文件失败:", err)
	}
}
