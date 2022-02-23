package utils

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
	"path/filepath"
	"sort"

	"github.com/bhojpur/render/pkg/3d/gui"
	"github.com/bhojpur/render/pkg/3d/gui/assets/icon"
	"github.com/bhojpur/render/pkg/3d/math32"
	"github.com/bhojpur/render/pkg/app"
)

type FileSelect struct {
	gui.Panel
	path *gui.Label
	list *gui.List
	bok  *gui.Button
	bcan *gui.Button
}

func NewFileSelect(width, height float32) (*FileSelect, error) {

	fs := new(FileSelect)
	fs.Panel.Initialize(fs, width, height)
	fs.SetBorders(2, 2, 2, 2)
	fs.SetPaddings(4, 4, 4, 4)
	fs.SetColor(math32.NewColor("White"))
	fs.SetVisible(false)
	fs.SetBounded(false)

	// Set vertical box layout for the whole panel
	l := gui.NewVBoxLayout()
	l.SetSpacing(4)
	fs.SetLayout(l)

	// Creates path label
	fs.path = gui.NewLabel("path")
	fs.Add(fs.path)

	// Creates list
	fs.list = gui.NewVList(0, 0)
	fs.list.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 5, AlignH: gui.AlignWidth})
	fs.list.Subscribe(gui.OnChange, func(evname string, ev interface{}) {
		fs.onSelect()
	})
	fs.Add(fs.list)

	// Button container panel
	bc := gui.NewPanel(0, 0)
	bcl := gui.NewHBoxLayout()
	bcl.SetAlignH(gui.AlignWidth)
	bc.SetLayout(bcl)
	bc.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 1, AlignH: gui.AlignWidth})
	fs.Add(bc)

	// Creates OK button
	fs.bok = gui.NewButton("OK")
	fs.bok.SetLayoutParams(&gui.HBoxLayoutParams{Expand: 0, AlignV: gui.AlignCenter})
	fs.bok.Subscribe(gui.OnClick, func(evname string, ev interface{}) {
		fs.Dispatch("OnOK", nil)
	})
	bc.Add(fs.bok)

	// Creates Cancel button
	fs.bcan = gui.NewButton("Cancel")
	fs.bcan.SetLayoutParams(&gui.HBoxLayoutParams{Expand: 0, AlignV: gui.AlignCenter})
	fs.bcan.Subscribe(gui.OnClick, func(evname string, ev interface{}) {
		fs.Dispatch("OnCancel", nil)
	})
	bc.Add(fs.bcan)

	// Sets initial directory
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	} else {
		fs.SetPath(path)
	}
	return fs, nil
}

// Show shows or hide the file selection dialog
func (fs *FileSelect) Show(show bool) {

	if show {
		fs.SetVisible(true)
		width, height := app.App().GetSize()
		px := (float32(width) - fs.Width()) / 2
		py := (float32(height) - fs.Height()) / 2
		fs.SetPosition(px, py)
	} else {
		fs.SetVisible(false)
	}
}

func (fs *FileSelect) SetPath(path string) error {

	// Open path file or dir
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Checks if it is a directory
	files, err := f.Readdir(0)
	if err != nil {
		return err
	}
	fs.path.SetText(path)

	// Sort files by name
	sort.Sort(listFileInfo(files))

	// Reads directory contents and loads into the list
	fs.list.Clear()
	// Adds previous directory
	prev := gui.NewImageLabel("..")
	prev.SetIcon(icon.FolderOpen)
	fs.list.Add(prev)
	// Adds directory files
	for i := 0; i < len(files); i++ {
		item := gui.NewImageLabel(files[i].Name())
		if files[i].IsDir() {
			item.SetIcon(icon.FolderOpen)
		} else {
			item.SetIcon(icon.InsertPhoto)
		}
		fs.list.Add(item)
	}
	return nil
}

func (fs *FileSelect) Selected() string {

	selist := fs.list.Selected()
	if len(selist) == 0 {
		return ""
	}
	label := selist[0].(*gui.ImageLabel)
	text := label.Text()
	return filepath.Join(fs.path.Text(), text)
}

func (fs *FileSelect) onSelect() {

	// Get selected image label and its txt
	sel := fs.list.Selected()[0]
	label := sel.(*gui.ImageLabel)
	text := label.Text()

	// Checks if previous directory
	if text == ".." {
		dir, _ := filepath.Split(fs.path.Text())
		fs.SetPath(filepath.Dir(dir))
		return
	}

	// Checks if it is a directory
	path := filepath.Join(fs.path.Text(), text)
	s, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	if s.IsDir() {
		fs.SetPath(path)
	}
}

// For sorting array of FileInfo by Name
type listFileInfo []os.FileInfo

func (fi listFileInfo) Len() int      { return len(fi) }
func (fi listFileInfo) Swap(i, j int) { fi[i], fi[j] = fi[j], fi[i] }
func (fi listFileInfo) Less(i, j int) bool {

	return fi[i].Name() < fi[j].Name()
}