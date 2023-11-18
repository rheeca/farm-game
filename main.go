package main

import (
	"embed"
	"fmt"
	"guion-2d-project3/entity/animal"
	"guion-2d-project3/entity/environment"
	"guion-2d-project3/entity/game"
	"guion-2d-project3/entity/loader"
	"guion-2d-project3/entity/player"
	"guion-2d-project3/utils"
	"log"
	"os"
	"path"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
)

//go:embed assets/*
var EmbeddedAssets embed.FS

func loadMapFromEmbedded(name string) (*tiled.Map, error) {
	embeddedMap, err := tiled.LoadFile(name,
		tiled.WithFileSystem(EmbeddedAssets))
	if err != nil {
		return nil, err
	}
	return embeddedMap, nil
}

func loadWavFromEmbedded(name string, context *audio.Context) (soundPlayer *audio.Player, err error) {
	soundFile, err := EmbeddedAssets.Open(path.Join("assets", "sounds", name))
	if err != nil {
		log.Fatal("failed to load embedded audio:", soundFile, err)
	}
	sound, err := wav.DecodeWithoutResampling(soundFile)
	if err != nil {
		return soundPlayer, fmt.Errorf("failed to interpret sound file: %s", err)
	}
	soundPlayer, err = context.NewPlayer(sound)
	if err != nil {
		return soundPlayer, fmt.Errorf("failed to create sound player: %s", err)
	}
	return soundPlayer, nil
}

func main() {
	gameMap, err := loadMapFromEmbedded(path.Join("assets", utils.MapFile))
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	windowWidth := gameMap.Width * gameMap.TileWidth
	windowHeight := gameMap.Height * gameMap.TileHeight
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(utils.ProjectTitle)

	images, err := loader.NewImageCollection(EmbeddedAssets)
	if err != nil {
		log.Fatal(err.Error())
	}
	setConstants(gameMap, images)

	// load environment
	env := environment.NewEnvironment(EmbeddedAssets, []*tiled.Map{gameMap})

	// load audio
	audioContext := audio.NewContext(utils.SoundSampleRate)
	bgmPlayer, err := loadWavFromEmbedded(utils.FirstTownAudio, audioContext)
	if err != nil {
		log.Fatal(err.Error())
	}

	// load player
	embeddedFile, err := EmbeddedAssets.Open(path.Join("assets", utils.PlayerImg))
	if err != nil {
		log.Fatal("failed to load embedded image:", embeddedFile, err)
	}
	playerImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("error loading player image")
	}
	playerChar := player.NewPlayer(playerImage)

	// load chickens
	embeddedFile, err = EmbeddedAssets.Open(path.Join("assets", "animals", utils.ChickenImg))
	if err != nil {
		log.Fatal("failed to load embedded image:", embeddedFile, err)
	}
	chickenImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("error loading chicken image")
	}
	chicken := animal.NewChicken(chickenImage)

	gameObj := game.Game{
		Environment: env,
		CurrentMap:  utils.FarmMap,
		Player:      playerChar,
		Chickens:    []*animal.Chicken{chicken},
		Images:      images,
	}

	go func(player *audio.Player) {
		player.Play()
		time.Sleep(122 * time.Second)
		player.Rewind()
	}(bgmPlayer)
	err = ebiten.RunGame(&gameObj)
	if err != nil {
		fmt.Println("failed to run game:", err)
	}
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
	utils.FarmItemsColumns = images.FarmItems.Bounds().Dx() / utils.UnitSize
}
