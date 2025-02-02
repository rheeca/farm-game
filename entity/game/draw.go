package game

import (
	"fmt"
	"guion-2d-project3/entity/animal"
	"guion-2d-project3/entity/loader"
	"guion-2d-project3/entity/model"
	"guion-2d-project3/entity/player"
	"guion-2d-project3/utils"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/lafriks/go-tiled"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

func DrawMap(gMap *tiled.Map, tilesets map[string]*ebiten.Image, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for i, layer := range gMap.Layers {
		if layer.Name == utils.GuideOnlyLayer || i == utils.CollisionLayer || i == utils.FarmingLandLayer {
			continue
		}
		for tileY := 0; tileY < utils.MapRows; tileY += 1 {
			for tileX := 0; tileX < utils.MapColumns; tileX += 1 {
				// find img of tile to draw
				tileToDraw := layer.Tiles[tileY*utils.MapColumns+tileX]
				if tileToDraw.IsNil() {
					continue
				}

				tileToDrawX := int(tileToDraw.ID) % tileToDraw.Tileset.Columns
				tileToDrawY := int(tileToDraw.ID) / tileToDraw.Tileset.Columns

				ebitenTileToDraw := tilesets[tileToDraw.Tileset.Name].SubImage(image.Rect(tileToDrawX*utils.TileWidth,
					tileToDrawY*utils.TileHeight,
					tileToDrawX*utils.TileWidth+utils.TileWidth,
					tileToDrawY*utils.TileHeight+utils.TileHeight)).(*ebiten.Image)

				// draw tile
				drawOptions.GeoM.Reset()
				TileXpos := float64(utils.TileWidth * tileX)
				TileYpos := float64(utils.TileHeight * tileY)
				drawOptions.GeoM.Translate(TileXpos, TileYpos)
				screen.DrawImage(ebitenTileToDraw, &drawOptions)
			}
		}
	}
}

func DrawTrees(trees []model.Object, images loader.ImageCollection, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for _, t := range trees {
		if t.IsNil {
			continue
		}
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(t.XLoc), float64(t.YLoc))
		screen.DrawImage(images.TreeSprites.SubImage(image.Rect(t.Frame*t.Sprite.Width,
			0,
			t.Frame*t.Sprite.Width+t.Sprite.Width,
			t.Sprite.Height)).(*ebiten.Image), &drawOptions)
	}
}

