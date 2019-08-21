package main

import (
	"io"
	"log"
	"os"

	"github.com/zeebo/errs"
	"github.com/zeebo/iaca/internal/objfile"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	f, err := objfile.Open(os.Args[1])
	if err != nil {
		return errs.Wrap(err)
	}
	anns, err := f.Annotations()
	if err != nil {
		return errs.Wrap(err)
	}
	f.Close()

	fh, err := os.Open(os.Args[1])
	if err != nil {
		return errs.Wrap(err)
	}
	defer fh.Close()

	copied := uint64(0)
	for _, ann := range anns {
		_, err := io.CopyN(os.Stdout, fh, int64(ann.Address-copied))
		if err != nil {
			return errs.Wrap(err)
		}

		if ann.Start {
			_, err = os.Stdout.WriteString("\xbb\x6f\x00\x00\x00\x64\x67\x90")
		} else {
			_, err = os.Stdout.WriteString("\xbb\xde\x00\x00\x00\x64\x67\x90")
		}
		if err != nil {
			return errs.Wrap(err)
		}

		if _, err := io.ReadFull(fh, make([]byte, 8)); err != nil {
			return errs.Wrap(err)
		}

		copied = ann.Address + 8
	}

	_, err = io.Copy(os.Stdout, fh)
	return errs.Wrap(err)
}
