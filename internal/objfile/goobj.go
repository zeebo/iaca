// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Parsing of Go intermediate object files and archives.

package objfile

import (
	"fmt"
	"os"

	"github.com/zeebo/iaca/internal/goobj"
	"github.com/zeebo/iaca/internal/objabi"
)

type goobjFile struct {
	goobj *goobj.Package
	f     *os.File // the underlying .o or .a file
}

func openGoFile(r *os.File) (*File, error) {
	f, err := goobj.Parse(r, `""`)
	if err != nil {
		return nil, err
	}
	return &File{r: r, raw: &goobjFile{goobj: f, f: r}}, nil
}

func goobjName(id goobj.SymID) string {
	if id.Version == 0 {
		return id.Name
	}
	return fmt.Sprintf("%s<%d>", id.Name, id.Version)
}

func (f *goobjFile) symbols() ([]Sym, error) {
	seen := make(map[goobj.SymID]bool)

	var syms []Sym
	for _, s := range f.goobj.Syms {
		seen[s.SymID] = true
		sym := Sym{Addr: uint64(s.Data.Offset), Name: goobjName(s.SymID), Size: s.Size, Type: s.Type.Name, Code: '?'}
		switch s.Kind {
		case objabi.STEXT:
			sym.Code = 'T'
		case objabi.SRODATA:
			sym.Code = 'R'
		case objabi.SDATA:
			sym.Code = 'D'
		case objabi.SBSS, objabi.SNOPTRBSS, objabi.STLSBSS:
			sym.Code = 'B'
		}
		if s.Version != 0 {
			sym.Code += 'a' - 'A'
		}
		syms = append(syms, sym)
	}

	for _, s := range f.goobj.Syms {
		for _, r := range s.Reloc {
			if !seen[r.Sym] {
				seen[r.Sym] = true
				sym := Sym{Name: goobjName(r.Sym), Code: 'U'}
				if s.Version != 0 {
					// should not happen but handle anyway
					sym.Code = 'u'
				}
				syms = append(syms, sym)
			}
		}
	}

	return syms, nil
}

// We treat the whole object file as the text section.
func (f *goobjFile) text() (textStart, textOff uint64, text []byte, err error) {
	var info os.FileInfo
	info, err = f.f.Stat()
	if err != nil {
		return
	}
	text = make([]byte, info.Size())
	_, err = f.f.ReadAt(text, 0)
	return
}

func (f *goobjFile) goarch() string {
	return f.goobj.Arch
}
