# xploader

A [REXPaint](https://www.gridsagegames.com/rexpaint/) `.xp` file handler written in go.

## Overview
`xploader` allows you to load and manipulate REXPaint `.xp` files in
terminal-based or graphical projects.
- Fully supports both compressed (`gzip`) and uncompressed `.xp` files.
- Transparently handles REXPaint's unusual column-major format:
  - By default, data is reordered into row-major (line-by-line) ordering for
    easier use.
  - If desired, you can configure loader options to preserve column-major
    memory layout.
  - Regardless of storage layout, `Layer.GetCell(x, y)` always retrieves the
    expected cell at logical `(x,y)` coordinates.

Saving is supported:
- `XPFile` structs can be saved back to disk.
- If desired, you can configure saving options to save uncompressed files.
- Saved files are **100% compatible** with REXPaint (v1 format used by REXPaint
  1.70).

## Usage example
See [cmd/main.go](cmd/main.go)  for a fully functional example demonstrating:
- Loading a `.xp` file
- Displaying metadata
- Displaying its layers in a terminal (you can use the files in the `testdata`
  folder)
- Handling background and foreground colors properly

## Usage example using [tcell](https://github.com/gdamore/tcell)
```go
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/malc0mn/xploader"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <file.xp>", filepath.Base(os.Args[0]))
	}

	path := os.Args[1]

	// Load .xp file
	xpfile, err := xploader.LoadXPFile(path)
	if err != nil {
		log.Fatalf("Failed to load XP file: %v", err)
	}

	// Initialize tcell screen
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}
	if err = screen.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
	}
	defer screen.Fini()

	// Clear screen
	screen.Clear()

	// Draw the first layer of the XP file
	drawLayer(screen, &xpfile.Layers[0], 0, 0)

	// Show the result
	screen.Show()

	// Wait for a keypress before exiting
	waitForKeypress(screen)
}

// drawLayer renders a Layer at a given offset (originX, originY).
func drawLayer(screen tcell.Screen, layer *xploader.Layer, originX, originY int) {
	height := int(layer.Height)
	width := int(layer.Width)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cell := layer.GetCell(x, y)

			style := tcell.StyleDefault
			if !cell.Fg.IsInvisible() {
				style = style.Foreground(tcell.NewRGBColor(int32(cell.Fg.R), int32(cell.Fg.G), int32(cell.Fg.B)))
			}
			if !cell.Bg.IsInvisible() {
				style = style.Background(tcell.NewRGBColor(int32(cell.Bg.R), int32(cell.Bg.G), int32(cell.Bg.B)))
			}

			screen.SetContent(originX+x, originY+y, cell.Rune, nil, style)
		}
	}
}

// waitForKeypress blocks until a key event is received.
func waitForKeypress(screen tcell.Screen) {
	for {
		switch screen.PollEvent().(type) {
		case *tcell.EventKey:
			return
		case *tcell.EventResize:
			screen.Sync()
		}
	}
}
```