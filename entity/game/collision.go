package game

import (
	"github.com/co0p/tankism/lib/collision"
	"guion-2d-project3/entity/model"
	"guion-2d-project3/utils"
)

func hasMapCollisions(g *Game, dx, dy int, collisionBody model.CollisionBody) bool {
	for tileY := 0; tileY < utils.MapRows; tileY += 1 {
		for tileX := 0; tileX < utils.MapColumns; tileX += 1 {
			for _, layer := range utils.CollisionLayers {
				tile := g.Environment.Maps[g.CurrentMap].Layers[layer].Tiles[tileY*utils.MapColumns+tileX]
				if tile.ID == 0 {
					continue
				}
				tileXpos := utils.TileWidth * tileX
				tileYpos := utils.TileHeight * tileY

				tileCollision := model.CollisionBody{
					X:      tileXpos,
					Y:      tileYpos,
					Width:  utils.TileWidth,
					Height: utils.TileHeight,
				}
				if hasCollision(dx, dy, collisionBody, tileCollision) {
					return true
				}
			}
		}
	}
	return false
}

func hasCollision(dx, dy int, bodyA, bodyB model.CollisionBody) bool {
	// check if movement of bodyA collides with bodyB
	aBounds := collision.BoundingBox{
		X:      float64(bodyA.X + dx),
		Y:      float64(bodyA.Y + dy),
		Width:  float64(bodyA.Width),
		Height: float64(bodyA.Height),
	}
	bBounds := collision.BoundingBox{
		X:      float64(bodyB.X),
		Y:      float64(bodyB.Y),
		Width:  float64(bodyB.Width),
		Height: float64(bodyB.Height),
	}
	if collision.AABBCollision(aBounds, bBounds) {
		return true
	}
	return false
}

func isClicked(x, y int, body model.SpriteBody) bool {
	// check if mouse clicked on a body
	aBounds := collision.BoundingBox{
		X:      float64(x),
		Y:      float64(y),
		Width:  1,
		Height: 1,
	}
	bBounds := collision.BoundingBox{
		X:      float64(body.X),
		Y:      float64(body.Y),
		Width:  float64(body.Width),
		Height: float64(body.Height),
	}
	if collision.AABBCollision(aBounds, bBounds) {
		return true
	}
	return false
}

func isAtExit(g *Game, exitPoint int) bool {
	exitToAnimalMap := g.Environment.Maps[g.CurrentMap].Groups[0].ObjectGroups[exitPoint].Objects[0]
	pointCollision := model.CollisionBody{
		X:      int(exitToAnimalMap.X),
		Y:      int(exitToAnimalMap.Y),
		Width:  1,
		Height: 1,
	}
	if hasCollision(g.Player.Dx, g.Player.Dy, g.Player.Collision, pointCollision) {
		return true
	}
	return false
}

func changeMap(g *Game, newMap, entryPoint int) {
	g.CurrentMap = newMap
	point := g.Environment.Maps[newMap].Groups[0].ObjectGroups[entryPoint].Objects[0]
	g.Player.ChangeLocation(int(point.X), int(point.Y))

}

func playerHasCollisions(g *Game) bool {
	if hasMapCollisions(g, g.Player.Dx, g.Player.Dy, g.Player.Collision) {
		return true
	}

	// check for exits
	if g.CurrentMap == utils.FarmMap {
		if isAtExit(g, utils.FarmMapExitToAnimalMapPoint) {
			changeMap(g, utils.AnimalsMap, utils.AnimalMapEntryPoint)
		} else if isAtExit(g, utils.FarmMapExitToForestMapPoint) {
			changeMap(g, utils.ForestMap, utils.ForestMapEntryPoint)
		}
	} else if g.CurrentMap == utils.AnimalsMap {
		if isAtExit(g, utils.AnimalMapExitPoint) {
			changeMap(g, utils.FarmMap, utils.FarmMapEntryFromAnimalMapPoint)
		}
	} else if g.CurrentMap == utils.ForestMap {
		if isAtExit(g, utils.ForestMapExitPoint) {
			changeMap(g, utils.FarmMap, utils.FarmMapEntryFromForestMapPoint)
		}
	}

	// check for animated entities collisions
	for _, c := range g.Chickens {
		if hasCollision(g.Player.Dx, g.Player.Dy, g.Player.Collision, c.Collision) {
			return true
		}
	}
	for _, c := range g.Cows {
		if hasCollision(g.Player.Dx, g.Player.Dy, g.Player.Collision, c.Collision) {
			return true
		}
	}

	// check for trees
	for _, t := range g.Environment.Trees {
		if hasCollision(g.Player.Dx, g.Player.Dy, g.Player.Collision, t.Collision) {
			return true
		}
	}

	// check for objects
	for _, o := range g.Environment.Objects[g.CurrentMap] {
		if hasCollision(g.Player.Dx, g.Player.Dy, g.Player.Collision, o.Collision) && o.IsCollision {
			return true
		}
	}
	return false
}

func updateAnimals(g *Game) {
	for i := range g.Chickens {
		g.Chickens[i].UpdateFrame(g.CurrentFrame)
		g.Cows[i].UpdateFrame(g.CurrentFrame)
	}
}
