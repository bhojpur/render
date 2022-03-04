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
	"net"
)

// Net is a network writer used for logging.
type Net struct {
	conn net.Conn
}

// NewNet creates and returns a pointer to a new Net object along with any error that occurred.
func NewNet(network string, address string) (*Net, error) {

	n := new(Net)
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	n.conn = conn
	return n, nil
}

// Write writes the provided logger event to the network.
func (n *Net) Write(event *Event) {

	n.conn.Write([]byte(event.fmsg))
}

// Clone closes the network connection.
func (n *Net) Close() {

	n.conn.Close()
}

func (n *Net) Sync() {

}
