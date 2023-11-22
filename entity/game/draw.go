package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"guion-2d-project3/utils"
	"image"
)

func drawMap(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for _, layer := range g.Environment.Maps[g.CurrentMap].Layers {
		if layer.Name == utils.GuideOnlyLayer {
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

				ebitenTileToDraw := g.Environment.Tilesets[tileToDraw.Tileset.Name].SubImage(image.Rect(tileToDrawX*utils.TileWidth,
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

func drawTrees(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for _, t := range g.Environment.Trees {
		if t.IsNil {
			continue
		}
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(t.XLoc), float64(t.YLoc))
		screen.DrawImage(g.Images.TreeSprites.SubImage(image.Rect(t.Frame*t.Sprite.Width,
			0,
			t.Frame*t.Sprite.Width+t.Sprite.Width,
			t.Sprite.Height)).(*ebiten.Image), &drawOptions)
	}
}

func drawObjects(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for _, o := range g.Environment.Objects[g.CurrentMap] {
		if o.IsNil {
			continue
		}
		var x0, y0, x1, y1 int
		var objImage *ebiten.Image

		drawOptions.GeoM.Reset()
		if o.Type == utils.ItemCraftingTable {
			objImage = g.Images.CraftingTable
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
			objImage = g.Images.DoorSprites
			x0 = o.Frame * o.Sprite.Width
			y0 = animation * o.Sprite.Height
			x1 = o.Frame*o.Sprite.Width + o.Sprite.Width
			y1 = animation*o.Sprite.Height + o.Sprite.Height
		} else if o.Type == utils.ItemBedPink {
			objImage = g.Images.BedPink
			x0, y0 = 0, 0
			x1 = objImage.Bounds().Dx()
			y1 = objImage.Bounds().Dy()
		} else if o.Type == utils.ItemMapStone3 {
			objImage = g.Environment.Tilesets[utils.TilesetFlowersStones]
			x0 = utils.UnitSize * (utils.MapStone3 % 12)
			y0 = utils.UnitSize * (utils.MapStone3 / 12)
			x1 = x0 + utils.UnitSize
			y1 = y0 + utils.UnitSize
		} else if o.Type == utils.ItemMapWood {
			objImage = g.Environment.Tilesets[utils.TilesetTrees]
			x0 = utils.UnitSize * (utils.MapWood % 12)
			y0 = utils.UnitSize * (utils.MapWood / 12)
			x1 = x0 + utils.UnitSize
			y1 = y0 + utils.UnitSize
		}

		drawOptions.GeoM.Translate(float64(o.XLoc), float64(o.YLoc))
		screen.DrawImage(objImage.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image), &drawOptions)
	}
}

func drawCraftingUI(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	// draw box
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(0, 0)
	screen.DrawImage(g.Images.CraftingUI, &drawOptions)

	// draw selected tile
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(utils.CraftingUIFirstBoxX+(utils.CraftingUISpacing*(g.UIState.SelectedRecipe%utils.CraftingUIColumns))),
		float64(utils.CraftingUIFirstBoxY+(utils.CraftingUISpacing*(g.UIState.SelectedRecipe/utils.CraftingUIColumns))))
	screen.DrawImage(g.Images.SelectedItem, &drawOptions)

	// draw recipes
	var x int
	for i, itemID := range utils.Recipes {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(utils.CraftingUIFirstSlotX+(x%(utils.CraftingUIColumns*utils.CraftingUISpacing))),
			float64(utils.CraftingUIFirstSlotY+(utils.CraftingUISpacing*(i/utils.CraftingUIColumns))))
		screen.DrawImage(g.Images.FarmItems.SubImage(image.Rect((itemID%utils.FarmItemsColumns)*utils.UnitSize,
			(itemID/utils.FarmItemsColumns)*utils.UnitSize,
			(itemID%utils.FarmItemsColumns)*utils.UnitSize+utils.UnitSize,
			(itemID/utils.FarmItemsColumns)*utils.UnitSize+utils.UnitSize)).(*ebiten.Image), &drawOptions)
		x += utils.CraftingUISpacing
	}

	// draw recipe ingredients
	for i, item := range utils.RecipeDetails[utils.Recipes[g.UIState.SelectedRecipe]] {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(utils.RecipeItemX+(i*64)),
			float64(utils.RecipeItemY))
		screen.DrawImage(g.Images.FarmItems.SubImage(image.Rect((item.ID%utils.FarmItemsColumns)*utils.UnitSize,
			(item.ID/utils.FarmItemsColumns)*utils.UnitSize,
			(item.ID%utils.FarmItemsColumns)*utils.UnitSize+utils.UnitSize,
			(item.ID/utils.FarmItemsColumns)*utils.UnitSize+utils.UnitSize)).(*ebiten.Image), &drawOptions)
		DrawCenteredText(screen, basicfont.Face7x13, fmt.Sprintf("%v", item.Count), utils.RecipeItemX+(i*64)+32, utils.RecipeItemY)
	}
}

