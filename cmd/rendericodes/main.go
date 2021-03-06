package main

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
	"unicode"
)

// Command line options
var (
	oPackage = flag.String("pkg", "icon", "Package name")
)

// Program name and version
const (
	PROGNAME = "rendericodes"
	VMAJOR   = 0
	VMINOR   = 1
)

type constInfo struct {
	Name  string
	Value string
}

type templateData struct {
	Packname string
	Consts   []constInfo
}

func main() {

	// Parse command line parameters
	flag.Usage = usage
	flag.Parse()

	// Opens input file
	if len(flag.Args()) == 0 {
		log.Fatal("Input file not supplied")
		return
	}
	finput, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
		return
	}

	// Creates optional output file
	fout := os.Stdout
	if len(flag.Args()) > 1 {
		fout, err = os.Create(flag.Args()[1])
		if err != nil {
			log.Fatal(err)
			return
		}
		defer fout.Close()
	}

	// Parse input file
	var td templateData
	td.Packname = *oPackage
	err = parse(finput, &td)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Parses the template
	tmpl := template.New("templ")
	tmpl, err = tmpl.Parse(templText)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Expands template to buffer
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, &td)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Formats buffer as Go source code
	p, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
		return
	}

	// Writes formatted source to output file
	fout.Write(p)
}

func parse(fin io.Reader, td *templateData) error {

	// Read words from input reader and builds words map
	scanner := bufio.NewScanner(fin)
	for scanner.Scan() {
		// Read next line
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			return err
		}
		// Remove line terminator, spaces and ignore empty lines
		line = strings.Trim(line, "\n ")
		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			continue
		}
		name := parts[0]
		code := parts[1]
		nameParts := strings.Split(name, "_")
		for i := 0; i < len(nameParts); i++ {
			nameParts[i] = strings.Title(nameParts[i])
		}
		finalName := strings.Join(nameParts, "")
		// If name starts with number adds prefix
		runes := []rune(finalName)
		if unicode.IsDigit(runes[0]) {
			finalName = "N" + finalName
		}
		td.Consts = append(td.Consts, constInfo{Name: finalName, Value: "0x" + code})
	}
	return nil
}

// Shows application usage
func usage() {

	fmt.Fprintf(os.Stderr, "%s v%d.%d\n", PROGNAME, VMAJOR, VMINOR)
	fmt.Fprintf(os.Stderr, "usage: %s [options] <input file> <output file>\n", strings.ToLower(PROGNAME))
	flag.PrintDefaults()
	os.Exit(0)
}

const templText = `// Code generated by Bhojpur Render - 3D Graphics Engine ('rendericodes'). Do not edit.
// This file is based on the original 'codepoints' file from the material design icon fonts:
// https://github.com/google/material-design-icons

package {{.Packname}}

// Icon constants.
const (
	{{range .Consts}}
		{{.Name}} = string({{.Value}})
	{{- end}}
)

// Codepoint returns the codepoint for the specified icon name.
// Returns 0 if the name not found
func Codepoint(name string) string {

	return name2Codepoint[name]
}

// Maps icon name to codepoint
var name2Codepoint = map[string]string{
	{{range .Consts}}
		"{{.Name}}": string({{.Value}}),	
	{{- end}}
}
`
