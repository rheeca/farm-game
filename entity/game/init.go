package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"guion-2d-project3/entity/animal"
	"guion-2d-project3/entity/environment"
	"guion-2d-project3/entity/loader"
	"guion-2d-project3/entity/model"
	"guion-2d-project3/entity/player"
	"guion-2d-project3/utils"
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
	g.Player.UpdateFrame(g.CurrentFrame)
	getPlayerInput(g)
	if g.State == utils.GameStatePlay {
		g.CurrentFrame += 1
		if g.CurrentMap == utils.AnimalsMap {
			updateAnimals(g)
		}
		updateTrees(g)
		updateObjects(g)
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
	drawFarmPlots(g, screen, drawOptions)
	if g.CurrentMap == utils.ForestMap {
		drawTrees(g, screen, drawOptions)
	}
	drawObjects(g, screen, drawOptions)

	if g.CurrentMap == utils.AnimalsMap {
		drawChickens(g, screen, drawOptions)
		drawCows(g, screen, drawOptions)
	}

	drawPlayer(g, screen, drawOptions)
	drawBackpack(g, screen, drawOptions)

	if g.State == utils.GameStateCraft {
		drawCraftingUI(g, screen, drawOptions)
	}

	drawImageToShow(g, screen, drawOptions)
	drawErrorMessage(g, screen, drawOptions)
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
