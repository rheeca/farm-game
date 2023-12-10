package game

import (
	"guion-2d-project3/entity/model"
	"guion-2d-project3/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) UpdateClientInput(clientInput model.ClientInputPacket) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.clientInputs[clientInput.PlayerID] = clientInput
}

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

func getClientInputs(g *Game) {
	g.lock.Lock()
	defer g.lock.Unlock()
	for _, p := range g.Data.Players {
		if g.clientInputs[p.PlayerID].Input == utils.InputKeyW {
			p.Direction = utils.Back
			p.State = utils.WalkState
			p.Dy -= utils.MovementSpeed
			p.UpdateLocation()
		} else if g.clientInputs[p.PlayerID].Input == utils.InputKeyA {
			p.Direction = utils.Left
			p.State = utils.WalkState
			p.Dx -= utils.MovementSpeed
			p.UpdateLocation()
		} else if g.clientInputs[p.PlayerID].Input == utils.InputKeyS {
			p.Direction = utils.Front
			p.State = utils.WalkState
			p.Dy += utils.MovementSpeed
			p.UpdateLocation()
		} else if g.clientInputs[p.PlayerID].Input == utils.InputKeyD {
			p.Direction = utils.Right
			p.State = utils.WalkState
			p.Dx += utils.MovementSpeed
			p.UpdateLocation()
		} else if p.StateTTL == 0 {
			p.State = utils.IdleState
		}
	}
}