func DrawObjects(objects []model.Object, images loader.ImageCollection, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	tilesets := images.Tilesets
	for _, o := range objects {
		if o.IsNil {
			continue
		}
		var x0, y0, x1, y1 int
		var objImage *ebiten.Image

		drawOptions.GeoM.Reset()
		if o.Type == utils.ItemCraftingTable {
			objImage = images.CraftingTable
			x0, y0 = 0, 0
			x1 = objImage.Bounds().Dx()
			y1 = objImage.Bounds().Dy()
		} else if o.Type == utils.ItemDoor {
			var animation int
			if o.IsCollision {
				animation = 0
			} else {
				animation = 1
			}
			objImage = images.DoorSprites
			x0 = o.Frame * o.Sprite.Width
			y0 = animation * o.Sprite.Height
			x1 = o.Frame*o.Sprite.Width + o.Sprite.Width
			y1 = animation*o.Sprite.Height + o.Sprite.Height
		} else if o.Type == utils.ItemBedPink {
			objImage = images.BedPink
			x0, y0 = 0, 0
			x1 = objImage.Bounds().Dx()
			y1 = objImage.Bounds().Dy()
		} else if o.Type == utils.ItemMapStone3 {
			objImage = tilesets[utils.TilesetFlowersStones]
			x0 = utils.UnitSize * (utils.MapStone3 % 12)
			y0 = utils.UnitSize * (utils.MapStone3 / 12)
			x1 = x0 + utils.UnitSize
			y1 = y0 + utils.UnitSize
		} else if o.Type == utils.ItemMapWood {
			objImage = tilesets[utils.TilesetTrees]
			x0 = utils.UnitSize * (utils.MapWood % 12)
			y0 = utils.UnitSize * (utils.MapWood / 12)
			x1 = x0 + utils.UnitSize
			y1 = y0 + utils.UnitSize
		} else if o.Type == utils.MapSunflower {
			objImage = tilesets[utils.TilesetFlowersStones]
			x0 = utils.UnitSize * (utils.MapSunflower % 12)
			y0 = utils.UnitSize * (utils.MapSunflower / 12)
			x1 = x0 + utils.UnitSize
			y1 = y0 + utils.UnitSize
		} else if o.Type == utils.MapBlueflower {
			objImage = tilesets[utils.TilesetFlowersStones]
			x0 = utils.UnitSize * (utils.MapBlueflower % 12)
			y0 = utils.UnitSize * (utils.MapBlueflower / 12)
			x1 = x0 + utils.UnitSize
			y1 = y0 + utils.UnitSize
		} else if o.Type == utils.MapWeed {
			objImage = tilesets[utils.TilesetFlowersStones]
			x0 = utils.UnitSize * (utils.MapWeed % 12)
			y0 = utils.UnitSize * (utils.MapWeed / 12)
			x1 = x0 + utils.UnitSize
			y1 = y0 + utils.UnitSize
		} else if o.Type == utils.MapPinkDyeFlower {
			objImage = tilesets[utils.TilesetFlowersStones]
			x0 = utils.UnitSize * (utils.MapPinkDyeFlower % 12)
			y0 = utils.UnitSize * (utils.MapPinkDyeFlower / 12)
			x1 = x0 + utils.UnitSize
			y1 = y0 + utils.UnitSize
		} else if o.Type == utils.MapBlueDyeFlower {
			objImage = tilesets[utils.TilesetFlowersStones]
			x0 = utils.UnitSize * (utils.MapBlueDyeFlower % 12)
			y0 = utils.UnitSize * (utils.MapBlueDyeFlower / 12)
			x1 = x0 + utils.UnitSize
			y1 = y0 + utils.UnitSize
		}

		drawOptions.GeoM.Translate(float64(o.XLoc), float64(o.YLoc))
		screen.DrawImage(objImage.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image), &drawOptions)
	}
}

func DrawCraftingUI(player *player.Player, images loader.ImageCollection, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	// draw box
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(0, 0)
	screen.DrawImage(images.CraftingUI, &drawOptions)

	// draw selected tile
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(utils.CraftingUIFirstBoxX+(utils.CraftingUISpacing*(player.UIState.SelectedRecipe%utils.CraftingUIColumns))),
		float64(utils.CraftingUIFirstBoxY+(utils.CraftingUISpacing*(player.UIState.SelectedRecipe/utils.CraftingUIColumns))))
	screen.DrawImage(images.SelectedItem, &drawOptions)

	// draw recipes
	var x int
	for i, itemID := range utils.Recipes {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(utils.CraftingUIFirstSlotX+(x%(utils.CraftingUIColumns*utils.CraftingUISpacing))),
			float64(utils.CraftingUIFirstSlotY+(utils.CraftingUISpacing*(i/utils.CraftingUIColumns))))
		screen.DrawImage(images.FarmItems.SubImage(image.Rect((itemID%utils.FarmItemsColumns)*utils.UnitSize,
			(itemID/utils.FarmItemsColumns)*utils.UnitSize,
			(itemID%utils.FarmItemsColumns)*utils.UnitSize+utils.UnitSize,
			(itemID/utils.FarmItemsColumns)*utils.UnitSize+utils.UnitSize)).(*ebiten.Image), &drawOptions)
		x += utils.CraftingUISpacing
	}

	// draw recipe ingredients
	for i, item := range utils.RecipeDetails[utils.Recipes[player.UIState.SelectedRecipe]].Materials {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(utils.RecipeItemX+(i*64)),
			float64(utils.RecipeItemY))
		screen.DrawImage(images.FarmItems.SubImage(image.Rect((item.ID%utils.FarmItemsColumns)*utils.UnitSize,
			(item.ID/utils.FarmItemsColumns)*utils.UnitSize,
			(item.ID%utils.FarmItemsColumns)*utils.UnitSize+utils.UnitSize,
			(item.ID/utils.FarmItemsColumns)*utils.UnitSize+utils.UnitSize)).(*ebiten.Image), &drawOptions)
		DrawCenteredText(screen, basicfont.Face7x13, fmt.Sprintf("%v", item.Count), utils.RecipeItemX+(i*64)+32, utils.RecipeItemY)
	}
}

