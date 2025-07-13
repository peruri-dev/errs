package errs

import (
	"fmt"
	"runtime"
)

func getWithDept(depth int) string {
	pc, _, line, _ := runtime.Caller(depth)
	return fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), line)
}