func checkMouseOnPlayState(g *Game) {
	player := g.Data.Players[g.PlayerID]
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
				player.EquippedItem = i
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
			player.RemoveFromBackpackByIndex(player.EquippedItem)
		}

		// use tool
		if player.Backpack[player.EquippedItem].ID == utils.ItemHoe {
			g.Sounds.PlaySound(g.Sounds.SFXTillSoil)
			player.State = utils.HoeState
			player.Frame = 0
			player.StateTTL = utils.PlayerFrameCount

			tileX, tileY := calculateTargetTile(g)
			if isFarmLand(g, tileX, tileY) {
				g.Data.Environment.AddPlot(tileX, tileY)
			}
		} else if player.Backpack[player.EquippedItem].ID == utils.ItemAxe {
			player.State = utils.AxeState
			player.Frame = 0
			player.StateTTL = utils.PlayerFrameCount

			if g.CurrentMap == utils.ForestMap {

				for i, t := range g.Data.Environment.Trees {
					if t.IsNil {
						continue
					}
					// if tree is in target, chop tree
					if hasCollision(0, 0, player.CalcTargetBox(), t.Collision) {
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
								player.AddToBackpack(utils.ItemWood2, 5)
								g.Data.Environment.Trees[treeHit].IsNil = true
							})
					}
				}
			}
		} else if player.Backpack[player.EquippedItem].ID == utils.ItemWateringCan {
			g.Sounds.PlaySound(g.Sounds.SFXWateringCan)
			player.State = utils.WateringState
			player.Frame = 0
			player.StateTTL = utils.PlayerFrameCount

			tileX, tileY := calculateTargetTile(g)
			if isFarmLand(g, tileX, tileY) {
				g.Data.Environment.WaterPlot(tileX, tileY)
			}
		} else if utils.IsSeed(player.Backpack[player.EquippedItem].ID) {
			tileX, tileY := calculateTargetTile(g)
			if isFarmLand(g, tileX, tileY) {
				if g.Data.Environment.PlantSeedInPlot(tileX, tileY, utils.PlantTomato) {
					player.RemoveFromBackpackByIndexAndCount(player.EquippedItem, 1)
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
					if g.Data.Players[g.PlayerID].AddToBackpack(utils.ItemRock1, 1) {
						g.Data.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
				} else if o.Type == utils.ItemMapWood {
					if player.AddToBackpack(utils.ItemWood2, 1) {
						g.Data.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
				} else if o.Type == utils.MapSunflower {
					if player.AddToBackpack(utils.ItemSunflower, 1) {
						g.Data.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
				} else if o.Type == utils.MapBlueflower {
					if player.AddToBackpack(utils.ItemBlueflower, 1) {
						g.Data.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
				} else if o.Type == utils.MapWeed {
					if player.AddToBackpack(utils.ItemWeed, 1) {
						g.Data.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
				} else if o.Type == utils.MapPinkDyeFlower {
					if player.AddToBackpack(utils.ItemPinkDyeFlower, 1) {
						g.Data.Environment.Objects[g.CurrentMap][i].IsNil = true
					} else {
						g.SetErrorMessage("Backpack is full!")
					}
				} else if o.Type == utils.MapBlueDyeFlower {
					if player.AddToBackpack(utils.ItemBlueDyeFlower, 1) {
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

		// harvest plant
		tileX, tileY = calculateTargetTile(g)
		if isFarmLand(g, tileX, tileY) {
			hasHarvest, plantType := g.Data.Environment.HarvestPlant(tileX, tileY)
			if hasHarvest {
				player.AddToBackpack(utils.PlantItemMapping[plantType], 1)
			}
		}
	}
}

func checkKeyboardOnPlayState(g *Game) {
	player := g.Data.Players[g.PlayerID]
	// player movement
	if ebiten.IsKeyPressed(ebiten.KeyA) && player.Sprite.X > 0 {
		player.Direction = utils.Left
		player.State = utils.WalkState
		player.Dx -= utils.MovementSpeed
		if !playerHasCollisions(g) {
			player.UpdateLocation()
		} else {
			player.Dx = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyD) &&
		player.Sprite.X < utils.MapWidth-player.Sprite.Width {
		player.Direction = utils.Right
		player.State = utils.WalkState
		player.Dx += utils.MovementSpeed
		if !playerHasCollisions(g) {
			player.UpdateLocation()
		} else {
			player.Dx = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyW) && player.Sprite.Y > 0 {
		player.Direction = utils.Back
		player.State = utils.WalkState
		player.Dy -= utils.MovementSpeed
		if !playerHasCollisions(g) {
			player.UpdateLocation()
		} else {
			player.Dy = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyS) &&
		player.Sprite.Y < utils.MapHeight-player.Sprite.Height {
		player.Direction = utils.Front
		player.State = utils.WalkState
		player.Dy += utils.MovementSpeed
		if !playerHasCollisions(g) {
			player.UpdateLocation()
		} else {
			player.Dy = 0
		}
	} else if player.StateTTL == 0 {
		player.State = utils.IdleState
	}

	// equip item
	if ebiten.IsKeyPressed(ebiten.Key1) {
		player.EquippedItem = 0
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		player.EquippedItem = 1
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		player.EquippedItem = 2
	} else if ebiten.IsKeyPressed(ebiten.Key4) {
		player.EquippedItem = 3
	} else if ebiten.IsKeyPressed(ebiten.Key5) {
		player.EquippedItem = 4
	} else if ebiten.IsKeyPressed(ebiten.Key6) {
		player.EquippedItem = 5
	} else if ebiten.IsKeyPressed(ebiten.Key7) {
		player.EquippedItem = 6
	} else if ebiten.IsKeyPressed(ebiten.Key8) {
		player.EquippedItem = 7
	} else if ebiten.IsKeyPressed(ebiten.Key9) {
		player.EquippedItem = 8
	}
}

func checkMouseOnCraftState(g *Game) {
	player := g.Data.Players[g.PlayerID]
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
			if player.RemoveFromBackpack(items) {
				g.Sounds.PlaySound(g.Sounds.SFXCraft)
				player.AddToBackpack(utils.Recipes[g.UIState.SelectedRecipe], recipe.Count)
			} else {
				g.SetErrorMessage("Not enough materials!")
			}
		}
	}
}

func checkMouseOnCustomCharState(g *Game) {
	player := g.Data.Players[g.PlayerID]
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
			player.Spritesheet = g.UIState.SelectedCharacter
			g.State = utils.GameStatePlay
		}
	}
}