func drawCharacterCustomizationUI(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	// draw box
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(0, 0)
	screen.DrawImage(g.Images.CharacterCustomizationUI, &drawOptions)

	// draw selected tile
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(utils.CharacterUIFirstBoxX+(utils.CharacterUISpacing*(g.UIState.SelectedCharacter%utils.CharacterUIColumns))),
		float64(utils.CharacterUIFirstBoxY+(utils.CharacterUISpacing*(g.UIState.SelectedCharacter/utils.CharacterUIColumns))))
	screen.DrawImage(g.Images.SelectedCharacter, &drawOptions)

	// draw characters
	for i, img := range g.Images.Characters {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(utils.CharacterUIFirstSlotX+(utils.CharacterUISpacing*(i%utils.CharacterUIColumns))),
			float64(utils.CharacterUIFirstSlotY+(utils.CharacterUISpacing*(i/utils.CharacterUIColumns))))
		screen.DrawImage(img.SubImage(image.Rect(utils.UnitSize,
			utils.UnitSize,
			utils.UnitSize*2,
			utils.UnitSize*2)).(*ebiten.Image), &drawOptions)
	}
}

func drawPlayer(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(g.Player.XLoc), float64(g.Player.YLoc))
	screen.DrawImage(g.Player.Spritesheet.SubImage(image.Rect(g.Player.Frame*utils.PlayerSpriteWidth,
		(g.Player.State*utils.NumOfDirections+g.Player.Direction)*utils.PlayerSpriteHeight,
		g.Player.Frame*utils.PlayerSpriteWidth+utils.PlayerSpriteWidth,
		(g.Player.State*utils.NumOfDirections+g.Player.Direction)*utils.PlayerSpriteHeight+utils.PlayerSpriteHeight)).(*ebiten.Image), &drawOptions)
}

func drawChickens(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for _, c := range g.Chickens {
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

func drawCows(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for _, c := range g.Cows {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(c.XLoc), float64(c.YLoc))
		screen.DrawImage(c.Spritesheet.SubImage(image.Rect(c.Frame*utils.CowSpriteWidth,
			(c.State*utils.AnimalNumOfDirections+c.Direction)*utils.CowSpriteHeight,
			c.Frame*utils.CowSpriteWidth+utils.CowSpriteWidth,
			(c.State*utils.AnimalNumOfDirections+c.Direction)*utils.CowSpriteHeight+utils.CowSpriteHeight)).(*ebiten.Image), &drawOptions)
	}
}

func drawBackpack(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	// draw backpack
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(utils.ToolsUIX), float64(utils.ToolsUIY))
	screen.DrawImage(g.Images.ToolsUI, &drawOptions)

	// draw selected box
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(utils.ToolsFirstBoxX+((g.Player.EquippedItem)*utils.BackpackUIBoxWidth)),
		float64(utils.ToolsFirstBoxY))
	screen.DrawImage(g.Images.SelectedTool, &drawOptions)

	// draw items in backpack
	for i, item := range g.Player.Backpack {
		if item.ID == 0 {
			continue
		}
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(utils.ToolsFirstSlotX+(i*utils.BackpackUIBoxWidth)), float64(utils.ToolsFirstSlotY))

		screen.DrawImage(g.Images.FarmItems.SubImage(image.Rect((item.ID%utils.FarmItemsColumns)*utils.UnitSize,
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
	screen.DrawImage(g.Images.ButtonDelete, &drawOptions)
}

func drawImageToShow(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	if g.UIState.ImageTTL > 0 {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(0, 0)
		screen.DrawImage(g.Images.BlackScreen, &drawOptions)
		g.UIState.ImageTTL -= 1
		if g.UIState.ImageTTL == 0 {
			g.UIState.ImageToShow = nil
		}
	}
}

func drawErrorMessage(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	if g.UIState.ErrorMessageTTL > 0 {
		text.Draw(screen, g.UIState.ErrorMessage, utils.LoadFont(12), 12, 20, colornames.Brown)
		g.UIState.ErrorMessageTTL -= 1
		if g.UIState.ErrorMessageTTL == 0 {
			g.UIState.ErrorMessage = ""
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
