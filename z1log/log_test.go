package z1log

import (
	"fmt"
	"testing"
)

// debug->info->warn->error

func TestLog(t *testing.T) {
	fmt.Println("------")
	SetLogPath("./logs")
	SetOutput("console")
	SetJsonFormat(false)
	Debug("error %v", "Debug")
	Info("error %v", "info")
	Warn("error %v", "info")
	Error("error %v", "info")
}
