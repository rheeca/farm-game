package model

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

type Object struct {
	Type      int
	XLoc      int
	YLoc      int
	Sprite    SpriteBody
	Collision CollisionBody
	Frame     int
	Health    int
}
