package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const cloudflareAPIURL = "https://api.cloudflare.com/client/v4"

// DNSRecord 表示一个DNS记录
type DNSRecord struct {
	ID      string `json:"ID"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListZonesResponse struct {
	Success bool          `json:"success"`
	Errors  []interface{} `json:"errors"`
	Result  []Zone        `json:"result"`
}

// CloudflareResponse 表示Cloudflare API的响应
type CloudflareResponse struct {
	Success bool            `json:"success"`
	Errors  []interface{}   `json:"errors"`
	Result  json.RawMessage `json:"result"`
}

// EmptyFunction 是一个空函数，用于在其他包中调用
func EmptyFunction() {
	// 这里什么都不做
}

// AddDNSRecord 添加DNS记录
func AddDNSRecord(apiKey, domain string, record DNSRecord) error {
	zoneID, err := GetZoneID(apiKey, domain)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/zones/%s/dns_records", cloudflareAPIURL, zoneID)
	return sendRequest("POST", url, apiKey, record)
}

// UpdateDNSRecord 更新DNS记录
func UpdateDNSRecord(apiKey, domain, recordID string, record DNSRecord) error {
	zoneID, err := GetZoneID(apiKey, domain)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", cloudflareAPIURL, zoneID, recordID)
	return sendRequest("PUT", url, apiKey, record)
}

// DeleteDNSRecord 删除DNS记录
func DeleteDNSRecord(apiKey, domain, recordID string) error {
	zoneID, err := GetZoneID(apiKey, domain)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", cloudflareAPIURL, zoneID, recordID)
	return sendRequest("DELETE", url, apiKey, nil)
}

// GetZoneID 获取zoneID
func GetZoneID(apiKey, domain string) (string, error) {
	url := fmt.Sprintf("%s/zones?name=%s", cloudflareAPIURL, domain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
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

	var listZonesResp ListZonesResponse
	if err := json.Unmarshal(body, &listZonesResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if !listZonesResp.Success {
		return "", fmt.Errorf("request failed: %v", listZonesResp.Errors)
	}

	if len(listZonesResp.Result) == 0 {
		return "", fmt.Errorf("no zones found for domain: %s", domain)
	}

	return listZonesResp.Result[0].ID, nil
}

// GetDNSRecordID 获取DNS记录ID
func GetDNSRecordID(apiKey, domain string, record DNSRecord) (string, error) {
	zoneID, err := GetZoneID(apiKey, domain)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/zones/%s/dns_records?type=%s&name=%s&content=%s", cloudflareAPIURL, zoneID, record.Type, record.Name, record.Content)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
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

	var cloudflareResp CloudflareResponse
	if err := json.Unmarshal(body, &cloudflareResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if !cloudflareResp.Success {
		return "", fmt.Errorf("API request failed: %v", cloudflareResp.Errors)
	}

	var dnsRecords []DNSRecord
	if err := json.Unmarshal(cloudflareResp.Result, &dnsRecords); err != nil {
		return "", fmt.Errorf("failed to unmarshal DNS records: %v", err)
	}

	if len(dnsRecords) == 0 {
		return "", nil
	}

	return dnsRecords[0].ID, nil
}

// sendRequest 发送HTTP请求
func sendRequest(method, url, apiKey string, data interface{}) error {
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

	req.Header.Set("Authorization", "Bearer "+apiKey)
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

	var cloudflareResp CloudflareResponse
	if err := json.Unmarshal(body, &cloudflareResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if !cloudflareResp.Success {
		return fmt.Errorf("request failed: %v", cloudflareResp.Errors)
	}

	return nil
}