func DrawCharacterCustomizationUI(player *player.Player, images loader.ImageCollection, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	// draw box
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(0, 0)
	screen.DrawImage(images.CharacterCustomizationUI, &drawOptions)

	// draw selected tile
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(utils.CharacterUIFirstBoxX+(utils.CharacterUISpacing*(player.UIState.SelectedCharacter%utils.CharacterUIColumns))),
		float64(utils.CharacterUIFirstBoxY+(utils.CharacterUISpacing*(player.UIState.SelectedCharacter/utils.CharacterUIColumns))))
	screen.DrawImage(images.SelectedCharacter, &drawOptions)

	// draw characters
	for i, img := range images.Characters {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(utils.CharacterUIFirstSlotX+(utils.CharacterUISpacing*(i%utils.CharacterUIColumns))),
			float64(utils.CharacterUIFirstSlotY+(utils.CharacterUISpacing*(i/utils.CharacterUIColumns))))
		screen.DrawImage(img.SubImage(image.Rect(utils.UnitSize,
			utils.UnitSize,
			utils.UnitSize*2,
			utils.UnitSize*2)).(*ebiten.Image), &drawOptions)
	}
}

func DrawPlayers(currentMap int, players map[string]*player.Player, images loader.ImageCollection, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for _, player := range players {
		if currentMap != player.CurrentMap {
			continue
		}
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(player.XLoc), float64(player.YLoc))
		screen.DrawImage(images.Characters[player.Spritesheet].SubImage(image.Rect(player.Frame*utils.PlayerSpriteWidth,
			(player.State*utils.NumOfDirections+player.Direction)*utils.PlayerSpriteHeight,
			player.Frame*utils.PlayerSpriteWidth+utils.PlayerSpriteWidth,
			(player.State*utils.NumOfDirections+player.Direction)*utils.PlayerSpriteHeight+utils.PlayerSpriteHeight)).(*ebiten.Image), &drawOptions)
	}
}

func DrawChickens(chickens []*animal.Chicken, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for _, c := range chickens {
		drawOptions.GeoM.Reset()

		var spriteHeight, yLoc int
		if c.State == utils.ChickenHeartState {
			spriteHeight = utils.UnitSize * 2
			yLoc = c.YLoc - utils.UnitSize
		} else {
			spriteHeight = c.Sprite.Height
			yLoc = c.YLoc
		}
		drawOptions.GeoM.Translate(float64(c.XLoc), float64(yLoc))
		screen.DrawImage(c.Spritesheet.SubImage(image.Rect(c.Frame*c.Sprite.Width,
			(c.State*utils.AnimalNumOfDirections+c.Direction)*c.Sprite.Height,
			c.Frame*c.Sprite.Width+c.Sprite.Width,
			(c.State*utils.AnimalNumOfDirections+c.Direction)*c.Sprite.Height+spriteHeight)).(*ebiten.Image), &drawOptions)
	}
}

func DrawCows(cows []*animal.Cow, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for _, c := range cows {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(c.XLoc), float64(c.YLoc))
		screen.DrawImage(c.Spritesheet.SubImage(image.Rect(c.Frame*utils.CowSpriteWidth,
			(c.State*utils.AnimalNumOfDirections+c.Direction)*utils.CowSpriteHeight,
			c.Frame*utils.CowSpriteWidth+utils.CowSpriteWidth,
			(c.State*utils.AnimalNumOfDirections+c.Direction)*utils.CowSpriteHeight+utils.CowSpriteHeight)).(*ebiten.Image), &drawOptions)
	}
}

