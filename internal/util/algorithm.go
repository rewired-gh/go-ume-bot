package util

import "fmt"

func GetNaturalSet(n int64) (str string) {
	str = ""
	for i := int64(1); i <= n; i++ {
		if i == 1 {
			str = "{}"
		} else {
			str = fmt.Sprintf("%s, {%s}", str, str)
		}
	}
	str = fmt.Sprintf("{%s}", str)
	return
}
