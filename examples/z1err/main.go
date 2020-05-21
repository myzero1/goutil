package main

import (
	"errors"
	"fmt"

	"github.com/myzero1/goutil/z1err"
)

func main() {
	fmt.Println("-----------------start")

	fmt.Println("------------------1----------------")
	err := errors.New("1 err")
	z1err.ChkErr(err)

	fmt.Println("------------------2----------------")
	z1err.ChkErr(errors.New("2 err"), "second err")

	for i := 0; i < 5; i++ {
		if i < 4 {
			if z1err.ChkErr(errors.New("test for continue"), "test for") {
				fmt.Println("=====", i)
				continue
			}
		}

		if z1err.ChkErr(errors.New("test for break"), "test for") {
			fmt.Println("=====break")
			break
		}
	}

	fmt.Println("------------------return----------------")
	if z1err.ChkErr(errors.New("test return"), "test return") {
		return
	}

	fmt.Println("------------------panic----------------")
	z1err.ChkErr(errors.New("test return"), "test panic", "panic")

	fmt.Println("-----------------end")
}
