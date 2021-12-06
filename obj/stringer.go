// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

// This is a mini version of the stringer tool customized for the Anames table
// in the architecture support for obj.
// This version just generates the slice of strings, not the String method.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	input  = flag.String("i", "", "input file name")
	output = flag.String("o", "", "output file name")
	pkg    = flag.String("p", "", "package name")
)

var Are = regexp.MustCompile(`^\tA([A-Za-z0-9]+)`)

func main() {
	flag.Parse()
	if *input == "" || *output == "" || *pkg == "" {
		flag.Usage()
		os.Exit(2)
	}
	in, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	fd, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	out := bufio.NewWriter(fd)
	defer out.Flush()
	var on = false
	s := bufio.NewScanner(in)
	first := true
	for s.Scan() {
		line := s.Text()
		if !on {
			// First relevant line contains "= obj.ABase".
			// If we find it, delete the = so we don't stop immediately.
			const prefix = "= obj.ABase"
			index := strings.Index(line, prefix)
			if index < 0 {
				continue
			}
			// It's on. Start with the header.
			fmt.Fprintf(out, header, *input, *output, *pkg, *pkg)
			on = true
			line = line[:index]
		}
		// Strip comments so their text won't defeat our heuristic.
		index := strings.Index(line, "//")
		if index > 0 {
			line = line[:index]
		}
		index = strings.Index(line, "/*")
		if index > 0 {
			line = line[:index]
		}
		// Termination condition: Any line with an = changes the sequence,
		// so stop there, and stop at a closing brace.
		if strings.HasPrefix(line, "}") || strings.ContainsRune(line, '=') {
			break
		}
		sub := Are.FindStringSubmatch(line)
		if len(sub) < 2 {
			continue
		}
		if first {
			fmt.Fprintf(out, "\tobj.A_ARCHSPECIFIC: %q,\n", sub[1])
			first = false
		} else {
			fmt.Fprintf(out, "\t%q,\n", sub[1])
		}
	}
	fmt.Fprintln(out, "}")
	if s.Err() != nil {
		log.Fatal(err)
	}
}

const header = `// Code generated by stringer -i %s -o %s -p %s; DO NOT EDIT.

package %s

import "github.com/wdvxdr1123/golang-asm/obj"

var Anames = []string{
`
