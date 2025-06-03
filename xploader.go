// Package xploader provides support for reading and writing REXPaint .xp files, including decoding and encoding CP437
// character data.
//
// Note: The CP437-to-Unicode mapping included here reflects the default font used by REXPaint. If you use a custom font
// (e.g., one that redefines glyphs 254/255 or others), the rendered output in REXPaint may differ from the Unicode
// characters provided by this mapping. The data remains technically valid, but the visual result depends on the
// selected font.
package xploader

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	gzipID1 byte = 0x1F
	gzipID2 byte = 0x8B
)

// LoadOptions controls how XP files are loaded.
type LoadOptions struct {
	// ColumnMajor preserves REXPaint's native column-major layout when true. Defaults to false, which returns data in
	// row-major order.
	ColumnMajor bool

	// RuneDecoder overrides how CP437 code points (0–255) are decoded to Unicode runes.
	// If nil, CP437Decoder is used. Useful when working with custom fonts.
	RuneDecoder func(int32) rune
}

// LoadXPFile loads a REXPaint .xp file from a filesystem path with default options and returns a pointer to an XPFile
// struct containing the fully parsed XP stream.
func LoadXPFile(path string) (*XPFile, error) {
	return LoadXPFileWithOptions(path, LoadOptions{ColumnMajor: false, RuneDecoder: CP437Decoder})
}

// LoadXPFileWithOptions loads a REXPaint .xp file with the specified options and returns a pointer to an XPFile struct
// containing the fully parsed XP stream.
func LoadXPFileWithOptions(path string, opts LoadOptions) (*XPFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	return LoadXPFromReader(f, opts)
}

// LoadXPFromReader loads a REXPaint .xp file from a reader with options and returns a pointer to an XPFile struct
// containing the fully parsed XP stream.
func LoadXPFromReader(r io.Reader, opts LoadOptions) (*XPFile, error) {
	isGzip, r, err := detectGzip(r)
	if err != nil {
		return nil, fmt.Errorf("failed to detect gzip: %w", err)
	}

	if isGzip {
		return LoadGzippedXPFromReader(r, opts)
	}
	return LoadPlainXPFromReader(r, opts)
}

// LoadGzippedXPFromReader will wrap the given io.Reader with a gzip.Reader and will return a pointer to an XPFile
// struct containing the fully parsed XP stream. It will fail if the source data is not gzipped.
func LoadGzippedXPFromReader(r io.Reader, opts LoadOptions) (*XPFile, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return LoadPlainXPFromReader(gr, opts)
}

// LoadPlainXPFromReader loads a RexPaint .xp file from an io.Reader.
func LoadPlainXPFromReader(r io.Reader, opts LoadOptions) (*XPFile, error) {
	var version int32
	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
		return nil, fmt.Errorf("failed to read version: %w", err)
	}

	var layerCount uint32
	if err := binary.Read(r, binary.LittleEndian, &layerCount); err != nil {
		return nil, fmt.Errorf("failed to read layer count: %w", err)
	}

	xp := &XPFile{
		Version: version,
		Layers:  make([]Layer, 0, layerCount),
	}

	for i := uint32(0); i < layerCount; i++ {
		layer, err := readLayer(r, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to read layer %d: %w", i, err)
		}
		xp.Layers = append(xp.Layers, *layer)
	}

	return xp, nil
}

// detectGzip peaks at the first two bytes to detect gzip compression. The caller MUST use the returned reader to read
// from since the original reader will have two bytes consumed after our peak.
func detectGzip(r io.Reader) (bool, io.Reader, error) {
	var header [2]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return false, nil, err
	}

	reader := io.MultiReader(bytes.NewReader(header[:]), r)

	if header[0] == gzipID1 && header[1] == gzipID2 {
		return true, reader, nil
	}
	return false, reader, nil
}

