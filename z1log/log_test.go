package z1log

import (
	"fmt"
	"testing"
)

// debug->info->warn->error

func TestLog(t *testing.T) {
	fmt.Println("------")

	Info("error %v", "info")
}
