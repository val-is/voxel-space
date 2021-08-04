package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
)

type Camera struct {
	Coords
	Height         float64
	Heading        float64
	Pitch          float64
	FOV            float64
	RenderDistance float64
	buffer         *ebiten.Image
}

const (
	scaleHeight = 200
	horizon     = 0
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

const (
	scatterDrawDist   = 25
	scatterFOVPadding = 30 / 180 * math.Pi
)

func renderTerrain(c *Camera, m *Map, screen *ebiten.Image) error {
	sinphi := math.Sin(c.Heading)
	cosphi := math.Cos(c.Heading)

	screenW, screenH := screen.Size()

	hz := horizon + float64(screenH)*math.Sin(c.Pitch)

	dz := 0.01
	z := 1.0

	yBuffer := make([]int, screenW)
	for i := range yBuffer {
		yBuffer[i] = screenHeight
	}

	// get scatters within the FOV that should be rendered
	// scattersRendering := make([]*Scatter, 0)
	leftFOV := c.Heading - c.FOV/2 - scatterFOVPadding
	if leftFOV < -math.Pi {
		leftFOV += 2 * math.Pi
	}
	rightFOV := c.Heading + c.FOV/2 + scatterFOVPadding
	if rightFOV < -math.Pi {
		rightFOV += 2 * math.Pi
	}
	for _, scatter := range m.Scatters {
		dist := math.Hypot(c.X-scatter.X, c.Y-scatter.Y)
		if dist >= scatterDrawDist {
			continue
		}
		angle := math.Atan2(-(scatter.X - c.X), -(scatter.Y - c.Y))
		if (leftFOV < rightFOV) && !(leftFOV <= angle && angle <= rightFOV) ||
			(leftFOV > rightFOV) && !(leftFOV <= angle && angle <= math.Pi || angle <= rightFOV && -math.Pi <= angle) {
			continue
		}
		fmt.Printf("rendering %f <= %f <= %f\n", leftFOV, angle, rightFOV)
	}

	// render the terrain
	for z <= c.RenderDistance {
		left := Coords{c.X + (-cosphi*z - sinphi*z), c.Y + (sinphi*z - cosphi*z)}
		right := Coords{c.X + (cosphi*z - sinphi*z), c.Y + (-sinphi*z - cosphi*z)}

		dx := (right.X - left.X) / float64(screenW)
		dy := (right.Y - left.Y) / float64(screenW)

		for i := 0; i < screenW; i++ {
			mapHeight := m.HeightAt(left)
			height := (c.Height-mapHeight)/z*scaleHeight + hz
			drawVertLine(screen, float64(i), height, float64(yBuffer[i]), m.ColorAt(left))
			drawVertLine(c.buffer, float64(i), height, float64(yBuffer[i]), color.NRGBA{uint8(z), 0, 0, 255})
			if int(height) < yBuffer[i] {
				yBuffer[i] = int(height)
			}
			left.X += dx
			left.Y += dy
		}

		z += dz
		dz += 0.005
	}

	return nil
}

func renderScatters(c *Camera, m *Map, screen *ebiten.Image) error {
	for i := range m.Scatters {
		if err := drawScatter(c, m, screen, m.Scatters[i]); err != nil {
			return err
		}
	}
	return nil
}

func drawScatter(c *Camera, m *Map, screen *ebiten.Image, scatter *Scatter) error {
	scatterBuffer, err := ebiten.NewImageFromImage(c.buffer, ebiten.FilterDefault)
	if err != nil {
		return err
	}
	dist := math.Hypot(c.X-scatter.X, c.Y-scatter.Y)
	if dist > scatterDrawDist {
		return nil
	}
	w, h := scatterBuffer.Size()
	z, _, _, _ := scatterBuffer.At(200, 200).RGBA()
	z /= 256
	fmt.Println(z, dist)
	for x := 0; x <= w; x++ {
		for y := 0; y <= h; y++ {
			col := scatterBuffer.At(x, y)
			buffDist, _, _, _ := col.RGBA()
			buffDist /= 256
			if buffDist <= uint32(dist) || dist == 0 {
				scatterBuffer.Set(x, y, color.NRGBA{0, 0, 0, 0})
			} else {
				scatterBuffer.Set(x, y, color.NRGBA{255, 255, 255, 255})
			}
		}
	}

	screenW, screenH := screen.Size()

	sinphi := math.Sin(c.Heading)
	cosphi := math.Cos(c.Heading)

	leftX := c.X + (-cosphi*dist - sinphi*dist)
	rightX := c.X + (cosphi*dist - sinphi*dist)

	dx := (rightX - leftX) / float64(screenW)

	hz := horizon + float64(screenH)*math.Sin(c.Pitch)

	mapHeight := m.HeightAt(scatter.Coords)
	height := (c.Height-mapHeight-3)/dist*scaleHeight + hz
	op := ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Scale(scatterDrawDist/dist*scatter.Scale, scatterDrawDist/dist*scatter.Scale)
	op.GeoM.Translate((scatter.X-leftX)/dx, height)
	op.CompositeMode = ebiten.CompositeModeMultiply
	scatterBuffer.DrawImage(scatter.Sprite, &op)
	screen.DrawImage(scatterBuffer, nil)

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
	if c.buffer == nil {
		w, h := screen.Size()
		buf, err := ebiten.NewImage(w, h, ebiten.FilterDefault)
		if err != nil {
			return err
		}
		c.buffer = buf
	}
	if err := c.buffer.Clear(); err != nil {
		return err
	} else if err := screen.Clear(); err != nil {
		return err
	} else if err := renderSkybox(c, m, screen); err != nil {
		return err
	} else if err := renderTerrain(c, m, screen); err != nil {
		return err
	} else if err := renderScatters(c, m, screen); err != nil {
		return err
	}

	return nil
}
