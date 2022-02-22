package text

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

// StrCount returns the number of runes in the specified string
func StrCount(s string) int {

	count := 0
	for range s {
		count++
	}
	return count
}

// StrFind returns the start and length of the rune at the
// specified position in the string
func StrFind(s string, pos int) (start, length int) {

	count := 0
	for index := range s {
		if count == pos {
			start = index
			count++
			continue
		}
		if count == pos+1 {
			length = index - start
			break
		}
		count++
	}
	if length == 0 {
		length = len(s) - start
	}
	return start, length
}

// StrRemove removes the rune from the specified string and position
func StrRemove(s string, col int) string {

	start, length := StrFind(s, col)
	return s[:start] + s[start+length:]
}

// StrInsert inserts a string at the specified character position
func StrInsert(s, data string, col int) string {

	start, _ := StrFind(s, col)
	return s[:start] + data + s[start:]
}

// StrPrefix returns the prefix of the specified string up to
// the specified character position
func StrPrefix(text string, pos int) string {

	count := 0
	for index := range text {
		if count == pos {
			return text[:index]
		}
		count++
	}
	return text
}
