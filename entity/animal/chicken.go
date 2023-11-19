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
	Sprite       model.SpriteBody
	Collision    model.CollisionBody
	Destination  utils.Location
	AnimationTTL int
}

func NewChicken(spritesheet *ebiten.Image, tileX, tileY int) *Chicken {
	xLoc := tileX * utils.TileWidth
	yLoc := tileY * utils.TileHeight
	return &Chicken{
		Spritesheet: spritesheet,
		State:       utils.ChickenIdleState,
		Direction:   utils.AnimalRight,
		XLoc:        xLoc,
		YLoc:        yLoc,
		Sprite: model.SpriteBody{
			X:      xLoc,
			Y:      yLoc,
			Width:  utils.UnitSize,
			Height: utils.UnitSize,
		},
		Collision: model.CollisionBody{
			X:      xLoc + 8,
			Y:      yLoc + 16,
			Width:  16,
			Height: 14,
		},
	}
}

func (c *Chicken) UpdateFrame(currentFrame int) {
	if currentFrame%utils.AnimalFrameDelay == 0 {
		if c.AnimationTTL > 1 {
			c.AnimationTTL -= 1
		} else if c.AnimationTTL == 1 {
			c.AnimationTTL -= 1
			c.State = utils.ChickenIdleState
		}

		c.Frame += 1
		if c.Frame >= utils.AnimalFrameCount {
			c.Frame = 0
		}
	}
}
