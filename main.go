package main

import "github.com/hajimehoshi/ebiten"

func main() {
	m, err := LoadMap("maps/ice")
	if err != nil {
		panic(err)
	}

	cam := &Camera{
		Coords:         Coords{X: 0, Y: 0},
		Height:         30,
		Heading:        0,
		FOV:            90,
		RenderDistance: 1000,
	}

	g := Game{
		LoadedMap: m,
		Cam:       cam,
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Test Voxel Engine")
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}
}