// readLayer reads a single layer from the XP file.
func readLayer(r io.Reader, opts LoadOptions) (*Layer, error) {
	var width, height uint32

	if err := binary.Read(r, binary.LittleEndian, &width); err != nil {
		return nil, fmt.Errorf("failed to read layer width: %w", err)
	}
	if err := binary.Read(r, binary.LittleEndian, &height); err != nil {
		return nil, fmt.Errorf("failed to read layer height: %w", err)
	}

	// Memory allocation based on desired layout.
	var outer, inner uint32
	outer = height
	inner = width
	if opts.ColumnMajor {
		outer = width
		inner = height
	}

	cells := make([][]Cell, outer)
	for i := uint32(0); i < outer; i++ {
		cells[i] = make([]Cell, inner)
	}

	for x := uint32(0); x < width; x++ {
		for y := uint32(0); y < height; y++ {
			var codepoint int32
			var fg Color
			var bg Color

			if err := binary.Read(r, binary.LittleEndian, &codepoint); err != nil {
				return nil, fmt.Errorf("failed to read codepoint at (%d,%d): %w", x, y, err)
			}

			if err := binary.Read(r, binary.LittleEndian, &fg); err != nil {
				return nil, fmt.Errorf("failed to read foreground color at (%d,%d): %w", x, y, err)
			}

			if err := binary.Read(r, binary.LittleEndian, &bg); err != nil {
				return nil, fmt.Errorf("failed to read background color at (%d,%d): %w", x, y, err)
			}

			ru := codepoint
			if opts.RuneDecoder != nil {
				ru = opts.RuneDecoder(codepoint)
			}

			if ru == '\x00' {
				ru = ' '
			}

			cell := Cell{
				Rune: ru,
				Fg:   fg,
				Bg:   bg,
			}

			if opts.ColumnMajor {
				cells[x][y] = cell
			} else {
				cells[y][x] = cell
			}
		}
	}

	return &Layer{
		ColumnMajor: opts.ColumnMajor,
		Width:       width,
		Height:      height,
		Cells:       cells,
	}, nil
}

// SaveOptions controls how XP files are saved.
type SaveOptions struct {
	// Gzip enables gzip compression of the output file. Defaults to true.
	Gzip bool

	// GzipLevel sets the compression level. Uses flate constants (e.g., flate.BestCompression).
	GzipLevel int

	// RuneEncoder overrides how Unicode runes are mapped to CP437 code points. If nil, runes are written as-is.
	// Useful for supporting custom fonts.
	RuneEncoder func(rune) int32
}

// SaveXPFile saves the XPFile to the given path, always compressed (recommended standard).
func SaveXPFile(xp *XPFile, path string) error {
	return SaveXPFileWithOptions(xp, path, SaveOptions{Gzip: true, GzipLevel: flate.BestCompression, RuneEncoder: CP437Encoder})
}

// SaveXPFileWithOptions saves the XPFile with full control over compression.
func SaveXPFileWithOptions(xp *XPFile, path string, opts SaveOptions) error {
	data, err := Marshal(xp, opts)
	if err != nil {
		return fmt.Errorf("failed to marshal XP file: %w", err)
	}

	if opts.Gzip {
		data, err = GzipData(data, opts.GzipLevel)
		if err != nil {
			return fmt.Errorf("failed to compress XP data: %w", err)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("failed to write XP data: %w", err)
	}

	return nil
}

// Marshal serializes the XPFile into uncompressed binary format, always column-major.
func Marshal(xp *XPFile, opts SaveOptions) ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.LittleEndian, xp.Version); err != nil {
		return nil, fmt.Errorf("failed to write version: %w", err)
	}

	if err := binary.Write(&buf, binary.LittleEndian, uint32(len(xp.Layers))); err != nil {
		return nil, fmt.Errorf("failed to write layer count: %w", err)
	}

	for i, layer := range xp.Layers {
		if err := marshalLayer(&buf, &layer, opts); err != nil {
			return nil, fmt.Errorf("failed to write layer %d: %w", i, err)
		}
	}

	return buf.Bytes(), nil
}

// GzipData compresses the given raw binary data using the gzip format.
// The level should be a valid compression level constant from the compress/flate package
// (e.g., flate.DefaultCompression, flate.BestCompression, flate.BestSpeed).
func GzipData(data []byte, level int) ([]byte, error) {
	var buf bytes.Buffer
	gw, err := gzip.NewWriterLevel(&buf, level)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip writer: %w", err)
	}

	if _, err := gw.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write gzip data: %w", err)
	}
	if err := gw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return buf.Bytes(), nil
}

// marshalLayer writes a single layer in column-major order.
func marshalLayer(w io.Writer, layer *Layer, opts SaveOptions) error {
	if err := binary.Write(w, binary.LittleEndian, layer.Width); err != nil {
		return fmt.Errorf("failed to write layer width: %w", err)
	}
	if err := binary.Write(w, binary.LittleEndian, layer.Height); err != nil {
		return fmt.Errorf("failed to write layer height: %w", err)
	}

	width := int(layer.Width)
	height := int(layer.Height)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			cell := layer.GetCell(x, y)
			r := cell.Rune
			if opts.RuneEncoder != nil {
				r = opts.RuneEncoder(r)
			}

			if err := binary.Write(w, binary.LittleEndian, r); err != nil {
				return fmt.Errorf("failed to write rune: %w", err)
			}
			if err := binary.Write(w, binary.LittleEndian, cell.Fg); err != nil {
				return fmt.Errorf("failed to write foreground color: %w", err)
			}
			if err := binary.Write(w, binary.LittleEndian, cell.Bg); err != nil {
				return fmt.Errorf("failed to write background color: %w", err)
			}
		}
	}

	return nil
}
