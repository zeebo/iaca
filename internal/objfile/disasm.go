// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package objfile

import (
	"fmt"
	"sort"

	"golang.org/x/arch/x86/x86asm"
)

type Annotation struct {
	Address uint64
	Start   bool
}

func (f *File) Annotations() ([]Annotation, error) {
	if f.GOARCH() != "amd64" {
		return nil, fmt.Errorf("unsupported architecture")
	}

	syms, err := f.Symbols()
	if err != nil {
		return nil, err
	}

	lookup := func(addr uint64) (name string, base uint64) {
		i := sort.Search(len(syms), func(i int) bool { return addr < syms[i].Addr })
		if i > 0 {
			s := syms[i-1]
			if s.Addr != 0 && s.Addr <= addr && addr < s.Addr+uint64(s.Size) {
				return s.Name, s.Addr
			}
		}
		return "", 0
	}

	textStart, textOff, text, err := f.Text()
	if err != nil {
		return nil, err
	}
	textEnd := textStart + uint64(len(text))

	var anns []Annotation

	for _, sym := range syms {
		symStart := sym.Addr
		symEnd := sym.Addr + uint64(sym.Size)

		if sym.Code != 'T' && sym.Code != 't' || symStart < textStart {
			continue
		}
		if symEnd > textEnd {
			symEnd = textEnd
		}

		code := text[:symEnd-textStart]

		for pc := symStart; pc < symEnd; {
			i := pc - textStart
			text, size := disasm_amd64(code[i:], pc, lookup)
			switch text {
			case "CALL github.com/zeebo/iaca.padStart(SB)":
				anns = append(anns, Annotation{
					Address: pc - textStart + textOff - 3,
					Start:   true,
				})

			case "CALL github.com/zeebo/iaca.padStop(SB)":
				anns = append(anns, Annotation{
					Address: pc - textStart + textOff - 3,
					Start:   false,
				})
			}
			pc += uint64(size)
		}
	}

	sort.Slice(anns, func(i, j int) bool {
		return anns[i].Address < anns[j].Address
	})

	return anns, nil
}

type lookupFunc = func(addr uint64) (sym string, base uint64)

func disasm_amd64(code []byte, pc uint64, lookup lookupFunc) (string, int) {
	inst, err := x86asm.Decode(code, 64)
	var text string
	size := inst.Len
	if err != nil || size == 0 || inst.Op == 0 {
		size = 1
		text = "?"
	} else {
		text = x86asm.GoSyntax(inst, pc, lookup)
	}
	return text, size
}
