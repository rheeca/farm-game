package game

import (
	"guion-2d-project3/entity/model"
	"guion-2d-project3/utils"
)

func calculateTargetTile(g *Game) (tileX, tileY int) {
	tileX = g.Player.Collision.X / utils.TileWidth
	tileY = (g.Player.Collision.Y + g.Player.Collision.Height) / utils.TileHeight
	if g.Player.Direction == utils.Front {
		tileY += 1
	} else if g.Player.Direction == utils.Back {
		tileY -= 1
	} else if g.Player.Direction == utils.Left {
		tileX -= 1
	} else if g.Player.Direction == utils.Right {
		tileX += 1
	}
	return tileX, tileY
}

func isTile(g *Game, tileX, tileY, tileID int, tileset string) bool {
	if (int(g.Maps[g.CurrentMap].Layers[utils.GroundLayer].Tiles[tileY*utils.MapColumns+tileX].ID) == tileID) &&
		(g.Maps[g.CurrentMap].Layers[utils.GroundLayer].Tiles[tileY*utils.MapColumns+tileX].Tileset.Name ==
			tileset) {
		return true
	} else {
		return false
	}
}

func isMapObject(g *Game, tileX, tileY, tileID int, tileset string) bool {
	if (int(g.Maps[g.CurrentMap].Layers[utils.ObjectsLayer].Tiles[tileY*utils.MapColumns+tileX].ID) == tileID) &&
		(g.Maps[g.CurrentMap].Layers[utils.ObjectsLayer].Tiles[tileY*utils.MapColumns+tileX].Tileset.Name ==
			tileset) {
		return true
	} else {
		return false
	}
}

func isFarmLand(g *Game, tileX, tileY int) bool {
	if g.CurrentMap != utils.FarmMap {
		return false
	}
	if !g.Maps[utils.FarmMap].Layers[utils.FarmingLandLayer].Tiles[tileY*utils.MapColumns+tileX].IsNil() &&
		!hasMapCollisions(g, 0, 0, model.CollisionBody{X: tileX * utils.TileWidth, Y: tileY * utils.TileHeight,
			Width: utils.UnitSize, Height: utils.UnitSize}) {
		return true
	} else {
		return false
	}
}
