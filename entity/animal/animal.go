package animal

import (
	"github.com/co0p/tankism/lib/collision"
	"github.com/hajimehoshi/ebiten/v2"
	"guion-2d-project3/interfaces"
	"guion-2d-project3/utils"
)

type Animal struct {
	Spritesheet *ebiten.Image
	Frame       int
	Direction   int
	XLoc        int
	YLoc        int
	Dx          int
	Dy          int
	Width       int
	Height      int
	Path        []utils.Location
	Destination int
}

func NewAnimal(spritesheet *ebiten.Image, path []utils.Location) *Animal {
	return &Animal{
		Spritesheet: spritesheet,
		XLoc:        path[0].X * utils.TileWidth,
		YLoc:        path[0].Y * utils.TileHeight,
		Width:       spritesheet.Bounds().Dx() / utils.AnimFrameCount,
		Height:      spritesheet.Bounds().Dy() / utils.AnimFrameCount,
		Path:        path,
	}
}

func (a *Animal) HasCollisionWith(object interfaces.AnimatedSprite) bool {
	animalBounds := collision.BoundingBox{
		X:      float64(a.XLoc + a.Dx),
		Y:      float64(a.YLoc + a.Dy),
		Width:  float64(a.Width),
		Height: float64(a.Height),
	}

	objectBounds := collision.BoundingBox{
		X:      float64(object.GetXLoc() + object.GetDx()),
		Y:      float64(object.GetYLoc() + object.GetDy()),
		Width:  float64(object.GetWidth()),
		Height: float64(object.GetHeight()),
	}
	if collision.AABBCollision(animalBounds, objectBounds) {
		return true
	}
	return false
}

func (a *Animal) GetXLoc() int {
	return a.XLoc
}

func (a *Animal) GetYLoc() int {
	return a.YLoc
}

func (a *Animal) GetDx() int {
	return a.Dx
}

func (a *Animal) GetDy() int {
	return a.Dy
}

func (a *Animal) GetWidth() int {
	return a.Width
}

func (a *Animal) GetHeight() int {
	return a.Height
}

func (a *Animal) UpdateLocation() {
	a.XLoc += a.Dx
	a.YLoc += a.Dy
	a.Dx = 0
	a.Dy = 0
}

func (a *Animal) UpdateFrame(currentFrame int) {
	if currentFrame%utils.AnimalFrameDelay == 0 {
		a.Frame += 1
		if a.Frame >= utils.AnimFrameCount {
			a.Frame = 0
		}
	}
}
