package game

import (
	"guion-2d-project3/entity/model"
	"guion-2d-project3/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lafriks/go-tiled"
)

func getPlayerInput(g *Game) {
	if g.State == utils.GameStateCustomChar {
		checkMouseOnCustomCharState(g)
	} else if g.State == utils.GameStateCraft {
		checkMouseOnCraftState(g)
	} else if g.State == utils.GameStatePlay {
		checkMouseOnPlayState(g)
		checkKeyboardOnPlayState(g)
	}
}

func checkMouseOnPlayState(g *Game) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		// select item in backpack
		for i := 0; i < utils.BackpackSize; i++ {
			if isClicked(mouseX, mouseY, model.SpriteBody{
				X:      utils.ToolsFirstBoxX + (i * utils.BackpackUIBoxWidth),
				Y:      utils.ToolsFirstBoxY,
				Width:  utils.BackpackUIBoxWidth,
				Height: utils.BackpackUIBoxWidth,
			}) {
				g.Players[g.PlayerID].EquippedItem = i
				return
			}
		}

		// delete item in backpack
		if isClicked(mouseX, mouseY, model.SpriteBody{
			X:      utils.BackpackDeleteButtonX,
			Y:      utils.BackpackDeleteButtonY,
			Width:  utils.BackpackDeleteButtonWidth,
			Height: utils.BackpackDeleteButtonHeight,
		}) {
			g.Players[g.PlayerID].RemoveFromBackpackByIndex(g.Players[g.PlayerID].EquippedItem)
		}

		// use tool
		if g.Players[g.PlayerID].Backpack[g.Players[g.PlayerID].EquippedItem].ID == utils.ItemHoe {
			g.Sounds.PlaySound(g.Sounds.SFXTillSoil)
			g.Players[g.PlayerID].State = utils.HoeState
			g.Players[g.PlayerID].Frame = 0
			g.Players[g.PlayerID].StateTTL = utils.PlayerFrameCount

			tileX, tileY := calculateTargetTile(g)
			if isFarmLand(g, tileX, tileY) {
				g.Data.Environment.AddPlot(tileX, tileY)
			}
		} else if g.Players[g.PlayerID].Backpack[g.Players[g.PlayerID].EquippedItem].ID == utils.ItemAxe {
			g.Players[g.PlayerID].State = utils.AxeState
			g.Players[g.PlayerID].Frame = 0
			g.Players[g.PlayerID].StateTTL = utils.PlayerFrameCount

			if g.CurrentMap == utils.ForestMap {

				for i, t := range g.Data.Environment.Trees {
					if t.IsNil {
						continue
					}
					// if tree is in target, chop tree
					if hasCollision(0, 0, g.Players[g.PlayerID].CalcTargetBox(), t.Collision) {
						// if tree health reaches zero, set the delay function to be executed after the animation
						g.Data.Environment.Trees[i].Health -= 1
						var doDelayFcn bool
						if g.Data.Environment.Trees[i].Health <= 0 {
							doDelayFcn = true
						} else {
							doDelayFcn = false
						}
						treeHit := i
						g.Sounds.PlaySound(g.Sounds.SFXChopTree)
						g.Data.Environment.Trees[i].StartAnimation(utils.TreeHitAnimation, utils.FrameCountSix, utils.AnimationDelay,
							doDelayFcn,
							func() {
								g.Players[g.PlayerID].AddToBackpack(utils.ItemWood2, 5)
								g.Data.Environment.Trees[treeHit].IsNil = true
							})
					}
				}
			}
		} else if g.Players[g.PlayerID].Backpack[g.Players[g.PlayerID].EquippedItem].ID == utils.ItemWateringCan {
			g.Sounds.PlaySound(g.Sounds.SFXWateringCan)
			g.Players[g.PlayerID].State = utils.WateringState
			g.Players[g.PlayerID].Frame = 0
			g.Players[g.PlayerID].StateTTL = utils.PlayerFrameCount

			tileX, tileY := calculateTargetTile(g)
			if isFarmLand(g, tileX, tileY) {
				g.Data.Environment.WaterPlot(tileX, tileY)
			}
		} else if utils.IsSeed(g.Players[g.PlayerID].Backpack[g.Players[g.PlayerID].EquippedItem].ID) {
			tileX, tileY := calculateTargetTile(g)
			if isFarmLand(g, tileX, tileY) {
				if g.Data.Environment.PlantSeedInPlot(tileX, tileY, utils.PlantTomato) {
					g.Players[g.PlayerID].RemoveFromBackpackByIndexAndCount(g.Players[g.PlayerID].EquippedItem, 1)
				}
			}
		}
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		tileX, tileY := calculateTargetTile(g)
		mouseX, mouseY := ebiten.CursorPosition()

		// if target tile is an object
		for i, o := range g.Data.Environment.Objects[g.CurrentMap] {
			if isClicked(mouseX, mouseY, o.Sprite) {
				if o.Type == utils.ItemCraftingTable {
					g.UIState.SelectedRecipe = 0
					g.State = utils.GameStateCraft
				} else if o.Type == utils.ItemDoor {
					if o.IsCollision { // door is currently closed
						g.Sounds.PlaySound(g.Sounds.SFXOpenDoor)
						g.Data.Environment.Objects[g.CurrentMap][i].StartAnimation(utils.OpenDoorAnimation, utils.FrameCountSix, 0,
							true, func() {
								g.Data.Environment.Objects[g.CurrentMap][i].IsCollision = false
							})
					} else { // door is currently open
						g.Sounds.PlaySound(g.Sounds.SFXCloseDoor)
						g.Data.Environment.Objects[g.CurrentMap][i].StartAnimation(utils.CloseDoorAnimation, utils.FrameCountSix, 0,
							true, func() {
								g.Data.Environment.Objects[g.CurrentMap][i].IsCollision = true
							})
					}
				} else if o.Type == utils.ItemBedPink {
					g.ShowImage(g.Images.BlackScreen)
					g.Data.Environment.ResetDay()
				} else if o.Type == utils.ItemMapStone3 {
					if g.Players[g.PlayerID].AddToBackpack(utils.ItemRock1, 1) {
						g.Data.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
				} else if o.Type == utils.ItemMapWood {
					if g.Players[g.PlayerID].AddToBackpack(utils.ItemWood2, 1) {
						g.Data.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
				} else if o.Type == utils.MapSunflower {
					if g.Players[g.PlayerID].AddToBackpack(utils.ItemSunflower, 1) {
						g.Data.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
				} else if o.Type == utils.MapBlueflower {
					if g.Players[g.PlayerID].AddToBackpack(utils.ItemBlueflower, 1) {
						g.Data.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
				} else if o.Type == utils.MapWeed {
					if g.Players[g.PlayerID].AddToBackpack(utils.ItemWeed, 1) {
						g.Data.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
				} else if o.Type == utils.MapPinkDyeFlower {
					if g.Players[g.PlayerID].AddToBackpack(utils.ItemPinkDyeFlower, 1) {
						g.Data.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
				} else if o.Type == utils.MapBlueDyeFlower {
					if g.Players[g.PlayerID].AddToBackpack(utils.ItemBlueDyeFlower, 1) {
						g.Data.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
				}
				return
			}
		}

		// if target tile has an animated character
		if g.CurrentMap == utils.AnimalsMap {
			for _, c := range g.Data.Chickens {
				if isClicked(mouseX, mouseY, c.Sprite) {
					g.Sounds.PlaySound(g.Sounds.SFXChicken)
					c.State = utils.ChickenHeartState
					c.Frame = 0
					c.StateTTL = utils.AnimalFrameCount
				}
			}

			for _, c := range g.Data.Cows {
				if isClicked(mouseX, mouseY, c.Sprite) {
					g.Sounds.PlaySound(g.Sounds.SFXCow)
					c.State = utils.CowHeartState
					c.Frame = 0
					c.AnimationTTL = utils.AnimalFrameCount
				}
			}
		}

		// pick up objects from the map
		emptyTile := tiled.LayerTile{Nil: true}
		if isMapObject(g, tileX, tileY, utils.MapWood, utils.TilesetTrees) {
			if g.Players[g.PlayerID].AddToBackpack(utils.ItemWood2, 1) {
				g.Maps[g.CurrentMap].Layers[utils.ObjectsLayer].Tiles[tileY*utils.MapColumns+tileX] = &emptyTile
			} else {
				g.SetErrorMessage("Backpack is full!")
			}
		} else if isMapObject(g, tileX, tileY, utils.MapStone3, utils.TilesetFlowersStones) {
			if g.Players[g.PlayerID].AddToBackpack(utils.ItemRock1, 1) {
				g.Maps[g.CurrentMap].Layers[utils.ObjectsLayer].Tiles[tileY*utils.MapColumns+tileX] = &emptyTile
			} else {
				g.SetErrorMessage("Backpack is full!")
			}
		}

		// harvest plant
		tileX, tileY = calculateTargetTile(g)
		if isFarmLand(g, tileX, tileY) {
			hasHarvest, plantType := g.Data.Environment.HarvestPlant(tileX, tileY)
			if hasHarvest {
				g.Players[g.PlayerID].AddToBackpack(utils.PlantItemMapping[plantType], 1)
			}
		}
	}
}

func checkKeyboardOnPlayState(g *Game) {
	// player movement
	if ebiten.IsKeyPressed(ebiten.KeyA) && g.Players[g.PlayerID].Sprite.X > 0 {
		g.Players[g.PlayerID].Direction = utils.Left
		g.Players[g.PlayerID].State = utils.WalkState
		g.Players[g.PlayerID].Dx -= utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Players[g.PlayerID].UpdateLocation()
		} else {
			g.Players[g.PlayerID].Dx = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyD) &&
		g.Players[g.PlayerID].Sprite.X < utils.MapWidth-g.Players[g.PlayerID].Sprite.Width {
		g.Players[g.PlayerID].Direction = utils.Right
		g.Players[g.PlayerID].State = utils.WalkState
		g.Players[g.PlayerID].Dx += utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Players[g.PlayerID].UpdateLocation()
		} else {
			g.Players[g.PlayerID].Dx = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyW) && g.Players[g.PlayerID].Sprite.Y > 0 {
		g.Players[g.PlayerID].Direction = utils.Back
		g.Players[g.PlayerID].State = utils.WalkState
		g.Players[g.PlayerID].Dy -= utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Players[g.PlayerID].UpdateLocation()
		} else {
			g.Players[g.PlayerID].Dy = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyS) &&
		g.Players[g.PlayerID].Sprite.Y < utils.MapHeight-g.Players[g.PlayerID].Sprite.Height {
		g.Players[g.PlayerID].Direction = utils.Front
		g.Players[g.PlayerID].State = utils.WalkState
		g.Players[g.PlayerID].Dy += utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Players[g.PlayerID].UpdateLocation()
		} else {
			g.Players[g.PlayerID].Dy = 0
		}
	} else if g.Players[g.PlayerID].StateTTL == 0 {
		g.Players[g.PlayerID].State = utils.IdleState
	}

	// equip item
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.Players[g.PlayerID].EquippedItem = 0
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		g.Players[g.PlayerID].EquippedItem = 1
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		g.Players[g.PlayerID].EquippedItem = 2
	} else if ebiten.IsKeyPressed(ebiten.Key4) {
		g.Players[g.PlayerID].EquippedItem = 3
	} else if ebiten.IsKeyPressed(ebiten.Key5) {
		g.Players[g.PlayerID].EquippedItem = 4
	} else if ebiten.IsKeyPressed(ebiten.Key6) {
		g.Players[g.PlayerID].EquippedItem = 5
	} else if ebiten.IsKeyPressed(ebiten.Key7) {
		g.Players[g.PlayerID].EquippedItem = 6
	} else if ebiten.IsKeyPressed(ebiten.Key8) {
		g.Players[g.PlayerID].EquippedItem = 7
	} else if ebiten.IsKeyPressed(ebiten.Key9) {
		g.Players[g.PlayerID].EquippedItem = 8
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
			recipe := utils.RecipeDetails[utils.Recipes[g.UIState.SelectedRecipe]]
			for _, item := range recipe.Materials {
				items = append(items, model.BackpackItem{ID: item.ID, Count: item.Count})
			}
			if g.Players[g.PlayerID].RemoveFromBackpack(items) {
				g.Sounds.PlaySound(g.Sounds.SFXCraft)
				g.Players[g.PlayerID].AddToBackpack(utils.Recipes[g.UIState.SelectedRecipe], recipe.Count)
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
			g.Players[g.PlayerID].Spritesheet = g.Images.Characters[g.UIState.SelectedCharacter]
			g.State = utils.GameStatePlay
		}
	}
}
