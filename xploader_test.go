package xploader

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

const testDataDir = "testdata/"

func newXpFile(layer Layer, t *testing.T) *XPFile {
	t.Helper()

	return &XPFile{
		Version: -1,
		Layers: []Layer{
			layer,
		},
	}
}

func newSimpleLayer(t *testing.T) *Layer {
	t.Helper()

	simple := NewEmptyLayer(10, 15)
	simple.Cells[0][0] = Cell{
		Rune: 'x',
		Fg:   Color{R: 255},
		Bg:   Color{G: 128, B: 255},
	}
	simple.Cells[0][1] = Cell{
		Rune: 'p',
		Fg:   Color{R: 255, G: 255},
		Bg:   Color{R: 191, B: 255},
	}
	simple.Cells[0][2] = Cell{
		Rune: 'l',
		Fg:   Color{R: 128, G: 255},
		Bg:   Color{R: 255, B: 191},
	}
	simple.Cells[0][3] = Cell{
		Rune: 'o',
		Fg:   Color{G: 255},
		Bg:   Color{R: 255, B: 128},
	}
	simple.Cells[0][4] = Cell{
		Rune: 'a',
		Fg:   Color{G: 255, B: 128},
		Bg:   Color{R: 255, B: 64},
	}
	simple.Cells[0][5] = Cell{
		Rune: 'd',
		Fg:   Color{G: 255, B: 191},
		Bg:   Color{R: 158, G: 158, B: 158},
	}
	simple.Cells[0][6] = Cell{
		Rune: 'e',
		Fg:   Color{G: 255, B: 255},
		Bg:   Color{R: 158, G: 134, B: 100},
	}
	simple.Cells[0][7] = Cell{
		Rune: 'r',
		Fg:   Color{G: 191, B: 255},
		Bg:   Color{R: 255, G: 255, B: 255},
	}

	return simple
}

func assertXPFileEqual(expected, actual *XPFile, t *testing.T) {
	t.Helper()

	if expected.Version != actual.Version {
		t.Fatalf("Expected version %d, got %d", expected.Version, actual.Version)
	}

	if len(expected.Layers) != len(actual.Layers) {
		t.Fatalf("Expected %d layers, got %d", len(expected.Layers), len(actual.Layers))
	}

	for i, expLayer := range expected.Layers {
		actLayer := actual.Layers[i]

		if expLayer.Width != actLayer.Width || expLayer.Height != actLayer.Height {
			t.Fatalf("Layer %d: expected dimensions %dx%d, got %dx%d", i, expLayer.Width, expLayer.Height, actLayer.Width, actLayer.Height)
		}

		if expLayer.ColumnMajor != actLayer.ColumnMajor {
			t.Fatalf("Layer %d: expected ColumnMajor=%v, got ColumnMajor=%v", i, expLayer.ColumnMajor, actLayer.ColumnMajor)
		}

		height := int(expLayer.Height)
		width := int(expLayer.Width)

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				expCell := expLayer.GetCell(x, y)
				actCell := actLayer.GetCell(x, y)

				if expCell.Rune != actCell.Rune {
					t.Errorf("Layer %d, Cell (%d,%d): expected rune '%c', got '%c'", i, x, y, expCell.Rune, actCell.Rune)
				}
				if expCell.Fg != actCell.Fg {
					t.Errorf("Layer %d, Cell (%d,%d): expected FG %+v, got %+v", i, x, y, expCell.Fg, actCell.Fg)
				}
				if expCell.Bg != actCell.Bg {
					t.Errorf("Layer %d, Cell (%d,%d): expected BG %+v, got %+v", i, x, y, expCell.Bg, actCell.Bg)
				}
			}
		}
	}
}

func TestLoadSimpleFiles(t *testing.T) {
	simple := newSimpleLayer(t)

	tests := []struct {
		name     string
		filename string
		expected *Layer
	}{
		{name: "Empty", filename: "empty.xp", expected: NewEmptyLayer(60, 60)},
		{name: "SimpleCompressed", filename: "simple.xp", expected: simple},
		{name: "SimpleUncompressed", filename: "simple_plain.xp", expected: simple},
	}

	for _, tt := range tests {
		tt := tt // MUST copy here to prevent awkward results due to reuse
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := testDataDir + tt.filename
			xp, err := LoadXPFile(f)
			if err != nil {
				t.Fatalf("Failed to load %q: %v", f, err)
			}

			assertXPFileEqual(newXpFile(*tt.expected, t), xp, t)
		})
	}
}

