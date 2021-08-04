package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Scatter struct {
	Coords
	Sprite *ebiten.Image
	Scale  float64
}

var (
	scatterCache = make(map[string]*ebiten.Image)
)

const (
	scatterPath = "sprites"
)

func loadScatterSprite(filename string) (*ebiten.Image, error) {
	if sprite, present := scatterCache[filename]; present {
		return sprite, nil
	} else {
		sprite, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("%s/%s", scatterPath, filename), ebiten.FilterDefault)
		if err != nil {
			return nil, err
		}
		scatterCache[filename] = sprite
		return sprite, nil
	}
}

func NewScatter(filename string, x, y, scale float64) (*Scatter, error) {
	if sprite, err := loadScatterSprite(filename); err != nil {
		return nil, err
	} else {
		return &Scatter{
			Coords: Coords{X: x, Y: y},
			Sprite: sprite,
			Scale:  scale,
		}, nil
	}
}
