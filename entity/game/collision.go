package game

import (
	"github.com/co0p/tankism/lib/collision"
	"guion-2d-project3/interfaces"
	"guion-2d-project3/utils"
)

func hasMapCollisions(g *Game, animObj interfaces.AnimatedSprite) bool {
	for tileY := 0; tileY < utils.MapRows; tileY += 1 {
		for tileX := 0; tileX < utils.MapColumns; tileX += 1 {
			tile := g.Environment.Maps[0].Layers[utils.ObjectsLayer].Tiles[tileY*utils.MapColumns+tileX]
			if tile.ID == 0 {
				continue
			}
			tileXpos := utils.TileWidth * tileX
			tileYpos := utils.TileHeight * tileY

			newX := animObj.GetXLoc() + animObj.GetDx()
			newY := animObj.GetYLoc() + animObj.GetDy()
			animBounds := collision.BoundingBox{
				// bounding box for animated object made slightly smaller than the sprite
				X:      float64(newX + animObj.GetWidth()/4),
				Y:      float64(newY + animObj.GetHeight()/2),
				Width:  float64(animObj.GetWidth() / 2),
				Height: float64(animObj.GetHeight() / 2),
			}
			tileBounds := collision.BoundingBox{
				X:      float64(tileXpos),
				Y:      float64(tileYpos),
				Width:  float64(utils.TileWidth),
				Height: float64(utils.TileHeight),
			}
			if collision.AABBCollision(animBounds, tileBounds) {
				return true
			}
		}
	}
	return false
}

func playerHasCollisions(g *Game) bool {
	// TODO: check for map collisions
	//if hasMapCollisions(g, g.Player) {
	//	return true
	//}

	// check for animated entities collisions
	for _, a := range g.Animals {
		if g.Player.HasCollisionWith(a) {
			return true
		}
	}
	return false
}

func animalHasCollisions(g *Game, animObj interfaces.AnimatedSprite) bool {
	// TODO: check for map collisions
	//if hasMapCollisions(g, animObj) {
	//	return true
	//}

	// check for collision with player
	if animObj.HasCollisionWith(g.Player) {
		return true
	}
	return false
}

func updateAnimals(g *Game) {
	for i, a := range g.Animals {
		if a.XLoc == (a.Path[a.Destination].X*utils.TileWidth) &&
			a.YLoc == (a.Path[a.Destination].Y*utils.TileHeight) {

			// if animal has reached its destination, give it a new destination
			g.Animals[i].Destination = (g.Animals[i].Destination + 1) % len(a.Path)
		} else {
			// move animal towards destination
			if a.XLoc > a.Path[a.Destination].X*utils.TileWidth {
				g.Animals[i].Direction = utils.Left
				g.Animals[i].UpdateFrame(g.CurrentFrame)
				g.Animals[i].Dx -= utils.AnimalMovementSpeed
				if !animalHasCollisions(g, g.Animals[i]) {
					g.Animals[i].UpdateLocation()
				} else {
					g.Animals[i].Dx = 0
				}
			} else if a.XLoc < a.Path[a.Destination].X*utils.TileWidth {
				g.Animals[i].Direction = utils.Right
				g.Animals[i].UpdateFrame(g.CurrentFrame)
				g.Animals[i].Dx += utils.AnimalMovementSpeed
				if !animalHasCollisions(g, g.Animals[i]) {
					g.Animals[i].UpdateLocation()
				} else {
					g.Animals[i].Dx = 0
				}
			} else if a.YLoc > a.Path[a.Destination].Y*utils.TileHeight {
				g.Animals[i].Direction = utils.Back
				g.Animals[i].UpdateFrame(g.CurrentFrame)
				g.Animals[i].Dy -= utils.AnimalMovementSpeed
				if !animalHasCollisions(g, g.Animals[i]) {
					g.Animals[i].UpdateLocation()
				} else {
					g.Animals[i].Dy = 0
				}
			} else if a.YLoc < a.Path[a.Destination].Y*utils.TileHeight {
				g.Animals[i].Direction = utils.Front
				g.Animals[i].UpdateFrame(g.CurrentFrame)
				g.Animals[i].Dy += utils.AnimalMovementSpeed
				if !animalHasCollisions(g, g.Animals[i]) {
					g.Animals[i].UpdateLocation()
				} else {
					g.Animals[i].Dy = 0
				}
			}
		}
	}
}
