package main

var x [4]uint64

func main() {
	x[0]++
	x[1]++
	x[2]++
	x[3]++
	x[0]++
	x[1]++
	x[2]++
	x[3]++
}
