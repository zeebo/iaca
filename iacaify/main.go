package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	"github.com/zeebo/errs"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

const (
	startMarker     = "\x45\x17\x0a\x2a\xcd\xe2\x52\x4a"
	endMarker       = "\xa0\x28\x3d\x2f\x8d\x21\xac\x47"
	iacaMarkerStart = "\xbb\x6f\x00\x00\x00\x64\x67\x90"
	iacaMarkerEnd   = "\xbb\xde\x00\x00\x00\x64\x67\x90"
)

func run() error {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		return errs.Wrap(err)
	}

	// Assumptions:
	// 1. There's a NOP close before the constant due to inlining
	// 2. Right after the instruction is 7 bytes moving into the global

	for data := data; ; {
		idx := bytes.Index(data, []byte(startMarker))
		if idx < 0 {
			break
		}
		for j := idx + 7; data[j] != 0x90; j-- {
			data[j] = 0x90
		}
		copy(data[idx+7:], iacaMarkerStart)
		data = data[idx:]
	}

	for data := data; ; {
		idx := bytes.Index(data, []byte(endMarker))
		if idx < 0 {
			break
		}
		j := idx
		for data[j] != 0x90 {
			j--
		}
		copy(data[j:], iacaMarkerEnd)
		for j += 8; j < idx+8+7; j++ {
			data[j] = 0x90
		}
		data = data[idx:]
	}

	_, err = os.Stdout.Write(data)
	return errs.Wrap(err)
}
