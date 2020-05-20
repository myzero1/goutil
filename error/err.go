package error

import (
	"fmt"
	"runtime"
)

func CheckErr(err error, opts ...string) bool {
	if err != nil {
		msgPre := ""
		if len(opts) > 0 {
			msgPre = opts[0]
		}
		_, file, line, _ := runtime.Caller(1)
		fmt.Println("z1checkErr[", msgPre, "]", err, "(", line, file, ")")

		if len(opts) > 1 && opts[1] == "panic" {
			panic(err)
		}

		return true
	} else {
		return false
	}
}
