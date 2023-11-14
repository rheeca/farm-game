package interfaces

type AnimatedSprite interface {
	HasCollisionWith(object AnimatedSprite) bool
	GetXLoc() int
	GetYLoc() int
	GetDx() int
	GetDy() int
	GetWidth() int
	GetHeight() int
}
