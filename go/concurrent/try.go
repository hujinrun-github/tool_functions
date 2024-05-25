package concurrent

import (
	"fmt"
	"runtime"
)

func Try(fun func()) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}

		buf := make([]byte, 16*1024*1024)
		buf = buf[:runtime.Stack(buf, false)]
		fmt.Printf("[PANIC]%v\n%s\n", err, buf)
	}()
	fun()
}
