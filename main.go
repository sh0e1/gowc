package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
)

var c, l, m, w bool

func main() {
	flag.BoolVar(&c, "c", false, "The number of bytes in each input file is written to the standard output.  This will cancel out any prior usage of the -m option.")
	flag.BoolVar(&l, "l", false, "The number of lines in each input file is written to the standard output.")
	flag.BoolVar(&m, "m", false, "The number of characters in each input file is written to the standard output.  If the current locale does not support multibyte characters, this is equivalent to the -c option.  This will cancel out any prior usage of the -c option.")
	flag.BoolVar(&w, "w", false, "The number of words in each input file is written to the standard output.")
	flag.Parse()

	buf := &bytes.Buffer{}
	for _, arg := range flag.Args() {
		f, err := os.Open(arg)
		if err != nil {
			fmt.Printf("gowc: %v\n", err)
			continue
		}
		// nolint:errcheck
		defer f.Close()

		var sf bufio.SplitFunc
		switch {
		case c:
			sf = bufio.ScanBytes
		case l:
			sf = bufio.ScanLines
		case m:
			sf = bufio.ScanRunes
		case w:
			sf = bufio.ScanWords
		}

		s := bufio.NewScanner(f)
		s.Split(sf)

		var cnt int
		for s.Scan() {
			cnt++
		}
		fmt.Fprintf(buf, "%8d %s\n", cnt, arg)
	}

	fmt.Print(buf.String())
}