func TestLoadMultiLayerFile(t *testing.T) {
	xp, err := LoadXPFile(testDataDir + "multilayer.xp")
	if err != nil {
		t.Fatalf("Failed to load multilayer XP file: %v", err)
	}

	if len(xp.Layers) < 2 {
		t.Fatalf("Expected multiple layers, got %d", len(xp.Layers))
	}

	expected := newXpFile(*newSimpleLayer(t), t)

	layerTwo := NewEmptyLayer(10, 15)
	layerTwo.Cells[14][0] = Cell{
		Rune: 'E',
		Fg:   Color{R: 255, G: 255, B: 255},
		Bg:   Color{},
	}
	layerTwo.Cells[14][1] = Cell{
		Rune: 'X',
		Fg:   Color{R: 255, G: 255, B: 255},
		Bg:   Color{},
	}
	layerTwo.Cells[14][2] = Cell{
		Rune: 'P',
		Fg:   Color{R: 255, G: 255, B: 255},
		Bg:   Color{},
	}
	layerTwo.Cells[14][3] = Cell{
		Rune: 'L',
		Fg:   Color{R: 255, G: 255, B: 255},
		Bg:   Color{},
	}
	layerTwo.Cells[14][4] = Cell{
		Rune: 'O',
		Fg:   Color{R: 255, G: 255, B: 255},
		Bg:   Color{},
	}
	layerTwo.Cells[14][5] = Cell{
		Rune: 'A',
		Fg:   Color{R: 255, G: 255, B: 255},
		Bg:   Color{},
	}
	layerTwo.Cells[14][6] = Cell{
		Rune: 'D',
		Fg:   Color{R: 255, G: 255, B: 255},
		Bg:   Color{},
	}
	layerTwo.Cells[14][7] = Cell{
		Rune: 'E',
		Fg:   Color{R: 255, G: 255, B: 255},
		Bg:   Color{},
	}
	layerTwo.Cells[14][8] = Cell{
		Rune: 'R',
		Fg:   Color{R: 255, G: 255, B: 255},
		Bg:   Color{},
	}

	expected.AddLayer(*layerTwo)

	assertXPFileEqual(expected, xp, t)
}

func TestLoadBrokenFile(t *testing.T) {
	_, err := LoadXPFile(testDataDir + "broken.xp")
	if err == nil {
		t.Fatal("Expected error when loading corrupted XP file, got nil")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := LoadXPFile(testDataDir + "non_existent.xp")
	if err == nil {
		t.Fatal("Expected error when loading non-existent XP file, got nil")
	}
}

func TestSaveXPFile(t *testing.T) {
	sourceFile := filepath.Join(testDataDir, "simple.xp")

	// Load original XP file.
	xpOriginal, err := LoadXPFile(sourceFile)
	if err != nil {
		t.Fatalf("Failed to load original XP file: %v", err)
	}

	// Save to temporary file.
	tempDir := t.TempDir()
	savedFile := filepath.Join(tempDir, "simple_saved.xp")
	err = SaveXPFile(xpOriginal, savedFile)
	if err != nil {
		t.Fatalf("Failed to save XP file: %v", err)
	}

	// Reload the newly saved file.
	xpReloaded, err := LoadXPFile(savedFile)
	if err != nil {
		t.Fatalf("Failed to reload saved XP file: %v", err)
	}

	// Compare the structures.
	assertXPFileEqual(xpOriginal, xpReloaded, t)
}

func TestMarshal(t *testing.T) {
	path := filepath.Join(testDataDir, "simple_plain.xp")

	// Load raw, uncompressed REXPaint generated data.
	originalData, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read original plain XP file: %v", err)
	}

	// Load the same file to an XPFile struct.
	xp, err := LoadXPFile(path)
	if err != nil {
		t.Fatalf("Failed to load plain XP file: %v", err)
	}

	// Marshal the XPFile struct.
	marshaledData, err := Marshal(xp)
	if err != nil {
		t.Fatalf("Failed to marshal XP file: %v", err)
	}

	// Marshaled data MUST match original raw data byte for byte.
	if !bytes.Equal(originalData, marshaledData) {
		if len(originalData) != len(marshaledData) {
			t.Fatalf(
				"Marshaled data does not match original data (lengths: original=%d bytes, marshaled=%d bytes)",
				len(originalData), len(marshaledData),
			)
		}

		// Lengths are the same, so find where they differ.
		for i := range originalData {
			if originalData[i] != marshaledData[i] {
				t.Fatalf(
					"Marshaled data differs at byte %d: original=0x%02x, marshaled=0x%02x",
					i, originalData[i], marshaledData[i],
				)
			}
		}

		// Should be unreachable.
		t.Fatal("Marshaled data differs, but no exact difference found")
	}
}
