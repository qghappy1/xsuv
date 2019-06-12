
package util

import (
	"regexp"
	"strconv"
	"strings"	
)


func SubStr(s string, start int, length int) string {
	rs := []rune(s)
	len := len(rs)
	if start < 0 {
		start = 0
	}
	if start + length>len {
		return ""
	}
	return string(rs[start:start+length])
}

func ToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i 
}

func ToString(i int) string {
	return strconv.Itoa(i)
}

func ToBool(s string) bool {
	s = strings.ToLower(s)
	return s == "true"
}

func CreateRandomString(num int) string{
   str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
   bytes := []byte(str)
   result := []byte{}
   for i := 0; i < num; i++ {
      result = append(result, bytes[Intn(len(bytes))])
   }
   return string(result)
}

// 字符串SQL注入检查，false:不合法 true:合法
func CheckSQLInject(str string) bool {
	express := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	match, err := regexp.Compile(express)
	if err != nil {
		return false
	}
	return !match.MatchString(str)
}