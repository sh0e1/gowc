package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var dochar, doline, domulti, doword bool

func main() {
	flag.BoolVar(&dochar, "c", false, "The number of bytes in each input file is written to the standard output.  This will cancel out any prior usage of the -m option.")
	flag.BoolVar(&doline, "l", false, "The number of lines in each input file is written to the standard output.")
	flag.BoolVar(&domulti, "m", false, "The number of characters in each input file is written to the standard output.  If the current locale does not support multibyte characters, this is equivalent to the -c option.  This will cancel out any prior usage of the -c option.")
	flag.BoolVar(&doword, "w", false, "The number of words in each input file is written to the standard output.")
	flag.Parse()

	if flag.NFlag() == 0 {
		dochar, doline, doword = true, true, true
	}

	w := &bytes.Buffer{}
	for _, arg := range flag.Args() {
		f, err := os.Open(arg)
		if err != nil {
			fmt.Printf("gowc: %v\n", err)
			continue
		}
		// nolint:errcheck
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Printf("gowc: %v\n", err)
			continue
		}

		if doline {
			fmt.Fprintf(w, "%8d", count(bytes.NewReader(b), bufio.ScanLines))
		}
		if doword {
			fmt.Fprintf(w, "%8d", count(bytes.NewReader(b), bufio.ScanWords))
		}
		if domulti {
			dochar = false
			fmt.Fprintf(w, "%8d", count(bytes.NewReader(b), bufio.ScanRunes))
		}
		if dochar {
			fmt.Fprintf(w, "%8d", count(bytes.NewReader(b), bufio.ScanBytes))
		}

		fmt.Fprintf(w, " %s\n", arg)
	}

	fmt.Print(w.String())
}

func count(r io.Reader, sf bufio.SplitFunc) int {
	s := bufio.NewScanner(r)
	s.Split(sf)

	var cnt int
	for s.Scan() {
		cnt++
	}
	return cnt
}
