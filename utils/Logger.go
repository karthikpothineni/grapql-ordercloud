package utils

import (
	"fmt"
)

func Info(str string) {
	GetLog().Info(str)
}
func Error(str string) {
	GetLog().Error(str)
}
func Logger(levelOptional string, _str ...any) {
	str := ""
	for i := 0; i < len(_str); i++ {
		s := fmt.Sprint(_str[i])
		fmt.Println(_str[i])
		str += s
	}
	if str != "" {
		m := map[string]interface{}{
			"i": Info,
			"e": Error,
		}
		m[levelOptional].(func(string))(str)
	}
}
