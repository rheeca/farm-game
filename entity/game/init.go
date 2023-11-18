package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
	"guion-2d-project3/entity/animal"
	"guion-2d-project3/entity/environment"
	"guion-2d-project3/entity/loader"
	"guion-2d-project3/entity/player"
	"guion-2d-project3/utils"
	"image"
)

type Game struct {
	Environment  *environment.Environment
	Player       *player.Player
	Chickens     []*animal.Chicken
	CurrentMap   int
	CurrentFrame int
	Images       loader.ImageCollection
}

func (g *Game) Update() error {
	g.CurrentFrame += 1
	getPlayerInput(g)
	updateAnimals(g)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawOptions := ebiten.DrawImageOptions{}

	drawMap(g, screen, drawOptions, g.CurrentMap)

	// draw chickens
	for _, c := range g.Chickens {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(c.XLoc), float64(c.YLoc))
		screen.DrawImage(c.Spritesheet.SubImage(image.Rect(c.Frame*c.SpriteWidth,
			(c.State*utils.AnimalNumOfDirections+c.Direction)*c.SpriteHeight,
			c.Frame*c.SpriteWidth+c.SpriteWidth,
			(c.State*utils.NumOfDirections+c.Direction)*c.SpriteHeight+c.SpriteHeight)).(*ebiten.Image), &drawOptions)
	}

	// draw player
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(g.Player.XLoc), float64(g.Player.YLoc))
	screen.DrawImage(g.Player.Spritesheet.SubImage(image.Rect(g.Player.Frame*g.Player.SpriteWidth,
		(g.Player.State*utils.NumOfDirections+g.Player.Direction)*g.Player.SpriteHeight,
		g.Player.Frame*g.Player.SpriteWidth+g.Player.SpriteWidth,
		(g.Player.State*utils.NumOfDirections+g.Player.Direction)*g.Player.SpriteHeight+g.Player.SpriteHeight)).(*ebiten.Image), &drawOptions)

	// draw tools ui
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(utils.ToolsUIX), float64(utils.ToolsUIY))
	screen.DrawImage(g.Images.ToolsUI, &drawOptions)

	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(utils.ToolsFirstBoxX+((g.Player.EquippedItem)*utils.ToolsUIBoxSize)),
		float64(utils.ToolsFirstBoxY))
	screen.DrawImage(g.Images.SelectedTool, &drawOptions)

	for i, item := range g.Player.Backpack {
		if item.ID == 0 {
			continue
		}
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(utils.ToolsFirstSlotX+(i*utils.ToolsUIBoxSize)), float64(utils.ToolsFirstSlotY))

		screen.DrawImage(g.Images.FarmItems.SubImage(image.Rect((item.ID%utils.FarmItemsColumns)*utils.UnitSize,
			(item.ID/utils.FarmItemsColumns)*utils.UnitSize,
			(item.ID%utils.FarmItemsColumns)*utils.UnitSize+utils.UnitSize,
			(item.ID/utils.FarmItemsColumns)*utils.UnitSize+utils.UnitSize)).(*ebiten.Image), &drawOptions)

		// draw item count
		if item.Count > 1 {
			DrawCenteredText(screen, basicfont.Face7x13, fmt.Sprintf("%d", item.Count),
				utils.ToolsFirstSlotX+(i*utils.ToolsUIBoxSize)+utils.UnitSize, utils.ToolsFirstSlotY)
		}
	}
}

func (g *Game) Layout(oWidth, oHeight int) (sWidth, sHeight int) {
	return oWidth, oHeight
}
