package main

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
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
)

//go:embed assets/*
var EmbeddedAssets embed.FS

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
	gameMap, err := utils.LoadMapFromEmbedded(EmbeddedAssets, path.Join("assets", utils.FarmMapFile))
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	windowWidth := gameMap.Width * gameMap.TileWidth
	windowHeight := gameMap.Height * gameMap.TileHeight
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(utils.ProjectTitle)

	images := loader.NewImageCollection(EmbeddedAssets)
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
	embeddedFile, err := EmbeddedAssets.Open(path.Join("assets", "player", utils.DefaultPlayerImg))
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
	embeddedFile, err = EmbeddedAssets.Open(path.Join("assets", "animals", utils.ChickenImg))
	if err != nil {
		log.Fatal("failed to load embedded image:", embeddedFile, err)
	}
	chickenImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("error loading chicken image")
	}
	var chickens []*animal.Chicken
	for _, v := range utils.ChickenLocations {
		chicken := animal.NewChicken(chickenImage, v.X, v.Y)
		chickens = append(chickens, chicken)
	}

	// load cows
	embeddedFile, err = EmbeddedAssets.Open(path.Join("assets", "animals", utils.CowImg))
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

	gameObj := game.Game{
		State:       utils.GameStateCustomChar,
		Environment: env,
		CurrentMap:  utils.FarmMap,
		Player:      playerChar,
		Chickens:    chickens,
		Cows:        cows,
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
