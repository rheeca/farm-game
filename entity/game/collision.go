package game

import (
	"github.com/co0p/tankism/lib/collision"
	"guion-2d-project3/entity/model"
)

//func hasMapCollisions(g *Game, animObj interfaces.AnimatedSprite) bool {
//	for tileY := 0; tileY < utils.MapRows; tileY += 1 {
//		for tileX := 0; tileX < utils.MapColumns; tileX += 1 {
//			tile := g.Environment.Maps[0].Layers[utils.ObjectsLayer].Tiles[tileY*utils.MapColumns+tileX]
//			if tile.ID == 0 {
//				continue
//			}
//			tileXpos := utils.TileWidth * tileX
//			tileYpos := utils.TileHeight * tileY
//
//			newX := animObj.GetXLoc() + animObj.GetDx()
//			newY := animObj.GetYLoc() + animObj.GetDy()
//			animBounds := collision.BoundingBox{
//				// bounding box for animated object made slightly smaller than the sprite
//				X:      float64(newX + animObj.GetWidth()/4),
//				Y:      float64(newY + animObj.GetHeight()/2),
//				Width:  float64(animObj.GetWidth() / 2),
//				Height: float64(animObj.GetHeight() / 2),
//			}
//			tileBounds := collision.BoundingBox{
//				X:      float64(tileXpos),
//				Y:      float64(tileYpos),
//				Width:  float64(utils.TileWidth),
//				Height: float64(utils.TileHeight),
//			}
//			if collision.AABBCollision(animBounds, tileBounds) {
//				return true
//			}
//		}
//	}
//	return false
//}

func hasCollision(dx, dy int, bodyA, bodyB model.CollisionBody) bool {
	// check if movement of bodyA collides with bodyB
	aBounds := collision.BoundingBox{
		X:      float64(bodyA.X0 + dx),
		Y:      float64(bodyA.Y0 + dy),
		Width:  float64(bodyA.X1 - bodyA.X0),
		Height: float64(bodyA.Y1 - bodyA.Y0),
	}
	bBounds := collision.BoundingBox{
		X:      float64(bodyB.X0),
		Y:      float64(bodyB.Y0),
		Width:  float64(bodyB.X1 - bodyB.X0),
		Height: float64(bodyB.Y1 - bodyB.Y0),
	}
	if collision.AABBCollision(aBounds, bBounds) {
		return true
	}
	return false
}

func playerHasCollisions(g *Game) bool {
	// TODO: check for map collisions
	//if hasMapCollisions(g, g.Player) {
	//	return true
	//}

	// check for animated entities collisions
	for _, c := range g.Chickens {
		if hasCollision(g.Player.Dx, g.Player.Dy, g.Player.Collision, c.Collision) {
			return true
		}
	}
	return false
}

func updateAnimals(g *Game) {
	for i, _ := range g.Chickens {
		g.Chickens[i].UpdateFrame(g.CurrentFrame)
	}
}