func DrawFarmPlots(plots []model.Plot, images loader.ImageCollection, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for _, p := range plots {
		// draw soil
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(p.XTile*utils.TileWidth), float64(p.YTile*utils.TileHeight))

		var tileID int
		var tileset string
		if p.IsWatered {
			tileID = 12
			tileset = utils.TilesetDarkerSoilGround
		} else {
			tileID = 12
			tileset = utils.TilesetSoilGround
		}
		tileToDrawX := tileID % 11
		tileToDrawY := tileID / 11
		screen.DrawImage(images.Tilesets[tileset].SubImage(image.Rect(tileToDrawX*utils.TileWidth,
			tileToDrawY*utils.TileHeight,
			tileToDrawX*utils.TileWidth+utils.TileWidth,
			tileToDrawY*utils.TileHeight+utils.TileHeight)).(*ebiten.Image), &drawOptions)

		// draw plant
		if p.HasPlant && !p.ReadyForHarvest {
			drawOptions.GeoM.Reset()
			drawOptions.GeoM.Translate(float64(p.XTile*utils.TileWidth), float64(p.YTile*utils.TileHeight))
			screen.DrawImage(images.FarmingPlants.SubImage(image.Rect(0,
				utils.UnitSize*utils.PlantTomato,
				utils.UnitSize,
				(utils.UnitSize*utils.PlantTomato)+utils.UnitSize)).(*ebiten.Image), &drawOptions)
		} else if p.ReadyForHarvest {
			drawOptions.GeoM.Reset()
			drawOptions.GeoM.Translate(float64(p.XTile*utils.TileWidth), float64(p.YTile*utils.TileHeight))
			screen.DrawImage(images.FarmingPlants.SubImage(image.Rect(96,
				utils.UnitSize*utils.PlantTomato,
				96+utils.UnitSize,
				(utils.UnitSize*utils.PlantTomato)+utils.UnitSize)).(*ebiten.Image), &drawOptions)
		}
	}
}

func DrawBackpack(player *player.Player, images loader.ImageCollection, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	// draw backpack
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(utils.ToolsUIX), float64(utils.ToolsUIY))
	screen.DrawImage(images.ToolsUI, &drawOptions)

	// draw selected box
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(utils.ToolsFirstBoxX+((player.EquippedItem)*utils.BackpackUIBoxWidth)),
		float64(utils.ToolsFirstBoxY))
	screen.DrawImage(images.SelectedTool, &drawOptions)

	// draw items in backpack
	for i, item := range player.Backpack {
		if item.ID == 0 {
			continue
		}
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(utils.ToolsFirstSlotX+(i*utils.BackpackUIBoxWidth)), float64(utils.ToolsFirstSlotY))

		screen.DrawImage(images.FarmItems.SubImage(image.Rect((item.ID%utils.FarmItemsColumns)*utils.UnitSize,
			(item.ID/utils.FarmItemsColumns)*utils.UnitSize,
			(item.ID%utils.FarmItemsColumns)*utils.UnitSize+utils.UnitSize,
			(item.ID/utils.FarmItemsColumns)*utils.UnitSize+utils.UnitSize)).(*ebiten.Image), &drawOptions)

		// draw item count
		if item.Count > 1 {
			DrawCenteredText(screen, basicfont.Face7x13, fmt.Sprintf("%d", item.Count),
				utils.ToolsFirstSlotX+(i*utils.BackpackUIBoxWidth)+utils.UnitSize, utils.ToolsFirstSlotY)
		}
	}

	// draw delete button
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(utils.BackpackDeleteButtonX), float64(utils.BackpackDeleteButtonY))
	screen.DrawImage(images.ButtonDelete, &drawOptions)
}

func DrawImageToShow(p *player.Player, images loader.ImageCollection, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	if p.UIState.ImageTTL > 0 {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(0, 0)
		screen.DrawImage(images.BlackScreen, &drawOptions)
		p.UIState.ImageTTL -= 1
		if p.UIState.ImageTTL == 0 {
			p.UIState.ImageToShow = nil
		}
	}
}

func DrawErrorMessage(p *player.Player, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	if p.UIState.ErrorMessageTTL > 0 {
		text.Draw(screen, p.UIState.ErrorMessage, utils.LoadFont(12), 12, 20, colornames.Brown)
		p.UIState.ErrorMessageTTL -= 1
		if p.UIState.ErrorMessageTTL == 0 {
			p.UIState.ErrorMessage = ""
		}
	}
}

// DrawCenteredText %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
// from https://github.com/sedyh/ebitengine-cheatsheet
func DrawCenteredText(screen *ebiten.Image, font font.Face, s string, cx, cy int) {
	bounds := text.BoundString(font, s)
	x, y := cx-bounds.Min.X-bounds.Dx()/2, cy-bounds.Min.Y-bounds.Dy()/2
	text.Draw(screen, s, font, x, y, colornames.Brown)
}

// %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
