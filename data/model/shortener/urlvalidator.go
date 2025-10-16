package shortener

import (
	"net/url"
	"strings"
)

// ValidateURL 验证URL是否有效
func ValidateURL(urlStr string) (bool, error) {
	// 检查URL是否为空
	if urlStr == "" {
		return false, nil
	}

	// 解析URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false, err
	}

	// 检查URL是否有效
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false, nil
	}

	// 检查主机名是否有效
	if parsedURL.Host == "" {
		return false, nil
	}

	return true, nil
}

// IsShortURL 检查URL是否是短链接
func IsShortURL(urlStr string, domain string) bool {
	// 如果URL为空，则不是短链接
	if urlStr == "" {
		return false
	}

	// 解析URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// 检查是否包含短链接域名
	return strings.Contains(parsedURL.Host, domain)
}
