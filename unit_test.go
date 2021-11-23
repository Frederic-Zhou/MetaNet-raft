package main

import (
	"fmt"
	"testing"
)

func TestLanscan(t *testing.T) {
	// fmt.Println("开始")

	// var timer *time.Timer

	// timer = time.NewTimer(0 * time.Second)
	// timer.Reset(5 * time.Second)

	// select {
	// case <-timer.C:
	// case <-time.After(3 * time.Second):
	// 	// timer.Stop()
	// }

	// time.Sleep(1 * time.Second)
	// timer.Reset(5 * time.Second)
	// <-timer.C

	// fmt.Println("应该要等10s")
	for i := 0.0; i < 100; i++ {
		a := 1 + (i * 0.01)

		fmt.Println(a)
	}

	a := 1.14 * 1.0
	fmt.Println(a)

}
