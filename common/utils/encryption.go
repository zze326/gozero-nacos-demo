package utils

import (
	"crypto/md5"
	"encoding/hex"
)

/**
 * @Author: zze
 * @Date: 2022/5/19 17:05
 * @Desc: 加密相关
 */

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
