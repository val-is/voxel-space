package main

import (
	"math"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 400
	screenHeight = 400
)

var lastCursorX, lastCursorY int

type Game struct {
	LoadedMap *Map
	Cam       *Camera
}

func (g *Game) Update(screen *ebiten.Image) error {
	vX := 0.0
	vY := 0.0
	velScaling := 0.1
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		vY = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		vY = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		vX = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		vX = 1
	}
	// if ebiten.IsKeyPressed(ebiten.KeyShift) {
	// 	g.Cam.Height -= 1
	// } else if ebiten.IsKeyPressed(ebiten.KeySpace) {
	// 	g.Cam.Height += 1
	// }
	g.Cam.Height = g.LoadedMap.HeightAt(g.Cam.Coords) + 1

	cursorX, cursorY := ebiten.CursorPosition()
	cDx := cursorX - lastCursorX
	cDy := cursorY - lastCursorY
	lastCursorX = cursorX
	lastCursorY = cursorY

	cDx = int(math.Max(math.Min(200, float64(cDx)), -200))
	cDy = int(math.Max(math.Min(200, float64(cDy)), -200))

	g.Cam.Heading += -0.002 * float64(cDx)
	g.Cam.Pitch += -0.002 * float64(cDy)

	if g.Cam.Heading > math.Pi {
		g.Cam.Heading = -math.Pi + (g.Cam.Heading - math.Pi)
	} else if g.Cam.Heading < -math.Pi {
		g.Cam.Heading = math.Pi - (-g.Cam.Heading - math.Pi)
	}

	h := g.Cam.Heading
	for h > math.Pi*2 {
		h -= math.Pi * 2
	}
	h -= math.Pi
	h *= -1
	vX *= velScaling
	vY *= velScaling
	g.Cam.X -= vX*math.Cos(h) - vY*math.Sin(h)
	g.Cam.Y -= vX*math.Sin(h) + vY*math.Cos(h)
	g.Cam.Pitch = math.Min(math.Max(g.Cam.Pitch, -math.Pi/2), math.Pi/2)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if err := Render(g.Cam, g.LoadedMap, screen); err != nil {
		panic(err)
	}
}

func (g *Game) Layout(width, height int) (int, int) {
	return screenWidth, screenHeight
}
