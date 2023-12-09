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

	"github.com/gofrs/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
)

type Game struct {
	State        int
	Data         *GameData
	Maps         []*tiled.Map
	CurrentMap   int
	CurrentFrame int
	PlayerID     string
	Images       loader.ImageCollection
	Sounds       loader.SoundCollection
	UIState      model.UIState
}

type GameData struct {
	Environment *environment.Environment
	Players     map[string]*player.Player
	Chickens    []*animal.Chicken
	Cows        []*animal.Cow
}

func NewGame(embeddedAssets embed.FS) Game {
	gameMaps := LoadMaps(embeddedAssets, path.Join("client", "assets"))
	currentMap := utils.FarmMap
	windowWidth := gameMaps[currentMap].Width * gameMaps[currentMap].TileWidth
	windowHeight := gameMaps[currentMap].Height * gameMaps[currentMap].TileHeight
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(utils.ProjectTitle)

	images := loader.NewImageCollection(embeddedAssets, path.Join("client", "assets"))
	SetConstants(gameMaps[currentMap], images)

	// load environment
	env := environment.NewEnvironment(embeddedAssets, gameMaps)

	// load audio
	sounds := loader.NewSoundCollection(embeddedAssets, path.Join("client", "assets"))

	// load player
	players := map[string]*player.Player{}
	playerID := uuid.Must(uuid.NewV4()).String()
	spawnPoint := gameMaps[currentMap].Groups[0].ObjectGroups[utils.FarmMapSpawnPoint].Objects[0]
	playerChar := player.NewPlayer(playerID, int(spawnPoint.X), int(spawnPoint.Y), images)
	players[playerChar.PlayerID] = playerChar

	// load chickens
	embeddedFile, err := embeddedAssets.Open(path.Join("client", "assets", "animals", utils.ChickenImg))
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
	embeddedFile, err = embeddedAssets.Open(path.Join("client", "assets", "animals", utils.CowImg))
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
		State: utils.GameStateCustomChar,
		Data: &GameData{
			Environment: env,
			Players:     players,
			Chickens:    chickens,
			Cows:        cows,
		},
		Maps:       gameMaps,
		CurrentMap: currentMap,
		PlayerID:   playerChar.PlayerID,
		Images:     images,
		Sounds:     sounds,
	}
}

func (g *Game) Update() error {
	g.Data.Players[g.PlayerID].UpdateFrame(g.CurrentFrame)
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

	DrawMap(g.Maps[g.CurrentMap], g.Images.Tilesets, screen, drawOptions)
	drawFarmPlots(g, screen, drawOptions)
	if g.CurrentMap == utils.ForestMap {
		drawTrees(g, screen, drawOptions)
	}
	drawObjects(g, screen, drawOptions)

	if g.CurrentMap == utils.AnimalsMap {
		drawChickens(g, screen, drawOptions)
		drawCows(g, screen, drawOptions)
	}

	drawPlayers(g, screen, drawOptions)
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

func SetConstants(gameMap *tiled.Map, images loader.ImageCollection) {
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

func LoadMaps(embeddedAssets embed.FS, assetPath string) (gameMaps []*tiled.Map) {
	farmMap, err := utils.LoadMapFromEmbedded(embeddedAssets,
		path.Join(assetPath, utils.FarmMapFile))
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	animalsMap, err := utils.LoadMapFromEmbedded(embeddedAssets,
		path.Join(assetPath, utils.AnimalsMapFile))
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	forestMap, err := utils.LoadMapFromEmbedded(embeddedAssets,
		path.Join(assetPath, utils.ForestMapFile))
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}

	gameMaps = append(gameMaps, farmMap, animalsMap, forestMap)
	return gameMaps
}
