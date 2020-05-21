package main

import "github.com/myzero1/goutil/z1log"

func main() {
	z1log.SetLogPath("../logs")
	z1log.Debug("It just a test for debug")
	z1log.Debugf("It just a test for debugf(%v)", 1234)
}
