package animal

import (
	"github.com/hajimehoshi/ebiten/v2"
	"guion-2d-project3/entity/model"
	"guion-2d-project3/utils"
	"math/rand"
)

type Chicken struct {
	Spritesheet *ebiten.Image
	Frame       int
	State       int
	Direction   int
	XLoc        int
	YLoc        int
	Dx          int
	Dy          int
	Sprite      model.SpriteBody
	Collision   model.CollisionBody
	Destination utils.Location
	StateTTL    int
}

func NewChicken(spritesheet *ebiten.Image, loc utils.Location) *Chicken {
	xLoc := loc.X * utils.TileWidth
	yLoc := loc.Y * utils.TileHeight
	chicken := Chicken{
		Spritesheet: spritesheet,
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
	chicken.RandomMovement(0)
	return &chicken
}

func (c *Chicken) UpdateFrame(currentFrame int) {
	if currentFrame%utils.AnimalFrameDelay == 0 {
		if c.StateTTL > 0 {
			c.StateTTL -= 1
		}

		c.Frame += 1
		if c.Frame >= utils.AnimalFrameCount {
			c.Frame = 0
		}
	}
}

func (c *Chicken) RandomMovement(currentFrame int) {
	if !(currentFrame%utils.AnimalFrameDelay == 0) {
		return
	}

	if c.StateTTL > 0 {
		// don't set random state if a current state is not yet finished
		return
	}
	c.State = rand.Intn(2) + 1
	c.Direction = rand.Intn(2)
	c.StateTTL = (rand.Intn(3) + 1) * utils.AnimalFrameCount

	if c.State == utils.ChickenWalkState {
		x := c.XLoc + (rand.Intn(4)-2)*utils.TileWidth
		y := c.YLoc + (rand.Intn(4)-2)*utils.TileHeight
		if x < 0 {
			x = 0
		}
		if y < 0 {
			y = 0
		}
		c.Destination = utils.Location{
			X: x % utils.MapWidth,
			Y: y % utils.MapHeight,
		}
	}
}

func (c *Chicken) UpdateLocation() {
	c.XLoc += c.Dx
	c.YLoc += c.Dy
	c.Sprite.X += c.Dx
	c.Sprite.Y += c.Dy
	c.Collision.X += c.Dx
	c.Collision.Y += c.Dy
	c.Dx = 0
	c.Dy = 0
}
