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
	State        int
	Environment  *environment.Environment
	Player       *player.Player
	Chickens     []*animal.Chicken
	Cows         []*animal.Cow
	CurrentMap   int
	CurrentFrame int
	Images       loader.ImageCollection
}

func (g *Game) Update() error {
	if g.State == utils.GameStateCraft {
		return nil
	}

	g.CurrentFrame += 1
	getPlayerInput(g)
	updateAnimals(g)
	for i := range g.Environment.Trees {
		g.Environment.Trees[i].UpdateFrame(g.CurrentFrame)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawOptions := ebiten.DrawImageOptions{}

	drawMap(g, screen, drawOptions)
	if g.CurrentMap == utils.ForestMap {
		drawTrees(g, screen, drawOptions)
	}
	drawObjects(g, screen, drawOptions)

	// draw chickens
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

	// draw cows
	for _, c := range g.Cows {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(c.XLoc), float64(c.YLoc))
		screen.DrawImage(c.Spritesheet.SubImage(image.Rect(c.Frame*utils.CowSpriteWidth,
			(c.State*utils.AnimalNumOfDirections+c.Direction)*utils.CowSpriteHeight,
			c.Frame*utils.CowSpriteWidth+utils.CowSpriteWidth,
			(c.State*utils.AnimalNumOfDirections+c.Direction)*utils.CowSpriteHeight+utils.CowSpriteHeight)).(*ebiten.Image), &drawOptions)
	}

	// draw player
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(g.Player.XLoc), float64(g.Player.YLoc))
	screen.DrawImage(g.Player.Spritesheet.SubImage(image.Rect(g.Player.Frame*utils.PlayerSpriteWidth,
		(g.Player.State*utils.NumOfDirections+g.Player.Direction)*utils.PlayerSpriteHeight,
		g.Player.Frame*utils.PlayerSpriteWidth+utils.PlayerSpriteWidth,
		(g.Player.State*utils.NumOfDirections+g.Player.Direction)*utils.PlayerSpriteHeight+utils.PlayerSpriteHeight)).(*ebiten.Image), &drawOptions)

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

	if g.State == utils.GameStateCraft {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(0, 0)
		screen.DrawImage(g.Images.CraftingUI, &drawOptions)
	}
}

func (g *Game) Layout(oWidth, oHeight int) (sWidth, sHeight int) {
	return oWidth, oHeight
}
