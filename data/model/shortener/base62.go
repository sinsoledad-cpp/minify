package shortener

import (
	"strings"
)

const (
	// Base62字符集: 0-9, a-z, A-Z
	base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	base        = 62
)

// ToBase62 将整数转换为Base62编码的字符串
func ToBase62(num int64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	var result strings.Builder
	n := num

	for n > 0 {
		result.WriteByte(base62Chars[n%base])
		n /= base
	}

	// 反转字符串
	return reverseString(result.String())
}

// FromBase62 将Base62编码的字符串转换回整数
func FromBase62(str string) int64 {
	var result int64
	for _, char := range str {
		result = result*base + int64(strings.IndexByte(base62Chars, byte(char)))
	}
	return result
}

// reverseString 反转字符串
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}