package player

import (
	"github.com/co0p/tankism/lib/collision"
	"github.com/hajimehoshi/ebiten/v2"
	"guion-2d-project3/interfaces"
	"guion-2d-project3/utils"
)

type Player struct {
	Spritesheet *ebiten.Image
	Frame       int
	Direction   int
	XLoc        int
	YLoc        int
	Dy          int
	Dx          int
	Width       int
	Height      int
}

func NewPlayer(spritesheet *ebiten.Image) *Player {
	return &Player{
		Spritesheet: spritesheet,
		XLoc:        utils.StartingX * utils.TileWidth,
		YLoc:        utils.StartingY * utils.TileWidth,
		Width:       spritesheet.Bounds().Dx() / utils.AnimFrameCount,
		Height:      spritesheet.Bounds().Dy() / utils.AnimFrameCount,
	}
}

func (p *Player) HasCollisionWith(object interfaces.AnimatedSprite) bool {
	playerBounds := collision.BoundingBox{
		X:      float64(p.XLoc + p.Dx),
		Y:      float64(p.YLoc + p.Dy),
		Width:  float64(p.Width),
		Height: float64(p.Height),
	}
	objectBounds := collision.BoundingBox{
		X:      float64(object.GetXLoc() + object.GetDx()),
		Y:      float64(object.GetYLoc() + object.GetDy()),
		Width:  float64(object.GetWidth()),
		Height: float64(object.GetHeight()),
	}
	if collision.AABBCollision(playerBounds, objectBounds) {
		return true
	}
	return false
}

func (p *Player) GetXLoc() int {
	return p.XLoc
}

func (p *Player) GetYLoc() int {
	return p.YLoc
}

func (p *Player) GetDx() int {
	return p.Dx
}

func (p *Player) GetDy() int {
	return p.Dy
}

func (p *Player) GetWidth() int {
	return p.Width
}

func (p *Player) GetHeight() int {
	return p.Height
}

func (p *Player) UpdateLocation() {
	p.XLoc += p.Dx
	p.YLoc += p.Dy
	p.Dx = 0
	p.Dy = 0
}

func (p *Player) UpdateFrame(currentFrame int) {
	if currentFrame%utils.FrameDelay == 0 {
		p.Frame += 1
		if p.Frame >= utils.AnimFrameCount {
			p.Frame = 0
		}
	}
}
