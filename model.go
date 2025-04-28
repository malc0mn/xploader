package xploader

var (
	// DefaultForegroundColor is the initial foreground color of blank cells in REXPaint.
	DefaultForegroundColor = Color{R: 0, G: 0, B: 0}

	// InvisibleColor is the special color (absolute magenta) used internally by REXPaint. It is also the initial
	// background color of blank cells.
	// When assigned as background or foreground, the color will not be rendered visually.
	InvisibleColor = Color{R: 255, G: 0, B: 255}
)

// Color represents an RGB color.
type Color struct {
	R, G, B uint8
}

// IsInvisible will return true when the color is an absolute magenta. Absolute magenta is NEVER rendered: not as
// foreground, not as background.
func (c Color) IsInvisible() bool {
	return c == InvisibleColor
}

// Cell represents a single cell in a layer.
type Cell struct {
	Rune rune
	Fg   Color
	Bg   Color
}

// IsEmpty will return true when the artist did not paint this cell in REXPaint but left it untouched.
func (c Cell) IsEmpty() bool {
	return c.Rune == ' ' &&
		c.Fg == DefaultForegroundColor &&
		c.Bg.IsInvisible()
}

// NewEmptyCell returns a cell as REXPaint initialises it by default without the artist having touched it.
func NewEmptyCell() Cell {
	return Cell{
		Rune: ' ',
		Fg:   DefaultForegroundColor,
		Bg:   InvisibleColor,
	}
}

// Layer represents one layer of the XP file.
type Layer struct {
	ColumnMajor bool
	Width       uint32
	Height      uint32
	Cells       [][]Cell
}

// GetCell returns the cell at logical coordinates (x, y) based on the layer's memory layout.
// It automatically adjusts indexing if the layer was loaded in column-major order.
func (l *Layer) GetCell(x, y int) Cell {
	if l.ColumnMajor {
		return l.Cells[x][y]
	}
	return l.Cells[y][x]
}

// NewEmptyLayer returns a new layer of the given dimensions initialized with empty cells.
func NewEmptyLayer(width, height int) *Layer {
	return &Layer{
		Width:  uint32(width),
		Height: uint32(height),
		Cells: func() [][]Cell {
			tiles := make([][]Cell, height)
			for y := 0; y < height; y++ {
				row := make([]Cell, width)
				for x := 0; x < width; x++ {
					row[x] = NewEmptyCell()
				}
				tiles[y] = row
			}
			return tiles
		}(),
	}
}

// XPFile represents the parsed REXPaint .xp file.
type XPFile struct {
	Version int32
	Layers  []Layer
}

// AddLayer adds the given layer to the XPFile.
func (xp *XPFile) AddLayer(layer Layer) {
	xp.Layers = append(xp.Layers, layer)
}
