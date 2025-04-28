package xploader

import (
	"testing"
)

func TestDefaultForegroundColor(t *testing.T) {
	black := Color{R: 0, G: 0, B: 0}
	if DefaultForegroundColor != black {
		t.Fatalf("DefaultForegroundColor is incorrect, got %#v", DefaultForegroundColor)
	}
}

func TestInvisibleColor(t *testing.T) {
	magenta := Color{R: 255, G: 0, B: 255}
	if InvisibleColor != magenta {
		t.Fatalf("InvisibleColor is incorrect, got %#v", InvisibleColor)
	}
}

func TestColorIsInvisible(t *testing.T) {
	invisible := InvisibleColor
	notInvisible := Color{R: 10, G: 10, B: 10}

	if !invisible.IsInvisible() {
		t.Fatal("Expected InvisibleColor to be invisible")
	}

	if notInvisible.IsInvisible() {
		t.Fatal("Expected non-magenta color to not be invisible")
	}
}

func TestNewEmptyCell(t *testing.T) {
	cell := NewEmptyCell()

	if cell.Rune != ' ' {
		t.Fatalf("Expected Rune ' ', got '%c'", cell.Rune)
	}
	if cell.Fg != DefaultForegroundColor {
		t.Fatalf("Expected foreground color %+v, got %+v", DefaultForegroundColor, cell.Fg)
	}
	if cell.Bg != InvisibleColor {
		t.Fatalf("Expected background color %+v, got %+v", InvisibleColor, cell.Bg)
	}
	if !cell.IsEmpty() {
		t.Fatal("Expected NewEmptyCell to be empty")
	}
}

func TestCellIsEmpty(t *testing.T) {
	cell := Cell{
		Rune: ' ',
		Fg:   DefaultForegroundColor,
		Bg:   InvisibleColor,
	}

	if !cell.IsEmpty() {
		t.Fatal("Expected cell to be empty")
	}

	cellNotEmpty := Cell{
		Rune: 'X',
		Fg:   DefaultForegroundColor,
		Bg:   InvisibleColor,
	}
	if cellNotEmpty.IsEmpty() {
		t.Fatal("Expected non-blank rune cell to not be empty")
	}
}

func TestNewEmptyLayer(t *testing.T) {
	width := 5
	height := 3
	layer := NewEmptyLayer(width, height)

	if int(layer.Width) != width {
		t.Fatalf("Expected width %d, got %d", width, layer.Width)
	}
	if int(layer.Height) != height {
		t.Fatalf("Expected height %d, got %d", height, layer.Height)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cell := layer.GetCell(x, y)
			if !cell.IsEmpty() {
				t.Fatalf("Expected cell at (%d,%d) to be empty", x, y)
			}
		}
	}
}
