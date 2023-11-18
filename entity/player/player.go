package player

import (
	"github.com/co0p/tankism/lib/collision"
	"github.com/hajimehoshi/ebiten/v2"
	"guion-2d-project3/interfaces"
	"guion-2d-project3/utils"
)

type Player struct {
	Spritesheet  *ebiten.Image
	Frame        int
	State        int
	StateTTL     int
	Direction    int
	XLoc         int
	YLoc         int
	Dy           int
	Dx           int
	SpriteWidth  int
	SpriteHeight int
	Collision    CollisionBody
	Backpack     [utils.BackpackSize]int
	EquippedItem int
}

type CollisionBody struct {
	X0 int
	Y0 int
	X1 int
	Y1 int
}

func NewPlayer(spritesheet *ebiten.Image) *Player {
	return &Player{
		Spritesheet:  spritesheet,
		XLoc:         utils.StartingX * utils.TileWidth,
		YLoc:         utils.StartingY * utils.TileWidth,
		SpriteWidth:  utils.PlayerSpriteWidth,
		SpriteHeight: utils.PlayerSpriteHeight,
		Collision: CollisionBody{
			X0: 39,
			Y0: 50,
			X1: 57,
			Y1: 64,
		},
		Backpack:     [utils.BackpackSize]int{2, 3, 10, 19, 42},
		EquippedItem: 0,
	}
}

func (p *Player) HasCollisionWith(object interfaces.AnimatedSprite) bool {
	playerBounds := collision.BoundingBox{
		X:      float64(p.XLoc + p.Dx),
		Y:      float64(p.YLoc + p.Dy),
		Width:  float64(p.SpriteWidth),
		Height: float64(p.SpriteHeight),
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
	return p.SpriteWidth
}

func (p *Player) GetHeight() int {
	return p.SpriteHeight
}

func (p *Player) UpdateLocation() {
	p.XLoc += p.Dx
	p.YLoc += p.Dy
	p.Dx = 0
	p.Dy = 0
}

func (p *Player) UpdateFrame(currentFrame int) {
	if currentFrame%utils.PlayerFrameDelay == 0 {
		if p.StateTTL > 1 {
			p.StateTTL -= 1
		} else if p.StateTTL == 1 {
			p.StateTTL -= 1
			p.State = utils.IdleState
		}

		p.Frame += 1
		if p.Frame >= utils.PlayerFrameCount {
			p.Frame = 0
		}
	}
}
