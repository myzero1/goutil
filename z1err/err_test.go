package z1err

import (
	"errors"
	"fmt"
	"testing"
)

func TestErr(t *testing.T) {
	fmt.Println("-----------------start")

	fmt.Println("------------------1----------------")
	err := errors.New("1 err")
	ChkErr(err)

	fmt.Println("------------------2----------------")
	ChkErr(errors.New("2 err"), "second err")

	for i := 0; i < 5; i++ {
		if i < 4 {
			if ChkErr(errors.New("test for continue"), "test for") {
				fmt.Println("=====", i)
				continue
			}
		}

		if ChkErr(errors.New("test for break"), "test for") {
			fmt.Println("=====break")
			break
		}
	}

	fmt.Println("------------------return----------------")
	if ChkErr(errors.New("test return"), "test return") {
		return
	}

	fmt.Println("------------------panic----------------")
	ChkErr(errors.New("test return"), "test panic", "panic")

	fmt.Println("-----------------end")
}
