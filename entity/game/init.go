package game

import (
	"embed"
	"fmt"
	"guion-2d-project3/entity/animal"
	"guion-2d-project3/entity/environment"
	"guion-2d-project3/entity/loader"
	"guion-2d-project3/entity/model"
	"guion-2d-project3/entity/player"
	"guion-2d-project3/utils"
	"log"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
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

func NewGame(embeddedAssets embed.FS) Game {
	gameMap, err := utils.LoadMapFromEmbedded(embeddedAssets, path.Join("assets", utils.FarmMapFile))
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	windowWidth := gameMap.Width * gameMap.TileWidth
	windowHeight := gameMap.Height * gameMap.TileHeight
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(utils.ProjectTitle)

	images := loader.NewImageCollection(embeddedAssets)
	setConstants(gameMap, images)

	// load environment
	env := environment.NewEnvironment(embeddedAssets, []*tiled.Map{gameMap})

	// load audio
	sounds := loader.NewSoundCollection(embeddedAssets)

	// load player
	embeddedFile, err := embeddedAssets.Open(path.Join("assets", "player", utils.DefaultPlayerImg))
	if err != nil {
		log.Fatal("failed to load embedded image:", embeddedFile, err)
	}
	playerImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("error loading player image")
	}
	spawnPoint := gameMap.Groups[0].ObjectGroups[utils.FarmMapSpawnPoint].Objects[0]
	playerChar := player.NewPlayer(playerImage, int(spawnPoint.X), int(spawnPoint.Y))

	// load chickens
	embeddedFile, err = embeddedAssets.Open(path.Join("assets", "animals", utils.ChickenImg))
	if err != nil {
		log.Fatal("failed to load embedded image:", embeddedFile, err)
	}
	chickenImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("error loading chicken image")
	}
	var chickens []*animal.Chicken
	for _, v := range utils.ChickenLocations {
		chicken := animal.NewChicken(chickenImage, v)
		chickens = append(chickens, chicken)
	}

	// load cows
	embeddedFile, err = embeddedAssets.Open(path.Join("assets", "animals", utils.CowImg))
	if err != nil {
		log.Fatal("failed to load embedded image:", embeddedFile, err)
	}
	cowImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("error loading cow image")
	}
	var cows []*animal.Cow
	for _, v := range utils.CowLocations {
		cow := animal.NewCow(cowImage, v.X, v.Y)
		cows = append(cows, cow)
	}
	return Game{
		State:       utils.GameStateCustomChar,
		Environment: env,
		CurrentMap:  utils.FarmMap,
		Player:      playerChar,
		Chickens:    chickens,
		Cows:        cows,
		Images:      images,
		Sounds:      sounds,
	}
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

func setConstants(gameMap *tiled.Map, images loader.ImageCollection) {
	utils.MapWidth = gameMap.Width * gameMap.TileWidth
	utils.MapHeight = gameMap.Height * gameMap.TileHeight
	utils.TileWidth = gameMap.TileWidth
	utils.TileHeight = gameMap.TileHeight
	utils.MapColumns = gameMap.Width
	utils.MapRows = gameMap.Height
	utils.ToolsUIX = (utils.MapWidth / 2) - images.ToolsUI.Bounds().Dx()/2
	utils.ToolsUIY = utils.MapHeight - 60
	utils.ToolsFirstSlotX = utils.ToolsUIX + 22
	utils.ToolsFirstSlotY = utils.ToolsUIY + 10
	utils.ToolsFirstBoxX = utils.ToolsUIX + 14
	utils.ToolsFirstBoxY = utils.ToolsUIY + 2
}
