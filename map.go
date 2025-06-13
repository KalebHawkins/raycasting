package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Map struct {
	Grid     []int
	Rows     int
	Columns  int
	TileSize int
}

func NewMap(rows, columns, tileSize int) *Map {
	m := &Map{
		Grid: []int{
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1,
			1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1,
			1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		},
		Rows:     rows,
		Columns:  columns,
		TileSize: tileSize,
	}

	return m
}

func (m *Map) Draw(dst *ebiten.Image) {
	for index := range m.Grid {
		row := index / m.Columns
		col := index % m.Columns

		tileX := col * m.TileSize
		tileY := row * m.TileSize

		clr := color.Black
		if m.Grid[index] == 0 {
			clr = color.White
		}

		vector.DrawFilledRect(dst, MiniMapScaleFactor*float32(tileX), MiniMapScaleFactor*float32(tileY), MiniMapScaleFactor*float32(m.TileSize), MiniMapScaleFactor*float32(m.TileSize), clr, true)
		vector.StrokeRect(dst, MiniMapScaleFactor*float32(tileX), MiniMapScaleFactor*float32(tileY), MiniMapScaleFactor*float32(m.TileSize), MiniMapScaleFactor*float32(m.TileSize), 0.5, color.Black, true)
	}
}

func (m *Map) AtPixel(x, y int) (int, error) {
	if x < 0 || y < 0 || x >= m.Columns*m.TileSize || y >= m.Rows*m.TileSize {
		return 0, fmt.Errorf("index out of range for map position (%d, %d)", x, y)
	}

	gridX := x / m.TileSize
	gridY := y / m.TileSize
	idx := gridY*m.Columns + gridX
	return m.Grid[idx], nil
}
