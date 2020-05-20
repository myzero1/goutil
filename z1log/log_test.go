package z1log

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	fmt.Println("------")
	Errorf("error %v", "error11")
	Infof("error %v", "info")
	fmt.Println("========")
}
