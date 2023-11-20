package environment

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"guion-2d-project3/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
)

type Environment struct {
	Maps     []*tiled.Map
	Tilesets map[string]*ebiten.Image
}

func NewEnvironment(embeddedAssets embed.FS, gameMaps []*tiled.Map) *Environment {
	animalsMap, err := utils.LoadMapFromEmbedded(embeddedAssets,
		path.Join("assets", utils.AnimalsMapFile))
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	gameMaps = append(gameMaps, animalsMap)
	forestMap, err := utils.LoadMapFromEmbedded(embeddedAssets,
		path.Join("assets", utils.ForestMapFile))
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	gameMaps = append(gameMaps, forestMap)

	return &Environment{
		Maps:     gameMaps,
		Tilesets: loadTilesets(embeddedAssets),
	}
}

func loadTilesets(embeddedAssets embed.FS) map[string]*ebiten.Image {
	tilesets := map[string]*ebiten.Image{}
	for _, tsPath := range utils.Tilesets {
		embeddedFile, err := embeddedAssets.Open(path.Join("assets", "tilesets", tsPath))
		if err != nil {
			log.Fatal("failed to load embedded image:", embeddedFile, err)
		}
		tsImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
		if err != nil {
			fmt.Println("error loading tileset image")
		}
		tilesets[strings.Split(tsPath, ".")[0]] = tsImage
	}
	return tilesets
}
