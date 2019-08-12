package recovery

import (
	"fmt"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"runtime/debug"
)

func Recovery() bool {
	if r := recover(); r != nil {
		logging.Info(fmt.Sprintf("recovered:%s", string(debug.Stack())))
		return true
	}
	return false
}
