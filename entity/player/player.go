package player

import (
	"github.com/hajimehoshi/ebiten/v2"
	"guion-2d-project3/entity/model"
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
	Collision    model.CollisionBody
	Backpack     [utils.BackpackSize]BackpackItem
	EquippedItem int
}

type BackpackItem struct {
	ID    int
	Count int
}

func NewPlayer(spritesheet *ebiten.Image) *Player {
	xLoc := utils.StartingX * utils.TileWidth
	yLoc := utils.StartingY * utils.TileWidth
	return &Player{
		Spritesheet:  spritesheet,
		XLoc:         xLoc,
		YLoc:         yLoc,
		SpriteWidth:  utils.PlayerSpriteWidth,
		SpriteHeight: utils.PlayerSpriteHeight,
		Collision: model.CollisionBody{
			X0: xLoc + 39,
			Y0: yLoc + 50,
			X1: xLoc + 57,
			Y1: yLoc + 64,
		},
		Backpack: [utils.BackpackSize]BackpackItem{
			{ID: 2, Count: 1},
			{ID: 3, Count: 1},
			{ID: 10, Count: 1},
		},
		EquippedItem: 0,
	}
}

func (p *Player) AddToBackpack(itemID int) bool {
	// if item already exists in backpack and is not a tool, add count by 1
	for i, v := range p.Backpack {
		if v.ID == itemID && !isTool(v.ID) {
			p.Backpack[i].Count += 1
			return true
		}
	}
	// if there is an empty slot, put new item there
	for i, v := range p.Backpack {
		if v.ID == 0 {
			p.Backpack[i] = BackpackItem{
				ID:    itemID,
				Count: 1,
			}
			return true
		}
	}
	return false
}

func isTool(itemID int) bool {
	if itemID == 2 || itemID == 3 || itemID == 10 {
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
	p.Collision.X0 += p.Dx
	p.Collision.X1 += p.Dx
	p.Collision.Y0 += p.Dy
	p.Collision.Y1 += p.Dy
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
