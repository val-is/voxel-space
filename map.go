package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Coords struct {
	X, Y float64
}

type Map struct {
	Colormap      *ebiten.Image
	Heightmap     *ebiten.Image
	Width, Height float64
	NormalizedAlt float64
	Scatters      []*Scatter
}

func (m Map) getMapCoords(coords Coords) (x, y int) {
	x = int(math.Round(coords.X)) % int(m.Width)
	if x < 0 {
		x += int(m.Width)
	}
	y = int(math.Round(coords.Y)) % int(m.Height)
	if y < 0 {
		y += int(m.Height)
	}
	return x, y
}

func (m *Map) HeightAt(coords Coords) float64 {
	col := m.Heightmap.At(m.getMapCoords(coords))
	value, _, _, alpha := col.RGBA()
	return float64(value) / float64(alpha) * m.NormalizedAlt
}

func (m *Map) ColorAt(coords Coords) color.Color {
	return m.Colormap.At(m.getMapCoords(coords))
}

const (
	heightmapFilename = "depth.png"
	colormapFilename  = "color.png"

	mapAlt = 100
)

func LoadMap(dir string) (*Map, error) {
	heightmap, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("%s/%s", dir, heightmapFilename), ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}

	colormap, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("%s/%s", dir, colormapFilename), ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}

	width, height := colormap.Size()

	return &Map{
		Colormap:      colormap,
		Heightmap:     heightmap,
		Width:         float64(width),
		Height:        float64(height),
		NormalizedAlt: mapAlt,
		Scatters:      make([]*Scatter, 0),
	}, nil
}
