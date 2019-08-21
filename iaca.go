package main

import (
	"bufio"
	"encoding/hex"
	"log"
	"os"
	"strings"

	"github.com/zeebo/errs"
)

const (
	iacaStart = "\xbb\x6f\x00\x00\x00\x64\x67\x90"
	iacaEnd   = "\xbb\xde\x00\x00\x00\x64\x67\x90"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() (err error) {
	if _, err := os.Stdout.WriteString(iacaStart); err != nil {
		return errs.Wrap(err)
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		parts := strings.Fields(sc.Text())
		if len(parts) == 0 {
			continue
		} else if len(parts) < 3 {
			return errs.New("invalid line: %q", sc.Text())
		}
		if parts[0] == "TEXT" {
			continue
		}
		data, err := hex.DecodeString(parts[2])
		if err != nil {
			return errs.Wrap(err)
		}
		if _, err := os.Stdout.Write(data); err != nil {
			return errs.Wrap(err)
		}
	}
	if err := sc.Err(); err != nil {
		return errs.Wrap(err)
	}

	if _, err := os.Stdout.WriteString(iacaEnd); err != nil {
		return errs.Wrap(err)
	}

	return nil
}
