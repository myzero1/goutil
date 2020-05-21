package z1log

import (
	"fmt"
	"testing"
)

// debug->info->warn->error

func TestLog(t *testing.T) {
	fmt.Println("------")
	// set logPath,default ./logs
	SetLogPath("../logs")
	Debug("error %v", "Debug")
	Info("error %v", "info")
	Warn("error %v", "info")
	Error("error %v", "info")
}
