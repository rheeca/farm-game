package environment

import (
	"embed"
	"fmt"
	"guion-2d-project3/entity/model"
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
	Trees    []model.Object
	Objects  [][]model.Object
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
		Trees:    loadTrees(forestMap),
		Objects:  loadObjects(gameMaps),
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

func loadTrees(tMap *tiled.Map) (trees []model.Object) {
	treePoints := tMap.Groups[0].ObjectGroups[utils.ForestTreePoints].Objects
	for _, t := range treePoints {
		tree := model.Object{
			Type: 0,
			XLoc: int(t.X),
			YLoc: int(t.Y),
			Sprite: model.SpriteBody{
				X:      int(t.X),
				Y:      int(t.Y),
				Width:  64,
				Height: 64,
			},
			Collision: model.CollisionBody{
				X:      int(t.X) + 22,
				Y:      int(t.Y) + 44,
				Width:  20,
				Height: 16,
			},
			Health:      3,
			IsCollision: true,
		}
		trees = append(trees, tree)
	}
	return trees
}

func loadObjects(gameMaps []*tiled.Map) (objects [][]model.Object) {
	var farmObjects, animalsObjects, forestObjects []model.Object
	// crafting table
	ctObj := gameMaps[utils.FarmMap].Groups[0].ObjectGroups[utils.FarmMapCraftingTablePoint].Objects[0]
	craftingTable := model.Object{
		Type: utils.ItemCraftingTable,
		XLoc: int(ctObj.X),
		YLoc: int(ctObj.Y),
		Sprite: model.SpriteBody{
			X:      int(ctObj.X) + 2,
			Y:      int(ctObj.Y) + 20,
			Width:  60,
			Height: 40,
		},
		Collision: model.CollisionBody{
			X:      int(ctObj.X) + 2,
			Y:      int(ctObj.Y) + 20,
			Width:  60,
			Height: 40,
		},
		IsCollision: true,
	}
	farmObjects = append(farmObjects, craftingTable)

	// door
	doorObj := gameMaps[utils.FarmMap].Groups[0].ObjectGroups[utils.FarmMapDoorPoint].Objects[0]
	door := model.Object{
		Type: utils.ItemDoor,
		XLoc: int(doorObj.X),
		YLoc: int(doorObj.Y),
		Sprite: model.SpriteBody{
			X:      int(doorObj.X),
			Y:      int(doorObj.Y),
			Width:  utils.UnitSize,
			Height: utils.UnitSize,
		},
		Collision: model.CollisionBody{
			X:      int(doorObj.X),
			Y:      int(doorObj.Y),
			Width:  utils.UnitSize,
			Height: utils.UnitSize,
		},
		IsCollision: true,
	}
	farmObjects = append(farmObjects, door)

	objects = append(objects, farmObjects, animalsObjects, forestObjects)
	return objects
}
