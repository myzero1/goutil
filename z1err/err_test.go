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

	err = errors.New("2 err")
	ChkErr(err, "second err")

	for i := 0; i < 5; i++ {
		if i < 4 {
			err := errors.New("test for continue")
			if ChkErr(err, "test for") {
				fmt.Println("=====", i)
				continue
			}
		}

		err := errors.New("test for break")
		if ChkErr(err, "test for") {
			fmt.Println("=====break")
			break
		}
	}

	fmt.Println("------------------return----------------")

	err = errors.New("test return")
	if ChkErr(err, "test return") {
		return
	}

	fmt.Println("------------------panic----------------")

	err = errors.New("test return")
	ChkErr(err, "test panic", "panic")

	fmt.Println("-----------------end")
}
