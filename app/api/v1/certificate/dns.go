package certificate

import (
	"fmt"
	"nav-web-site/app/aliyun"
	"nav-web-site/app/cloudflare"
	"strings"
)

func CallEmptyFunctions() {
	aliyun.EmptyFunction()
	cloudflare.EmptyFunction()
}

func AddDNSRecord(domain string) error {
	dnsProvider, err := getDNSProvider(domain)
	if err != nil {
		return fmt.Errorf("failed to get DNS provider: %v", err)
	}

	switch dnsProvider {
	case "cloudflare":
		apiKey := ""
		record := cloudflare.DNSRecord{
			Type:    "",
			Name:    "",
			Content: "",
			TTL:     60,
		} // 初始化 record 变量

		// 检查DNS记录是否存在
		existingRecordID, err := cloudflare.GetDNSRecordID(apiKey, domain, record)
		if err != nil {
			return fmt.Errorf("failed to check DNS record: %v", err)
		}

		if existingRecordID != "" {
			// 更新DNS记录
			return cloudflare.UpdateDNSRecord(apiKey, domain, existingRecordID, record)
		} else {
			// 添加DNS记录
			return cloudflare.AddDNSRecord(apiKey, domain, record)
		}
	case "aliyun":
		accessKeyId := ""
		accessKeySecret := ""
		record := aliyun.DNSRecord{
			Type:    "",
			Name:    "",
			Content: "",
			TTL:     60,
		}

		// 检查DNS记录是否存在
		existingRecordID, err := aliyun.GetDNSRecordID(accessKeyId, accessKeySecret, domain, record)
		if err != nil {
			return fmt.Errorf("failed to check DNS record: %v", err)
		}

		if existingRecordID != "" {
			// 更新DNS记录
			return aliyun.UpdateDNSRecord(accessKeyId, accessKeySecret, existingRecordID, record)
		} else {
			// 添加DNS记录
			return aliyun.AddDNSRecord(accessKeyId, accessKeySecret, domain, record)
		}
	default:
		return fmt.Errorf("unsupported DNS provider: %s", dnsProvider)
	}
}

func getDNSProvider(domain string) (string, error) {
	// 这里实现检测域名 DNS 服务商的逻辑
	// 例如，通过 DNS 查询或 API 调用来确定服务商
	// 返回 "cloudflare" 或 "aliyun" 等
	// 这是一个示例实现，具体实现需要根据实际情况编写
	if strings.HasSuffix(domain, ".cf") {
		return "cloudflare", nil
	} else if strings.HasSuffix(domain, ".ali") {
		return "aliyun", nil
	}
	return "", fmt.Errorf("unknown DNS provider for domain: %s", domain)
}
