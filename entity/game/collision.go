package game

import (
	"guion-2d-project3/entity/model"
	"guion-2d-project3/entity/player"
	"guion-2d-project3/utils"

	"github.com/co0p/tankism/lib/collision"
)

func hasMapCollisions(g *Game, currentMap, dx, dy int, collisionBody model.CollisionBody) bool {
	for tileY := 0; tileY < utils.MapRows; tileY += 1 {
		for tileX := 0; tileX < utils.MapColumns; tileX += 1 {
			for _, layer := range utils.CollisionLayers {
				tile := g.Maps[currentMap].Layers[layer].Tiles[tileY*utils.MapColumns+tileX]
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

func isAtExit(g *Game, player *player.Player, exitPoint int) bool {
	exitToAnimalMap := g.Maps[player.CurrentMap].Groups[0].ObjectGroups[exitPoint].Objects[0]
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

func changeMap(g *Game, player *player.Player, newMap, entryPoint int) {
	if player.PlayerID == g.PlayerID {
		g.CurrentMap = newMap
	}
	point := g.Maps[newMap].Groups[0].ObjectGroups[entryPoint].Objects[0]
	player.ChangeLocation(newMap, int(point.X), int(point.Y))

}

func playerHasCollisions(g *Game, player *player.Player) bool {
	if hasMapCollisions(g, player.CurrentMap, player.Dx, player.Dy, player.Collision) {
		return true
	}

	// check for exits
	if player.CurrentMap == utils.FarmMap {
		if isAtExit(g, player, utils.FarmMapExitToAnimalMapPoint) {
			g.Sounds.PlaySound(g.Sounds.SFXChangeMap)
			changeMap(g, player, utils.AnimalsMap, utils.AnimalMapEntryPoint)
		} else if isAtExit(g, player, utils.FarmMapExitToForestMapPoint) {
			g.Sounds.PlaySound(g.Sounds.SFXChangeMap)
			changeMap(g, player, utils.ForestMap, utils.ForestMapEntryPoint)
		}
	} else if player.CurrentMap == utils.AnimalsMap {
		if isAtExit(g, player, utils.AnimalMapExitPoint) {
			g.Sounds.PlaySound(g.Sounds.SFXChangeMap)
			changeMap(g, player, utils.FarmMap, utils.FarmMapEntryFromAnimalMapPoint)
		}
	} else if player.CurrentMap == utils.ForestMap {
		if isAtExit(g, player, utils.ForestMapExitPoint) {
			g.Sounds.PlaySound(g.Sounds.SFXChangeMap)
			changeMap(g, player, utils.FarmMap, utils.FarmMapEntryFromForestMapPoint)
		}
	}

	// check for animated entities collisions
	if player.CurrentMap == utils.AnimalsMap {
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
	if player.CurrentMap == utils.ForestMap {
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
	for _, o := range g.Data.Environment.Objects[player.CurrentMap] {
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
	if hasMapCollisions(g, utils.AnimalsMap, g.Data.Chickens[chicken].Dx, g.Data.Chickens[chicken].Dy, g.Data.Chickens[chicken].Collision) {
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
