package logger

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
	"os"
)

// ANSI terminal color codes
const (
	csi      = "\x1B["
	black    = "30m"
	red      = "31m"
	green    = "32m"
	yellow   = "33m"
	blue     = "34m"
	magenta  = "35m"
	cyan     = "36m"
	white    = "37m"
	bblack   = "30m"
	bred     = "31;1m"
	bgreen   = "32;1m"
	byellow  = "33;1m"
	bblue    = "34;1m"
	bmagenta = "35;1m"
	bcyan    = "36;1m"
	bwhite   = "37;1m"
)

// Maps log level to color sequence
var colorMap = map[int]string{
	DEBUG: white,
	INFO:  green,
	WARN:  byellow,
	ERROR: bred,
	FATAL: bmagenta,
}

// Console is a console writer used for logging.
type Console struct {
	writer *os.File
	color  bool
}

// NewConsole creates and returns a new logger Console writer
// If color is true, this writer uses Ansi codes to write
// log messages in color accordingly to its level.
func NewConsole(color bool) *Console {

	return &Console{os.Stdout, color}
}

// Write writes the provided logger event to the console.
func (w *Console) Write(event *Event) {

	if w.color {
		w.writer.Write([]byte(csi))
		w.writer.Write([]byte(colorMap[event.level]))
	}
	w.writer.Write([]byte(event.fmsg))
	if w.color {
		w.writer.Write([]byte(csi))
		w.writer.Write([]byte(white))
	}
}

func (w *Console) Close() {

}

func (w *Console) Sync() {

}
