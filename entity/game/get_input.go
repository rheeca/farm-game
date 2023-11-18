package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lafriks/go-tiled"
	"guion-2d-project3/utils"
)

func checkMouse(g *Game) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.Player.Backpack[g.Player.EquippedItem].ID == utils.ItemHoe {
			g.Player.State = utils.HoeState
			g.Player.Frame = 0
			g.Player.StateTTL = utils.PlayerFrameCount

			tileX, tileY := calculateTargetTile(g)
			if isTile(g, tileX, tileY, 12, utils.TilesetGrassHill) {
				// if target tile is a grass tile, make tile into tilled ground
				g.Environment.Maps[g.CurrentMap].Layers[utils.GroundLayer].Tiles[tileY*utils.MapColumns+tileX].ID = 12
				g.Environment.Maps[g.CurrentMap].Layers[utils.GroundLayer].Tiles[tileY*utils.MapColumns+tileX].Tileset =
					g.Environment.Maps[g.CurrentMap].Tilesets[1]
			}
		} else if g.Player.Backpack[g.Player.EquippedItem].ID == utils.ItemAxe {
			g.Player.State = utils.AxeState
			g.Player.Frame = 0
			g.Player.StateTTL = utils.PlayerFrameCount
		} else if g.Player.Backpack[g.Player.EquippedItem].ID == utils.ItemWateringCan {
			g.Player.State = utils.WateringState
			g.Player.Frame = 0
			g.Player.StateTTL = utils.PlayerFrameCount

			tileX, tileY := calculateTargetTile(g)
			if isTile(g, tileX, tileY, 12, utils.TilesetSoilGround) {
				// if target tile is tilled ground, make tile into watered ground
				g.Environment.Maps[g.CurrentMap].Layers[utils.GroundLayer].Tiles[tileY*utils.MapColumns+tileX].ID = 12
				g.Environment.Maps[g.CurrentMap].Layers[utils.GroundLayer].Tiles[tileY*utils.MapColumns+tileX].Tileset =
					g.Environment.Maps[g.CurrentMap].Tilesets[4]
			}
		}
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		tileX, tileY := calculateTargetTile(g)

		// pick up objects from the map
		emptyTile := tiled.LayerTile{Nil: true}
		if isMapObject(g, tileX, tileY, utils.MapWood, utils.TilesetTrees) {
			if g.Player.AddToBackpack(utils.ItemWood2) {
				g.Environment.Maps[g.CurrentMap].Layers[utils.ObjectsLayer].Tiles[tileY*utils.MapColumns+tileX] = &emptyTile
			} else {
				// TODO: alert player that backpack is full
			}
		} else if isMapObject(g, tileX, tileY, utils.MapStone3, utils.TilesetFlowersStones) {
			if g.Player.AddToBackpack(utils.ItemRock1) {
				g.Environment.Maps[g.CurrentMap].Layers[utils.ObjectsLayer].Tiles[tileY*utils.MapColumns+tileX] = &emptyTile
			} else {
				// TODO: alert player that backpack is full
			}
		}
	}
}

func getPlayerInput(g *Game) {
	g.Player.UpdateFrame(g.CurrentFrame)

	checkMouse(g)
	if ebiten.IsKeyPressed(ebiten.KeyA) && g.Player.XLoc > 0 {
		g.Player.Direction = utils.Left
		g.Player.State = utils.WalkState
		g.Player.Dx -= utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dx = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyD) &&
		g.Player.XLoc < utils.MapWidth-g.Player.SpriteWidth {
		g.Player.Direction = utils.Right
		g.Player.State = utils.WalkState
		g.Player.Dx += utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dx = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyW) && g.Player.YLoc > 0 {
		g.Player.Direction = utils.Back
		g.Player.State = utils.WalkState
		g.Player.Dy -= utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dy = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyS) &&
		g.Player.YLoc < utils.MapHeight-g.Player.SpriteHeight {
		g.Player.Direction = utils.Front
		g.Player.State = utils.WalkState
		g.Player.Dy += utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dy = 0
		}
	} else if g.Player.StateTTL == 0 {
		g.Player.State = utils.IdleState
	}

	// Equip item
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.Player.EquippedItem = 0
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		g.Player.EquippedItem = 1
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		g.Player.EquippedItem = 2
	} else if ebiten.IsKeyPressed(ebiten.Key4) {
		g.Player.EquippedItem = 3
	} else if ebiten.IsKeyPressed(ebiten.Key5) {
		g.Player.EquippedItem = 4
	} else if ebiten.IsKeyPressed(ebiten.Key6) {
		g.Player.EquippedItem = 5
	} else if ebiten.IsKeyPressed(ebiten.Key7) {
		g.Player.EquippedItem = 6
	} else if ebiten.IsKeyPressed(ebiten.Key8) {
		g.Player.EquippedItem = 7
	} else if ebiten.IsKeyPressed(ebiten.Key9) {
		g.Player.EquippedItem = 8
	}
}
