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
	Sprite       model.SpriteBody
	Collision    model.CollisionBody
	Backpack     [utils.BackpackSize]BackpackItem
	EquippedItem int
}

type BackpackItem struct {
	ID    int
	Count int
}

func NewPlayer(spritesheet *ebiten.Image, startingX, startingY int) *Player {
	return &Player{
		Spritesheet: spritesheet,
		XLoc:        startingX - 39,
		YLoc:        startingY - 35,
		Sprite: model.SpriteBody{
			X:      startingX,
			Y:      startingY,
			Width:  20,
			Height: 30,
		},
		Collision: model.CollisionBody{
			X:      startingX,
			Y:      startingY + 15,
			Width:  18,
			Height: 16,
		},
		Backpack: [utils.BackpackSize]BackpackItem{
			{ID: 2, Count: 1},
			{ID: 3, Count: 1},
			{ID: 10, Count: 1},
		},
		EquippedItem: 0,
	}
}

func (p *Player) AddToBackpack(itemID, count int) bool {
	// if item already exists in backpack and is not a tool, add count
	for i, v := range p.Backpack {
		if v.ID == itemID && !isTool(v.ID) {
			p.Backpack[i].Count += count
			return true
		}
	}
	// if there is an empty slot, put new item there
	for i, v := range p.Backpack {
		if v.ID == 0 {
			p.Backpack[i] = BackpackItem{
				ID:    itemID,
				Count: count,
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
	return p.Sprite.Width
}

func (p *Player) GetHeight() int {
	return p.Sprite.Height
}

func (p *Player) UpdateLocation() {
	p.XLoc += p.Dx
	p.YLoc += p.Dy
	p.Sprite.X += p.Dx
	p.Sprite.Y += p.Dy
	p.Collision.X += p.Dx
	p.Collision.Y += p.Dy
	p.Dx = 0
	p.Dy = 0
}

func (p *Player) ChangeLocation(x, y int) {
	p.XLoc = x - 39
	p.YLoc = y - 35
	p.Sprite.X = x
	p.Sprite.Y = y
	p.Collision.X = x
	p.Collision.Y = y + 15
	p.Dx = 0
	p.Dy = 0
}

func (p *Player) CalcTargetBox() model.CollisionBody {
	var xLoc, yLoc int
	if p.Direction == utils.Front {
		xLoc = p.XLoc + utils.UnitSize
		yLoc = p.YLoc + utils.UnitSize + utils.UnitSize/2
	} else if p.Direction == utils.Back {
		xLoc = p.XLoc + utils.UnitSize
		yLoc = p.YLoc + utils.UnitSize/2
	} else if p.Direction == utils.Left {
		xLoc = p.XLoc + utils.UnitSize/2
		yLoc = p.YLoc + utils.UnitSize
	} else if p.Direction == utils.Right {
		xLoc = p.XLoc + utils.UnitSize + utils.UnitSize/2
		yLoc = p.YLoc + utils.UnitSize
	}
	return model.CollisionBody{
		X:      xLoc,
		Y:      yLoc,
		Width:  utils.UnitSize,
		Height: utils.UnitSize,
	}
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
