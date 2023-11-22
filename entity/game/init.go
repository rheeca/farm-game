package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"guion-2d-project3/entity/animal"
	"guion-2d-project3/entity/environment"
	"guion-2d-project3/entity/loader"
	"guion-2d-project3/entity/model"
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
	Sounds       loader.SoundCollection
	UIState      model.UIState
}

func (g *Game) Update() error {
	getPlayerInput(g)
	if g.State == utils.GameStateCustomChar || g.State == utils.GameStateCraft {
		return nil
	}

	g.CurrentFrame += 1
	updateAnimals(g)
	for i, t := range g.Environment.Trees {
		if t.IsNil {
			continue
		}
		g.Environment.Trees[i].UpdateFrame(g.CurrentFrame)
	}
	for i, o := range g.Environment.Objects[g.CurrentMap] {
		if o.IsNil {
			continue
		}
		if o.Type == utils.ItemDoor {
			g.Environment.Objects[g.CurrentMap][i].UpdateFrame(g.CurrentFrame)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawOptions := ebiten.DrawImageOptions{}
	if g.State == utils.GameStateCustomChar {
		drawCharacterCustomizationUI(g, screen, drawOptions)
		return
	}

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
		drawCraftingUI(g, screen, drawOptions)
	}

	// draw image to show, if any
	if g.UIState.ImageTTL > 0 {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(0, 0)
		screen.DrawImage(g.Images.BlackScreen, &drawOptions)
		g.UIState.ImageTTL -= 1
		if g.UIState.ImageTTL == 0 {
			g.UIState.ImageToShow = nil
		}
	}

	// draw error message, if any
	if g.UIState.ErrorMessageTTL > 0 {
		text.Draw(screen, g.UIState.ErrorMessage, utils.LoadFont(12), 12, 20, colornames.Brown)
		g.UIState.ErrorMessageTTL -= 1
		if g.UIState.ErrorMessageTTL == 0 {
			g.UIState.ErrorMessage = ""
		}
	}
}

func (g *Game) Layout(oWidth, oHeight int) (sWidth, sHeight int) {
	return oWidth, oHeight
}

func (g *Game) SetErrorMessage(message string) {
	g.UIState.ErrorMessage = message
	g.UIState.ErrorMessageTTL = 60
}

func (g *Game) ShowImage(image *ebiten.Image) {
	g.UIState.ImageToShow = image
	g.UIState.ImageTTL = 60
}
