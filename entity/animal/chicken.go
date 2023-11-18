package animal

import (
	"github.com/hajimehoshi/ebiten/v2"
	"guion-2d-project3/entity/model"
	"guion-2d-project3/utils"
)

type Chicken struct {
	Spritesheet  *ebiten.Image
	Frame        int
	State        int
	Direction    int
	XLoc         int
	YLoc         int
	Dx           int
	Dy           int
	SpriteWidth  int
	SpriteHeight int
	Collision    model.CollisionBody
	Destination  utils.Location
	AnimationTTL int
}

func NewChicken(spritesheet *ebiten.Image) *Chicken {
	xLoc := 5 * utils.TileWidth
	yLoc := 5 * utils.TileHeight
	return &Chicken{
		Spritesheet:  spritesheet,
		State:        utils.ChickenIdleState,
		Direction:    utils.AnimalRight,
		XLoc:         xLoc,
		YLoc:         yLoc,
		SpriteWidth:  utils.UnitSize,
		SpriteHeight: utils.UnitSize,
		Collision: model.CollisionBody{
			X0: xLoc + 8,
			Y0: yLoc + 16,
			X1: xLoc + 24,
			Y1: yLoc + 30,
		},
	}
}

func (c *Chicken) UpdateFrame(currentFrame int) {
	if currentFrame%utils.AnimalFrameDelay == 0 {
		c.Frame += 1
		if c.Frame >= utils.AnimalFrameCount {
			c.Frame = 0
		}
	}
}
