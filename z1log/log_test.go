package z1log

import (
	"fmt"
	"testing"
)

// debug->info->warn->error

func TestLog(t *testing.T) {
	fmt.Println("------")

	Debug("error %v", "Debug")
	Info("error %v", "info")
}
