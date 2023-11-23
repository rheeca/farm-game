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
	Plots    []model.Plot
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

	// bed
	bedObj := gameMaps[utils.FarmMap].Groups[0].ObjectGroups[utils.FarmMapBedPoint].Objects[0]
	bed := model.Object{
		Type: utils.ItemBedPink,
		XLoc: int(bedObj.X),
		YLoc: int(bedObj.Y),
		Sprite: model.SpriteBody{
			X:      int(bedObj.X) + 2,
			Y:      int(bedObj.Y) + 20,
			Width:  28,
			Height: 44,
		},
		Collision: model.CollisionBody{
			X:      int(bedObj.X) + 2,
			Y:      int(bedObj.Y) + 20,
			Width:  28,
			Height: 44,
		},
		IsCollision: true,
	}
	farmObjects = append(farmObjects, bed)

	// forest map
	// rocks
	for _, stoneObj := range gameMaps[utils.ForestMap].Groups[0].ObjectGroups[utils.ForestRockPoints].Objects {
		stone := model.Object{
			Type: utils.ItemMapStone3,
			XLoc: int(stoneObj.X),
			YLoc: int(stoneObj.Y),
			Sprite: model.SpriteBody{
				X:      int(stoneObj.X),
				Y:      int(stoneObj.Y),
				Width:  utils.UnitSize,
				Height: utils.UnitSize,
			},
			Collision: model.CollisionBody{
				X:      int(stoneObj.X),
				Y:      int(stoneObj.Y),
				Width:  utils.UnitSize,
				Height: utils.UnitSize,
			},
			IsCollision: true,
		}
		forestObjects = append(forestObjects, stone)
	}
	// wood
	for _, woodObj := range gameMaps[utils.ForestMap].Groups[0].ObjectGroups[utils.ForestWoodPoints].Objects {
		wood := model.Object{
			Type: utils.ItemMapWood,
			XLoc: int(woodObj.X),
			YLoc: int(woodObj.Y),
			Sprite: model.SpriteBody{
				X:      int(woodObj.X),
				Y:      int(woodObj.Y),
				Width:  utils.UnitSize,
				Height: utils.UnitSize,
			},
			Collision: model.CollisionBody{
				X:      int(woodObj.X),
				Y:      int(woodObj.Y),
				Width:  utils.UnitSize,
				Height: utils.UnitSize,
			},
			IsCollision: true,
		}
		forestObjects = append(forestObjects, wood)
	}

	forestObjects = loadObject32(forestObjects, utils.ForestSunflowerPoints, utils.MapSunflower, gameMaps[utils.ForestMap])
	forestObjects = loadObject32(forestObjects, utils.ForestBlueflowerPoints, utils.MapBlueflower, gameMaps[utils.ForestMap])
	forestObjects = loadObject32(forestObjects, utils.ForestWeedPoints, utils.MapWeed, gameMaps[utils.ForestMap])
	forestObjects = loadObject32(forestObjects, utils.ForestPinkDyeFlowerPoints, utils.MapPinkDyeFlower, gameMaps[utils.ForestMap])
	forestObjects = loadObject32(forestObjects, utils.ForestBlueDyeFlowerPoints, utils.MapBlueDyeFlower, gameMaps[utils.ForestMap])

	objects = append(objects, farmObjects, animalsObjects, forestObjects)
	return objects
}

func (e *Environment) ResetDay() {
	for i := range e.Trees {
		e.Trees[i].IsNil = false
	}
	for i, o := range e.Objects[utils.ForestMap] {
		if o.IsNil && (o.Type == utils.ItemMapStone3 || o.Type == utils.ItemMapWood) {
			e.Objects[utils.ForestMap][i].IsNil = false
		}
	}
	for i, p := range e.Plots {
		if p.IsWatered && p.HasPlant {
			e.Plots[i].ReadyForHarvest = true
		}
		e.Plots[i].IsWatered = false
	}
}

func (e *Environment) AddPlot(tileX, tileY int) {
	e.Plots = append(e.Plots, model.Plot{
		XTile: tileX,
		YTile: tileY,
	})
}

func (e *Environment) WaterPlot(tileX, tileY int) {
	for i, plot := range e.Plots {
		if plot.XTile == tileX && plot.YTile == tileY {
			e.Plots[i].IsWatered = true
		}
	}
}

func (e *Environment) PlantSeedInPlot(tileX, tileY, plantType int) bool {
	for i, plot := range e.Plots {
		if plot.XTile == tileX && plot.YTile == tileY && !e.Plots[i].HasPlant {
			e.Plots[i].HasPlant = true
			e.Plots[i].PlantType = plantType
			return true
		}
	}
	return false
}

func (e *Environment) HarvestPlant(tileX, tileY int) (harvested bool, plantType int) {
	for i, plot := range e.Plots {
		if plot.XTile == tileX && plot.YTile == tileY && e.Plots[i].ReadyForHarvest {
			e.Plots[i].HasPlant = false
			e.Plots[i].ReadyForHarvest = false
			return true, e.Plots[i].PlantType
		}
	}
	return false, 0
}
