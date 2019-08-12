package recovery

import (
	"fmt"
	"github.com/selinplus/go-dingtalk/pkg/logging"
)

func Recovery() {
	if r := recover(); r != nil {
		logging.Info(fmt.Sprintf("recovered:", r))
	}
}
