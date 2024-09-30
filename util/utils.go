// util/utils.go
package util

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"math/rand"
	"nav-web-site/config"
	"nav-web-site/util/log"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// api接口返回输出的结构体
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 封装的错误处理函数
func WrapError(err error, msg string) error {
	if err == nil {
		return nil
	}

	if config.Config.Base.Debug {
		// 获取调用者的文件名和行号
		_, file, line, ok := runtime.Caller(1)
		if ok {
			return fmt.Errorf("%s: %v (at %s:%d)", msg, err, file, line)
		}
	}

	// 如果配置文件里debur=false或者获取文件名和行号失败，返回普通错误信息
	return fmt.Errorf("%s: %v", msg, err)
}

// 获取当前时间的10位时间戳
// GetTimestamp 根据参数返回不同位数的时间戳
func GetTimestamp(digits int) int64 {
	now := time.Now()
	switch digits {
	case 10:
		return now.Unix()
	case 13:
		return now.UnixNano() / 1e6
	case 16:
		return now.UnixNano() / 1e3
	case 19:
		return now.UnixNano()
	default:
		return now.Unix() // 默认返回10位时间戳
	}
}

// MD5Hash 对字符串进行MD5加密
func MD5Hash(text string, salt string) string {
	hash := md5.New()
	io.WriteString(hash, text+salt)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Hash256 对字符串进行SHA-256哈希
func Hash256(text string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(text + salt))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// GenerateToken 生成 token	，用于验证请求的合法性	，把域名+ip+时间戳+随机数 进行 md5 加密
func GenerateToken(token_str string) string {
	hash := md5.New()                       // 创建一个新的 md5 哈希对象
	io.WriteString(hash, token_str)         // 将输入的字符串 token_str 写入到哈希对象中
	return fmt.Sprintf("%x", hash.Sum(nil)) // 计算哈希值并格式化为十六进制字符串返回
}

// IsNumeric 检查字符串是否可以转换为数字
func IsNumeric(s string) bool {
	_, err := strconv.Atoi(s) // 尝试转换为整数
	if err != nil {
		_, err = strconv.ParseFloat(s, 64) // 尝试转换为浮点数
		return err == nil
	}
	return true
}

// 获取客户端IP
func GetClientIP(r *http.Request) string {
	// 获取 X-Real-IP 或 X-Forwarded-For 的第一个 IP
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}
	// 如果没有 X-Real-IP 和 X-Forwarded-For，使用 RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "127.0.0.1"
	}
	return ip
}

// ConvertSliceToMap 将结构体切片转换为以指定字段为键的映射
func ConvertSliceToMap[T any](slice []T, keyField string) (map[string]T, error) {
	result := make(map[string]T)
	for _, item := range slice {
		value := reflect.ValueOf(item)
		if value.Kind() == reflect.Struct {
			field := value.FieldByName(keyField)
			if field.IsValid() {
				key := fmt.Sprintf("%v", field.Interface())
				result[key] = item
			} else {
				return nil, fmt.Errorf("字段 %s 在结构体中不存在", keyField)
			}
		} else {
			return nil, fmt.Errorf("输入切片的元素不是结构体")
		}
	}
	return result, nil
}

// port是字符串类型，结构有可能是如（80,443）或者（80,433,8000-8010）或者（80）等，需要转换为[]string
func ParsePortList(portStr string) []string {
	log.InfoLogger.Println("Parsing port list from string:", portStr)
	var portList []string
	portRanges := strings.Split(portStr, ",")
	for _, portRange := range portRanges {
		log.InfoLogger.Println("Processing port range:", portRange)
		if strings.Contains(portRange, "-") {
			ports := strings.Split(portRange, "-")
			start, err1 := strconv.Atoi(ports[0])
			end, err2 := strconv.Atoi(ports[1])
			if err1 != nil || err2 != nil {
				log.ErrorLogger.Println("Error parsing port range:", portRange, "Error:", err1, err2)
				continue
			}
			for port := start; port <= end; port++ {
				portList = append(portList, strconv.Itoa(port))
			}
		} else {
			port, err := strconv.Atoi(portRange)
			if err != nil {
				log.ErrorLogger.Println("Error parsing port:", portRange, "Error:", err)
				continue
			}
			portList = append(portList, strconv.Itoa(port))
		}
	}
	log.InfoLogger.Println("Parsed port list:", portList)
	return portList
}

// 生成随机字符串
func GenerateRandomString(length int, charType int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // 使用新的随机生成器
	var chars string
	switch charType {
	case 1:
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	case 2:
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?"
	case 3:
		chars = "0123456789"
	default:
		panic("无效的字符类型")
	}

	result := make([]byte, length)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}

	if charType == 3 && result[0] == '0' {
		result[0] = chars[r.Intn(len(chars)-1)+1] // 确保第一个字符不是0
	}

	return string(result)
}

