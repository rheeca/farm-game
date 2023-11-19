package animal

import (
	"github.com/hajimehoshi/ebiten/v2"
	"guion-2d-project3/entity/model"
	"guion-2d-project3/utils"
)

type Cow struct {
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

func NewCow(spritesheet *ebiten.Image, tileX, tileY int) *Cow {
	xLoc := tileX * utils.TileWidth
	yLoc := tileY * utils.TileHeight
	return &Cow{
		Spritesheet: spritesheet,
		State:       utils.CowIdleState,
		Direction:   utils.AnimalRight,
		XLoc:        xLoc,
		YLoc:        yLoc,
		Sprite: model.SpriteBody{
			X:      xLoc + 10,
			Y:      yLoc + 27,
			Width:  50,
			Height: 30,
		},
		Collision: model.CollisionBody{
			X0: xLoc + 12,
			Y0: yLoc + 30,
			X1: xLoc + 54,
			Y1: yLoc + 57,
		},
	}
}

func (c *Cow) UpdateFrame(currentFrame int) {
	if currentFrame%utils.AnimalFrameDelay == 0 {
		if c.AnimationTTL > 1 {
			c.AnimationTTL -= 1
		} else if c.AnimationTTL == 1 {
			c.AnimationTTL -= 1
			c.State = utils.CowIdleState
		}

		c.Frame += 1
		if c.Frame >= utils.AnimalFrameCount {
			c.Frame = 0
		}
	}
}
