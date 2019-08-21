package main

import "github.com/zeebo/iaca"

var x [4]uint64

func main() {
	iaca.Start()
	x[0]++
	x[1]++
	x[2]++
	x[3]++
	x[0]++
	x[1]++
	x[2]++
	x[3]++
	iaca.Stop()
}
