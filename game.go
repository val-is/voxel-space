package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 400
	screenHeight = 400
)

type Game struct {
	LoadedMap *Map
	Cam       *Camera
}

func (g *Game) Update(screen *ebiten.Image) error {
	fmt.Println(ebiten.CurrentFPS())
	// g.Cam.Y -= 1
	// g.Cam.X += 0.25
	vX := 0.0
	vY := 0.0
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
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		g.Cam.Height -= 1
	} else if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Cam.Height += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		g.Cam.Heading += 0.02
	} else if ebiten.IsKeyPressed(ebiten.KeyL) {
		g.Cam.Heading -= 0.02
	}
	h := g.Cam.Heading
	for h > math.Pi*2 {
		h -= math.Pi * 2
	}
	h -= math.Pi
	h *= -1
	g.Cam.X -= vX*math.Cos(h) - vY*math.Sin(h)
	g.Cam.Y -= vX*math.Sin(h) + vY*math.Cos(h)
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
