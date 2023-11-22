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
				if t.IsNil {
					continue
				}
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
					treeHit := i
					g.Environment.Trees[i].StartAnimation(utils.TreeHitAnimation, utils.FrameCountSix, utils.AnimationDelay,
						doDelayFcn,
						func() {
							g.Player.AddToBackpack(utils.ItemWood2, 5)
							g.Environment.Trees[treeHit].IsNil = true
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
		for i, o := range g.Environment.Objects[g.CurrentMap] {
			if isClicked(mouseX, mouseY, o.Sprite) {
				if o.Type == utils.ItemCraftingTable {
					g.UIState.SelectedRecipe = 0
					g.State = utils.GameStateCraft
				} else if o.Type == utils.ItemDoor {
					if o.IsCollision { // door is currently closed
						g.Environment.Objects[g.CurrentMap][i].StartAnimation(utils.OpenDoorAnimation, utils.FrameCountSix, 0,
							true, func() {
								g.Environment.Objects[g.CurrentMap][i].IsCollision = false
							})
					} else { // door is currently open
						g.Environment.Objects[g.CurrentMap][i].StartAnimation(utils.CloseDoorAnimation, utils.FrameCountSix, 0,
							true, func() {
								g.Environment.Objects[g.CurrentMap][i].IsCollision = true
							})
					}
				} else if o.Type == utils.ItemBedPink {
					g.ShowImage(g.Images.BlackScreen)
					g.Environment.ResetDay()
				} else if o.Type == utils.ItemMapStone3 {
					if g.Player.AddToBackpack(utils.ItemRock1, 1) {
						g.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
				} else if o.Type == utils.ItemMapWood {
					if g.Player.AddToBackpack(utils.ItemWood2, 1) {
						g.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
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
				g.SetErrorMessage("Backpack is full!")
			}
		} else if isMapObject(g, tileX, tileY, utils.MapStone3, utils.TilesetFlowersStones) {
			if g.Player.AddToBackpack(utils.ItemRock1, 1) {
				g.Environment.Maps[g.CurrentMap].Layers[utils.ObjectsLayer].Tiles[tileY*utils.MapColumns+tileX] = &emptyTile
			} else {
				g.SetErrorMessage("Backpack is full!")
			}
		}
	}
}

func getPlayerInput(g *Game) {
	g.Player.UpdateFrame(g.CurrentFrame)

	if g.State == utils.GameStateCustomChar {
		checkMouseOnCustomCharState(g)
		return
	} else if g.State == utils.GameStateCraft {
		checkMouseOnCraftState(g)
		return
	}

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

func checkMouseOnCraftState(g *Game) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		// exit button
		if isClicked(mouseX, mouseY, model.SpriteBody{X: 654, Y: 106, Width: 36, Height: 40}) {
			g.State = utils.GameStatePlay
			return
		}

		// select recipe
		for i := 0; i < utils.CraftingUIBoxCount; i++ {
			recipeBox := model.SpriteBody{
				X:      utils.CraftingUIBoxCollisionX + (utils.CraftingUISpacing * (i % utils.CraftingUIColumns)),
				Y:      utils.CraftingUIBoxCollisionY + (utils.CraftingUISpacing * (i / utils.CraftingUIColumns)),
				Width:  utils.CraftingUIBoxCollisionWidth,
				Height: utils.CraftingUIBoxCollisionHeight,
			}
			if isClicked(mouseX, mouseY, recipeBox) {
				g.UIState.SelectedRecipe = i
			}
		}

		// craft button
		if isClicked(mouseX, mouseY, model.SpriteBody{X: 486, Y: 452, Width: 180, Height: 54}) {
			var items []model.BackpackItem
			for _, item := range utils.RecipeDetails[utils.Recipes[g.UIState.SelectedRecipe]] {
				items = append(items, model.BackpackItem{ID: item.ID, Count: item.Count})
			}
			if g.Player.RemoveFromBackpack(items) {
				g.Player.AddToBackpack(utils.Recipes[g.UIState.SelectedRecipe], 1)
			} else {
				g.SetErrorMessage("Not enough materials!")
			}
		}
	}
}

func checkMouseOnCustomCharState(g *Game) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		// select character
		for i := 0; i < utils.CharacterUIBoxCount; i++ {
			recipeBox := model.SpriteBody{
				X:      utils.CharacterUIBoxCollisionX + (utils.CharacterUISpacing * (i % utils.CharacterUIColumns)),
				Y:      utils.CharacterUIBoxCollisionY + (utils.CharacterUISpacing * (i / utils.CharacterUIColumns)),
				Width:  utils.CharacterUIBoxCollisionWidth,
				Height: utils.CharacterUIBoxCollisionHeight,
			}
			if isClicked(mouseX, mouseY, recipeBox) {
				g.UIState.SelectedCharacter = i
			}
		}

		// play button
		if isClicked(mouseX, mouseY, model.SpriteBody{X: 294, Y: 387, Width: 212, Height: 55}) {
			g.Player.Spritesheet = g.Images.Characters[g.UIState.SelectedCharacter]
			g.State = utils.GameStatePlay
		}
	}
}
