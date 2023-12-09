package game

import (
	"guion-2d-project3/entity/model"
	"guion-2d-project3/utils"

	"github.com/co0p/tankism/lib/collision"
)

func hasMapCollisions(g *Game, dx, dy int, collisionBody model.CollisionBody) bool {
	for tileY := 0; tileY < utils.MapRows; tileY += 1 {
		for tileX := 0; tileX < utils.MapColumns; tileX += 1 {
			for _, layer := range utils.CollisionLayers {
				tile := g.Maps[g.CurrentMap].Layers[layer].Tiles[tileY*utils.MapColumns+tileX]
				if tile.IsNil() {
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
	player := g.Data.Players[g.PlayerID]
	exitToAnimalMap := g.Maps[g.CurrentMap].Groups[0].ObjectGroups[exitPoint].Objects[0]
	pointCollision := model.CollisionBody{
		X:      int(exitToAnimalMap.X),
		Y:      int(exitToAnimalMap.Y),
		Width:  1,
		Height: 1,
	}
	if hasCollision(player.Dx, player.Dy, player.Collision, pointCollision) {
		return true
	}
	return false
}

func changeMap(g *Game, newMap, entryPoint int) {
	g.CurrentMap = newMap
	point := g.Maps[newMap].Groups[0].ObjectGroups[entryPoint].Objects[0]
	g.Data.Players[g.PlayerID].ChangeLocation(newMap, int(point.X), int(point.Y))

}

func playerHasCollisions(g *Game) bool {
	player := g.Data.Players[g.PlayerID]
	if hasMapCollisions(g, player.Dx, player.Dy, player.Collision) {
		return true
	}

	// check for exits
	if g.CurrentMap == utils.FarmMap {
		if isAtExit(g, utils.FarmMapExitToAnimalMapPoint) {
			g.Sounds.PlaySound(g.Sounds.SFXChangeMap)
			changeMap(g, utils.AnimalsMap, utils.AnimalMapEntryPoint)
		} else if isAtExit(g, utils.FarmMapExitToForestMapPoint) {
			g.Sounds.PlaySound(g.Sounds.SFXChangeMap)
			changeMap(g, utils.ForestMap, utils.ForestMapEntryPoint)
		}
	} else if g.CurrentMap == utils.AnimalsMap {
		if isAtExit(g, utils.AnimalMapExitPoint) {
			g.Sounds.PlaySound(g.Sounds.SFXChangeMap)
			changeMap(g, utils.FarmMap, utils.FarmMapEntryFromAnimalMapPoint)
		}
	} else if g.CurrentMap == utils.ForestMap {
		if isAtExit(g, utils.ForestMapExitPoint) {
			g.Sounds.PlaySound(g.Sounds.SFXChangeMap)
			changeMap(g, utils.FarmMap, utils.FarmMapEntryFromForestMapPoint)
		}
	}

	// check for animated entities collisions
	if g.CurrentMap == utils.AnimalsMap {
		for _, c := range g.Data.Chickens {
			if hasCollision(player.Dx, player.Dy, player.Collision, c.Collision) {
				return true
			}
		}
		for _, c := range g.Data.Cows {
			if hasCollision(player.Dx, player.Dy, player.Collision, c.Collision) {
				return true
			}
		}
	}

	// check for trees
	if g.CurrentMap == utils.ForestMap {
		for _, t := range g.Data.Environment.Trees {
			if t.IsNil {
				continue
			}
			if hasCollision(player.Dx, player.Dy, player.Collision, t.Collision) {
				return true
			}
		}
	}

	// check for objects
	for _, o := range g.Data.Environment.Objects[g.CurrentMap] {
		if o.IsNil {
			continue
		}
		if hasCollision(player.Dx, player.Dy, player.Collision, o.Collision) && o.IsCollision {
			return true
		}
	}
	return false
}

func chickenHasCollisions(g *Game, chicken int) bool {
	if hasMapCollisions(g, g.Data.Chickens[chicken].Dx, g.Data.Chickens[chicken].Dy, g.Data.Chickens[chicken].Collision) {
		return true
	}
	for _, cow := range g.Data.Cows {
		if hasCollision(g.Data.Chickens[chicken].Dx, g.Data.Chickens[chicken].Dy, g.Data.Chickens[chicken].Collision,
			cow.Collision) {
			return true
		}
	}
	return false
}
