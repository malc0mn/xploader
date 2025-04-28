package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/malc0mn/xploder"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file.xp>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	path := os.Args[1]

	xp, err := xploader.LoadXPFile(path)
	if err != nil {
		log.Fatalf("Failed to load XP file: %v", err)
	}

	fmt.Printf("XP File: %s\n", path)
	fmt.Printf("Version: %d\n", xp.Version)
	fmt.Printf("Number of Layers: %d\n\n", len(xp.Layers))

	for idx, layer := range xp.Layers {
		fmt.Printf("Layer %d:\n", idx)
		fmt.Printf("  Dimensions: %dx%d\n", layer.Width, layer.Height)

		// Optional: count number of non-empty cells
		nonEmpty := 0
		for y := 0; y < int(layer.Height); y++ {
			for x := 0; x < int(layer.Width); x++ {
				cell := layer.GetCell(x, y)
				if !cell.IsEmpty() {
					nonEmpty++
				}
			}
		}
		fmt.Printf("  Non-empty cells: %d\n", nonEmpty)
		fmt.Println()
	}

	fmt.Println("Rendered layers")

	for layerIndex, layer := range xp.Layers {
		fmt.Printf("Layer %d (%dx%d):\n", layerIndex, layer.Width, layer.Height)

		height := int(layer.Height)
		width := int(layer.Width)

		fmt.Println("┌" + strings.Repeat("─", width) + "┐")
		for y := 0; y < height; y++ {
			fmt.Print("│")
			for x := 0; x < width; x++ {
				cell := layer.Cells[y][x]

				if cell.IsEmpty() {
					fmt.Print("\033[0m ")
					continue
				}

				// Set foreground and background colors.
				if !cell.Fg.IsInvisible() {
					fmt.Printf("\033[38;2;%d;%d;%dm", cell.Fg.R, cell.Fg.G, cell.Fg.B) // Foreground
				}
				if !cell.Bg.IsInvisible() {
					fmt.Printf("\033[48;2;%d;%d;%dm", cell.Bg.R, cell.Bg.G, cell.Bg.B) // Background
				}

				// Print rune.
				fmt.Printf("%c", cell.Rune)

				// Optionally: reset color after each rune (or after each line for optimization).
				fmt.Print("\033[0m")
			}
			fmt.Println("│")
		}
		fmt.Println("└" + strings.Repeat("─", width) + "┘")

		fmt.Println()
	}
}
