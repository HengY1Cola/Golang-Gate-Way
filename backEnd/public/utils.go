package public

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
)

// GenSaltPassword sha256加密
func GenSaltPassword(salt string, pass string) string {
	s1 := sha256.New()
	s1.Write([]byte(pass))
	str1 := fmt.Sprintf("%x", s1.Sum(nil))
	s2 := sha256.New()
	s2.Write([]byte(str1 + salt))
	return fmt.Sprintf("%x", s2.Sum(nil))
}

//MD5 md5加密
func MD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Obj2Json 定义json转换的方法
func Obj2Json(s interface{}) string {
	marshal, _ := json.Marshal(s)
	return string(marshal)
}

// InStringSlice 判断在不在里面
func InStringSlice(slice []string, str string) bool {
	for _, item := range slice {
		if str == item {
			return true
		}
	}
	return false
}
