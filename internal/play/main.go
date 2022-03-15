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
	"flag"
	"fmt"
	"os"
	"time"

	engine "github.com/bhojpur/render/pkg/app"
	"github.com/bhojpur/render/pkg/g3d/audio"
	"github.com/bhojpur/render/pkg/g3d/renderer"
)

// usage shows the application usage
func usage() {

	fmt.Fprintf(os.Stderr, "usage: renderplay <soundfile>\n")
}

func main() {

	// Parse command line parameters
	flag.Usage = usage
	flag.Parse()

	// Get file to play
	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}
	fpath := args[0]

	// Create a new 3D Rendering application
	engine.BhojpurApp3D()

	// Create player
	player, err := audio.NewPlayer(fpath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get total play time
	total := player.TotalTime()
	fmt.Printf("Bhojpur Render - playing audio:[%s] (%3.1f seconds)\n", fpath, total)

	// Start player
	player.Play()

	// Run the 3D Rendering application
	engine.BhojpurApp3D().Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {})
}
