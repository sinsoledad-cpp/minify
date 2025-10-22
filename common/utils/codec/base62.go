package codec

import (
	"strings"
)

// Base62 字符集 (0-9, a-z, A-Z)
const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base62Length = uint64(len(base62Chars))

// Base62Encode 将 uint64 转换为 Base62 字符串
func Base62Encode(number uint64) string {
	if number == 0 {
		return string(base62Chars[0])
	}
	var result strings.Builder
	result.Grow(11) // 64位无符号整数最大值转为62进制，长度不超过11
	for number > 0 {
		remainder := number % base62Length
		result.WriteByte(base62Chars[remainder])
		number /= base62Length
	}
	return reverseString(result.String())
}

// reverseString 反转字符串
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
