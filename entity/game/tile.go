package game

import (
	"guion-2d-project3/utils"
)

func calculateTargetTile(g *Game) (tileX, tileY int) {
	tileX = g.Player.Collision.X0 / utils.TileWidth
	tileY = g.Player.Collision.Y1 / utils.TileHeight
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
	if (int(g.Environment.Maps[g.CurrentMap].Layers[utils.GroundLayer].Tiles[tileY*utils.MapColumns+tileX].ID) == tileID) &&
		(g.Environment.Maps[g.CurrentMap].Layers[utils.GroundLayer].Tiles[tileY*utils.MapColumns+tileX].Tileset.Name ==
			tileset) {
		return true
	} else {
		return false
	}
}

func isMapObject(g *Game, tileX, tileY, tileID int, tileset string) bool {
	if (int(g.Environment.Maps[g.CurrentMap].Layers[utils.ObjectsLayer].Tiles[tileY*utils.MapColumns+tileX].ID) == tileID) &&
		(g.Environment.Maps[g.CurrentMap].Layers[utils.ObjectsLayer].Tiles[tileY*utils.MapColumns+tileX].Tileset.Name ==
			tileset) {
		return true
	} else {
		return false
	}
}
