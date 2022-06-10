package utils

import (
	"fmt"
	"strings"
)

/**
 * @Author: zze
 * @Date: 2022/5/27 13:34
 * @Desc: 字符串工具方法
 */

// RemoveHttpPrefix
// desc: 移除 url 中的 http 和 https 前缀
func RemoveHttpPrefix(url string) string {
	url = strings.TrimSpace(url)
	if len(url) > 7 && url[:7] == "http://" {
		return url[7:]
	}

	if len(url) > 8 && url[:8] == "https://" {
		return url[8:]
	}

	return url
}

// IsEmpty
// desc: 判断字符串是否为空
func IsEmpty(str string) bool {
	if len(strings.TrimSpace(fmt.Sprintf("%v", str))) == 0 {
		return true
	}
	return false
}

func ToString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}
