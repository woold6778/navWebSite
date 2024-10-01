package aliyun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const aliyunAPIURL = "https://alidns.aliyuncs.com"

// DNSRecord 表示一个DNS记录
type DNSRecord struct {
	Type    string `json:"Type"`
	Name    string `json:"RR"`
	Content string `json:"Value"`
	TTL     int    `json:"TTL"`
}

// AliyunResponse 表示阿里云API的响应
type AliyunResponse struct {
	RequestId string `json:"RequestId"`
	RecordId  string `json:"RecordId"`
	Code      string `json:"Code"`
	Message   string `json:"Message"`
}

// EmptyFunction 是一个空函数，用于在其他包中调用
func EmptyFunction() {
	// 这里什么都不做
}

// AddDNSRecord 添加DNS记录
func AddDNSRecord(accessKeyId, accessKeySecret, domain string, record DNSRecord) error {
	url := fmt.Sprintf("%s/?Action=AddDomainRecord&DomainName=%s&RR=%s&Type=%s&Value=%s&TTL=%d", aliyunAPIURL, domain, record.Name, record.Type, record.Content, record.TTL)
	return sendRequest("GET", url, accessKeyId, accessKeySecret, nil)
}

// UpdateDNSRecord 更新DNS记录
func UpdateDNSRecord(accessKeyId, accessKeySecret, recordId string, record DNSRecord) error {
	url := fmt.Sprintf("%s/?Action=UpdateDomainRecord&RecordId=%s&RR=%s&Type=%s&Value=%s&TTL=%d", aliyunAPIURL, recordId, record.Name, record.Type, record.Content, record.TTL)
	return sendRequest("GET", url, accessKeyId, accessKeySecret, nil)
}

// GetDNSRecordID 获取DNS记录ID
func GetDNSRecordID(accessKeyId, accessKeySecret, domain string, record DNSRecord) (string, error) {
	url := fmt.Sprintf("%s/?Action=DescribeDomainRecords&DomainName=%s&RRKeyWord=%s&Type=%s", aliyunAPIURL, domain, record.Name, record.Type)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "APPCODE "+accessKeyId+":"+accessKeySecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var aliyunResp struct {
		DomainRecords struct {
			Record []struct {
				RecordId string `json:"RecordId"`
			} `json:"Record"`
		} `json:"DomainRecords"`
	}
	if err := json.Unmarshal(body, &aliyunResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(aliyunResp.DomainRecords.Record) == 0 {
		return "", nil
	}

	return aliyunResp.DomainRecords.Record[0].RecordId, nil
}

// sendRequest 发送HTTP请求
func sendRequest(method, url, accessKeyId, accessKeySecret string, data interface{}) error {
	client := &http.Client{}
	var req *http.Request
	var err error

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %v", err)
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "APPCODE "+accessKeyId+":"+accessKeySecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var aliyunResp AliyunResponse
	if err := json.Unmarshal(body, &aliyunResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if aliyunResp.Code != "" {
		return fmt.Errorf("request failed: %v", aliyunResp.Message)
	}

	return nil
}
