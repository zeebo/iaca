// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package objfile implements portable access to OS-specific executable files.
package objfile

import (
	"fmt"
	"io"
	"os"
	"sort"
)

type rawFile interface {
	symbols() (syms []Sym, err error)
	text() (textStart, textOff uint64, text []byte, err error)
	goarch() string
}

// A File is an opened executable file.
type File struct {
	r   *os.File
	raw rawFile
}

// A Sym is a symbol defined in an executable file.
type Sym struct {
	Name string // symbol name
	Addr uint64 // virtual address of symbol
	Size int64  // size in bytes
	Code rune   // nm code (T for text, D for data, and so on)
	Type string // XXX?
}

var openers = []func(io.ReaderAt) (rawFile, error){
	openElf,
	openMacho,
}

// Open opens the named file.
// The caller must call f.Close when the file is no longer needed.
func Open(name string) (*File, error) {
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	if f, err := openGoFile(r); err == nil {
		return f, nil
	}
	for _, try := range openers {
		if raw, err := try(r); err == nil {
			return &File{r: r, raw: raw}, nil
		}
	}
	r.Close()
	return nil, fmt.Errorf("open %s: unrecognized object file", name)
}

func (f *File) Close() error {
	return f.r.Close()
}

func (f *File) Symbols() ([]Sym, error) {
	syms, err := f.raw.symbols()
	if err != nil {
		return nil, err
	}
	sort.Sort(byAddr(syms))
	return syms, nil
}

func (f *File) Text() (uint64, uint64, []byte, error) {
	return f.raw.text()
}

func (f *File) GOARCH() string {
	return f.raw.goarch()
}

type byAddr []Sym

func (x byAddr) Less(i, j int) bool { return x[i].Addr < x[j].Addr }
func (x byAddr) Len() int           { return len(x) }
func (x byAddr) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
