package error

import (
	"errors"
	"fmt"
	"testing"
)

func TestErr(t *testing.T) {
	fmt.Println("-----------------start")

	fmt.Println("------------------1----------------")
	err := errors.New("1 err")
	CheckErr(err)

	fmt.Println("------------------2----------------")
	CheckErr(errors.New("2 err"), "second err")

	for i := 0; i < 5; i++ {
		if i < 4 {
			if CheckErr(errors.New("test for continue"), "test for") {
				fmt.Println("=====", i)
				continue
			}
		}

		if CheckErr(errors.New("test for break"), "test for") {
			fmt.Println("=====break")
			break
		}
	}

	fmt.Println("------------------return----------------")
	if CheckErr(errors.New("test return"), "test return") {
		return
	}

	fmt.Println("------------------panic----------------")
	CheckErr(errors.New("test return"), "test panic", "panic")

	fmt.Println("-----------------end")
}
