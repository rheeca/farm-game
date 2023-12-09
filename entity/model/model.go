package model

import (
	"guion-2d-project3/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type CollisionBody struct {
	X      int
	Y      int
	Width  int
	Height int
}

type SpriteBody struct {
	X      int
	Y      int
	Width  int
	Height int
}

type UIState struct {
	SelectedCharacter int
	SelectedRecipe    int
	ErrorMessage      string
	ErrorMessageTTL   int
	ImageToShow       *ebiten.Image
	ImageTTL          int
}

type BackpackItem struct {
	ID    int
	Count int
}

type delayfcn func()

type Object struct {
	// Object type
	Type int

	// x and y location of object on the map
	XLoc int
	YLoc int

	Sprite      SpriteBody
	Collision   CollisionBody
	IsCollision bool
	IsNil       bool
	Frame       int

	// When object reaches zero health, object will be removed
	Health int

	// If IsAnimating is true, the CurrentAnimation will play
	// on the object
	IsAnimating      bool
	CurrentAnimation int

	// AnimationDelay is the number of frames to wait for before
	// the CurrentAnimation is played
	AnimationDelay int

	// AnimationTTL is how many frames the animation will play for
	AnimationTTL int

	// If DoDelayFcn is true, DelayFcn will be executed after
	// the end of the animation
	DoDelayFcn bool
	DelayFcn   delayfcn `json:"-"`
}

func (o *Object) StartAnimation(animation, animationTTL, animationDelay int, doDelayFcn bool, fcn delayfcn) {
	o.Frame = 0
	o.IsAnimating = true
	o.CurrentAnimation = animation
	o.AnimationTTL = animationTTL
	o.AnimationDelay = animationDelay
	o.DoDelayFcn = doDelayFcn
	o.DelayFcn = fcn
}

func (o *Object) UpdateFrame(currentFrame int) {
	if !o.IsAnimating {
		return
	}
	if currentFrame%utils.FrameDelay == 0 {
		if o.AnimationDelay > 0 {
			o.AnimationDelay -= 1
			return
		}

		o.Frame += 1
		if o.AnimationTTL > 1 {
			o.AnimationTTL -= 1
		} else if o.AnimationTTL == 1 {
			o.IsAnimating = false
			o.AnimationTTL = 0
			o.AnimationDelay = 0
			o.Frame = 0
			if o.DoDelayFcn {
				o.DelayFcn()
				o.DoDelayFcn = false
			}
		}
	}
}
