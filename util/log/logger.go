package log

import (
	"log"
	"os"
	"time"
)

// 全局日志变量
var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

// 日志初始化
func InitLoggers() {
	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	// 指定日志文件夹路径
	//logDir := cwd + "/log"
	currentMonth := time.Now().Format("01")
	currentDay := time.Now().Format("02")
	logDir := cwd + "/log/" + currentMonth

	// 检查目录是否存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// 创建目录
		if err := os.Mkdir(logDir, 0755); err != nil {
			log.Fatalf("Failed to create log directory: %v", err)
		}
	}
	infoLogPath := logDir + "/" + currentDay + ".log"
	errorLogPath := logDir + "/" + currentDay + "_error.log"

	// 打开日志文件
	infoLogFile, err := os.OpenFile(infoLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	errorLogFile, err := os.OpenFile(errorLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	// 创建 log.Logger 实例
	InfoLogger = log.New(infoLogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorLogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