/**
 * 字符串加密    主要用于url加密
 * @param string $str 要加密的字符串
 * @param int $encryption_number    加密用的加密码  1000000-9999999之间的整数
 */
func MyEncryption(str string, encryptionNumber int) string {
	if encryptionNumber < 1000000 || encryptionNumber > 9999999 {
		panic("加密码必须在1000000-9999999之间")
	}

	var result strings.Builder
	for _, char := range str {
		// 使用 Unicode 编码
		ascii := int(char)
		result.WriteString(fmt.Sprintf("%07d", ascii+encryptionNumber))
	}
	return result.String()
}

/**
 * 字符串解密    主要用于url解密
 * @param string $str 要解密的字符串
 * @param int $encryption_number    加密用的加密码  1000000-9999999之间的整数
 */
func MyDecryption(encryptinoStr string, encryptionNumber int) string {
	if encryptionNumber < 1000000 || encryptionNumber > 9999999 {
		panic("解密码必须在1000000-9999999之间")
	}
	var result strings.Builder
	for i := 0; i < len(encryptinoStr); i += 7 {
		// 使用 Unicode 编码
		ascii, _ := strconv.Atoi(encryptinoStr[i : i+7])
		char := string(rune(ascii - encryptionNumber))
		result.WriteString(char)
	}
	return result.String()
}

/**
 * 获取1000000-9999999之间的随机整数
 * @return int 返回1000000-9999999之间的随机整数
 */
func GetRandomEncryptionNumber() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(9000000) + 1000000
}

/**
 * 验证邮箱地址的合法性
 * @param string email 要验证的邮箱地址
 * @return bool 返回邮箱地址是否合法
 */
func IsValidEmail(email string) bool {
	// 使用正则表达式验证邮箱地址
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emailRegex := regexp.MustCompile(emailRegexPattern)
	return emailRegex.MatchString(email)
}

/**
 * 验证各国手机号的合法性
 * @param string phoneNumber 要验证的手机号
 * @param string countryCode 国家代码 (例如: "CN" 表示中国, "US" 表示美国)
 * @return bool 返回手机号是否合法
 */
func IsValidPhoneNumber(phoneNumber string, countryCode string) bool {
	var phoneRegexPattern string

	switch countryCode {
	case "CN": // 中国
		phoneRegexPattern = `^1[3-9]\d{9}$`
	case "US": // 美国
		phoneRegexPattern = `^\+1\d{10}$`
	case "IN": // 印度
		phoneRegexPattern = `^\+91\d{10}$`
	case "UK": // 英国
		phoneRegexPattern = `^\+44\d{10}$`
	case "JP": // 日本
		phoneRegexPattern = `^\+81\d{10}$`
	default:
		return false
	}

	phoneRegex := regexp.MustCompile(phoneRegexPattern)
	return phoneRegex.MatchString(phoneNumber)
}

/**
 * 生成设备指纹
 * @param *http.Request request 请求对象
 * @return string 返回生成的设备指纹
 */
func GenerateDeviceFingerprint(request *http.Request) string {
	// 获取客户端IP地址
	clientIP := request.RemoteAddr

	// 获取User-Agent
	userAgent := request.UserAgent()

	// 获取Accept-Language
	acceptLanguage := request.Header.Get("Accept-Language")

	// 获取请求头中的其他信息
	acceptEncoding := request.Header.Get("Accept-Encoding")
	connection := request.Header.Get("Connection")

	// 将以上信息拼接成一个字符串
	fingerprintSource := clientIP + userAgent + acceptLanguage + acceptEncoding + connection

	// 对拼接后的字符串进行MD5哈希
	fingerprint := MD5Hash(fingerprintSource, "")

	return fingerprint
}

// CalculateFileHash 计算文件内容的哈希值
func CalculateFileHash(fileContent io.Reader) string {
	hash := sha256.New()
	if _, err := io.Copy(hash, fileContent); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// MapToStructTags 将 map[string]string 转换为结构体字段的 tag 标签内容
func MapToStructTags[T any](inputMap map[string]string, tag string) map[string]string {
	result := make(map[string]string)
	var t T
	val := reflect.ValueOf(t)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typeOfT := val.Type()

	for key, value := range inputMap {
		field, found := typeOfT.FieldByName(key)
		if !found {
			continue
		}
		tagValue := field.Tag.Get(tag)
		if tagValue != "" {
			result[tagValue] = value
		} else {
			result[key] = value
		}
	}

	return result
}
