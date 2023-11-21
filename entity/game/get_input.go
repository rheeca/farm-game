package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lafriks/go-tiled"
	"guion-2d-project3/entity/model"
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

			for i, t := range g.Environment.Trees {
				// if tree is in target, chop tree
				if hasCollision(0, 0, g.Player.CalcTargetBox(), t.Collision) {
					// if tree health reaches zero, set the delay function to be executed after the animation
					g.Environment.Trees[i].Health -= 1
					var doDelayFcn bool
					if g.Environment.Trees[i].Health <= 0 {
						doDelayFcn = true
					} else {
						doDelayFcn = false
					}

					g.Environment.Trees[i].StartAnimation(utils.TreeHitAnimation, utils.FrameCountSix, utils.AnimationDelay,
						doDelayFcn,
						func() {
							g.Player.AddToBackpack(utils.ItemWood2, 5)
							// remove chopped trees
							var newTrees []model.Object
							for _, t := range g.Environment.Trees {
								if t.Health > 0 {
									newTrees = append(newTrees, t)
								}
							}
							g.Environment.Trees = newTrees
						})
				}
			}
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
		mouseX, mouseY := ebiten.CursorPosition()

		// if target tile is an object
		for _, o := range g.Environment.Objects[g.CurrentMap] {
			if isClicked(mouseX, mouseY, o.Sprite) {
				if o.Type == utils.ItemCraftingTable {
					g.State = utils.GameStateCraft
				}
				return
			}
		}

		// if target tile has an animated character
		for _, c := range g.Chickens {
			if isClicked(mouseX, mouseY, c.Sprite) {
				c.State = utils.ChickenHeartState
				c.Frame = 0
				c.AnimationTTL = utils.AnimalFrameCount
			}
		}

		for _, c := range g.Cows {
			if isClicked(mouseX, mouseY, c.Sprite) {
				c.State = utils.CowHeartState
				c.Frame = 0
				c.AnimationTTL = utils.AnimalFrameCount
			}
		}

		// pick up objects from the map
		emptyTile := tiled.LayerTile{Nil: true}
		if isMapObject(g, tileX, tileY, utils.MapWood, utils.TilesetTrees) {
			if g.Player.AddToBackpack(utils.ItemWood2, 1) {
				g.Environment.Maps[g.CurrentMap].Layers[utils.ObjectsLayer].Tiles[tileY*utils.MapColumns+tileX] = &emptyTile
			} else {
				// TODO: alert player that backpack is full
			}
		} else if isMapObject(g, tileX, tileY, utils.MapStone3, utils.TilesetFlowersStones) {
			if g.Player.AddToBackpack(utils.ItemRock1, 1) {
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
	if ebiten.IsKeyPressed(ebiten.KeyA) && g.Player.Sprite.X > 0 {
		g.Player.Direction = utils.Left
		g.Player.State = utils.WalkState
		g.Player.Dx -= utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dx = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyD) &&
		g.Player.Sprite.X < utils.MapWidth-g.Player.Sprite.Width {
		g.Player.Direction = utils.Right
		g.Player.State = utils.WalkState
		g.Player.Dx += utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dx = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyW) && g.Player.Sprite.Y > 0 {
		g.Player.Direction = utils.Back
		g.Player.State = utils.WalkState
		g.Player.Dy -= utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dy = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyS) &&
		g.Player.Sprite.Y < utils.MapHeight-g.Player.Sprite.Height {
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
