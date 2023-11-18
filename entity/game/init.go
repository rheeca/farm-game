package game

import (
	"github.com/hajimehoshi/ebiten/v2"
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
	Animals      []*animal.Animal
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

	// draw animals
	for _, a := range g.Animals {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(a.XLoc), float64(a.YLoc))
		screen.DrawImage(a.Spritesheet.SubImage(image.Rect(a.Frame*a.Width,
			a.Direction*a.Height,
			a.Frame*a.Width+a.Width,
			a.Direction*a.Height+a.Height)).(*ebiten.Image), &drawOptions)
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
		if item == 0 {
			continue
		}
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(utils.ToolsFirstSlotX+(i*utils.ToolsUIBoxSize)), float64(utils.ToolsFirstSlotY))

		screen.DrawImage(g.Images.FarmItems.SubImage(image.Rect((item%utils.FarmItemsColumns)*utils.UnitSize,
			(item/utils.FarmItemsColumns)*utils.UnitSize,
			(item%utils.FarmItemsColumns)*utils.UnitSize+utils.UnitSize,
			(item/utils.FarmItemsColumns)*utils.UnitSize+utils.UnitSize)).(*ebiten.Image), &drawOptions)
	}
}

func (g *Game) Layout(oWidth, oHeight int) (sWidth, sHeight int) {
	return oWidth, oHeight
}
