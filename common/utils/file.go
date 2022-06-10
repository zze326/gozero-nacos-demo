package utils

import "os"

/**
 * @Author: zze
 * @Date: 2022/6/10 11:24
 * @Desc: 文件操作
 */

func FilePathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

func AbsentDir(path string) (bool, error) {
	if FilePathExist(path) {
		err := os.RemoveAll(path)
		if err != nil {
			return false, err
		}
		err = os.Remove(path)
		if err != nil {
			if !FilePathExist(path) {
				return true, nil
			}
			return false, err
		}

	}
	return true, nil
}
