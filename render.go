package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
)

type Camera struct {
	Coords
	Height         float64
	Heading        float64
	FOV            float64
	RenderDistance float64
}

const (
	scaleHeight = 200
	horizon     = 100
)

func drawVertLine(screen *ebiten.Image, x, bottom, top float64, col color.Color) {
	for i := bottom; i <= top; i++ {
		screen.Set(int(x), int(i), col)
	}
}

func drawHorizLine(screen *ebiten.Image, y, left, right float64, col color.Color) {
	for i := left; i <= right; i++ {
		screen.Set(int(i), int(y), col)
	}
}

func renderTerrain(c *Camera, m *Map, screen *ebiten.Image) error {
	sinphi := math.Sin(c.Heading)
	cosphi := math.Cos(c.Heading)

	screenW, _ := screen.Size()

	dz := 1.0
	z := 1.0

	yBuffer := make([]int, screenW)
	for i := range yBuffer {
		yBuffer[i] = screenHeight
	}

	// render the terrain
	for z <= c.RenderDistance {
		left := Coords{c.X + (-cosphi*z - sinphi*z), c.Y + (sinphi*z - cosphi*z)}
		right := Coords{c.X + (cosphi*z - sinphi*z), c.Y + (-sinphi*z - cosphi*z)}

		dx := (right.X - left.X) / float64(screenW)
		dy := (right.Y - left.Y) / float64(screenW)

		for i := 0; i < screenW; i++ {
			mapHeight := m.HeightAt(left)
			height := (c.Height-mapHeight)/z*scaleHeight + horizon
			drawVertLine(screen, float64(i), height, float64(yBuffer[i]), m.ColorAt(left))
			if int(height) < yBuffer[i] {
				yBuffer[i] = int(height)
			}
			left.X += dx
			left.Y += dy
		}

		z += dz
		dz += 0.05
	}

	return nil
}

var (
	topSkyboxColor    = color.RGBA{10, 22, 50, 255}
	bottomSkyboxColor = color.RGBA{88, 84, 124, 255}
)

func renderSkybox(c *Camera, m *Map, screen *ebiten.Image) error {
	width, height := screen.Size()
	dRed := float64(int(bottomSkyboxColor.R)-int(topSkyboxColor.R)) / float64(height)
	dGreen := float64(int(bottomSkyboxColor.G)-int(topSkyboxColor.G)) / float64(height)
	dBlue := float64(int(bottomSkyboxColor.B)-int(topSkyboxColor.B)) / float64(height)
	col := topSkyboxColor

	for y := 0; y <= height; y++ {
		col = color.RGBA{
			uint8(float64(topSkyboxColor.R) + dRed*float64(y)),
			uint8(float64(topSkyboxColor.G) + dGreen*float64(y)),
			uint8(float64(topSkyboxColor.B) + dBlue*float64(y)),
			col.A,
		}
		drawHorizLine(screen, float64(y), 0, float64(width), col)
	}
	return nil
}

func Render(c *Camera, m *Map, screen *ebiten.Image) error {
	if err := screen.Clear(); err != nil {
		return err
	} else if err := renderSkybox(c, m, screen); err != nil {
		return err
	} else if err := renderTerrain(c, m, screen); err != nil {
		return err
	}

	return nil
}
