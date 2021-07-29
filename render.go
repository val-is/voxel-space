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

func Render(c *Camera, m *Map, screen *ebiten.Image) error {
	if err := screen.Clear(); err != nil {
		return err
	}

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
		dz += 0.02
	}

	return nil
}
