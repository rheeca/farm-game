package utils

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/lafriks/go-tiled"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	ProjectTitle     = "Project 3"
	DefaultPlayerImg = "player_white.png"
	ChickenImg       = "chicken.png"
	CowImg           = "cow.png"
	SoundSampleRate  = 16000
	UnitSize         = 32
)

// game state
const (
	GameStatePlay = iota
	GameStateCustomChar
	GameStateCraft
)

// map
const (
	FarmMapFile    = "farm_map.tmx"
	AnimalsMapFile = "animals_map.tmx"
	ForestMapFile  = "forest_map.tmx"

	FarmMap    = 0
	AnimalsMap = 1
	ForestMap  = 2

	GroundLayer        = 1
	ObjectsLayer       = 2
	FixedObjectsLayer  = 5
	FixedObjects2Layer = 6
	GuideOnlyLayer     = "GuideOnly"

	FarmMapSpawnPoint              = 0
	FarmMapExitToAnimalMapPoint    = 1
	FarmMapExitToForestMapPoint    = 2
	FarmMapEntryFromAnimalMapPoint = 3
	FarmMapEntryFromForestMapPoint = 4
	FarmMapCraftingTablePoint      = 5
	FarmMapDoorPoint               = 6
	FarmMapBedPoint                = 7

	AnimalMapEntryPoint = 0
	AnimalMapExitPoint  = 1

	ForestMapEntryPoint = 0
	ForestMapExitPoint  = 1
	ForestTreePoints    = 2
	ForestRockPoints    = 3
	ForestWoodPoints    = 4
)

const (
	FrameDelay     = 4
	FrameCountSix  = 6
	AnimationDelay = 8
	MovementSpeed  = 2
)

const (
	TreeHitAnimation   = 1
	OpenDoorAnimation  = 0
	CloseDoorAnimation = 1
)

// Directions
const (
	Front = iota
	Back
	Left
	Right
	NumOfDirections
)

// Player sprite sheet
const (
	PlayerFrameCount   = 8
	PlayerFrameDelay   = 8
	PlayerSpriteWidth  = 96
	PlayerSpriteHeight = 96

	IdleState     = 0
	WalkState     = 1
	RunState      = 2
	HoeState      = 3
	AxeState      = 4
	WateringState = 5
)

// animal sprite sheet
const (
	AnimalRight           = 0
	AnimalLeft            = 1
	AnimalNumOfDirections = 2
	AnimalFrameCount      = 8
	AnimalFrameDelay      = 12
	AnimalMovementSpeed   = 1

	ChickenIdleState  = 0
	ChickenHeartState = 8

	CowIdleState    = 0
	CowHeartState   = 8
	CowSpriteWidth  = 64
	CowSpriteHeight = 64
)

// Tilesets
const (
	TilesetGrassHill        = "grass_hill"
	TilesetSoilGround       = "soil_ground"
	TilesetHills            = "hills"
	TilesetWater            = "water"
	TilesetDarkerSoilGround = "darker_soil_ground"
	TilesetFlowersStones    = "flowers_stones"
	TilesetTrees            = "trees"
)

type Location struct {
	X int
	Y int
}

var (
	MapWidth   int
	MapHeight  int
	TileWidth  int
	TileHeight int
	MapColumns int
	MapRows    int

	Tilesets = []string{
		"barn_structures.png",
		"chicken_houses.png",
		"darker_grass_hill.png",
		"darker_soil_ground.png",
		"door.png",
		"fences.png",
		"flowers_stones.png",
		"furniture.png",
		"grass_hill.png",
		"hills.png",
		"paths.png",
		"soil_ground.png",
		"trees.png",
		"water.png",
		"water_tray.png",
		"wood_bridge.png",
		"wooden_house.png",
	}
	ChickenLocations = []Location{
		{X: 5, Y: 5},
		{X: 6, Y: 7},
	}
	CowLocations = []Location{
		{X: 5, Y: 11},
		{X: 7, Y: 13},
	}
	CollisionLayers = []int{
		ObjectsLayer, FixedObjectsLayer, FixedObjects2Layer,
	}
)

// Error messages
const (
	ErrorLoadEmbeddedImage = "failed to load embedded image: %v"
	ErrorLoadEbitenImage   = "failed to load ebiten image: %v"
)

func LoadMapFromEmbedded(EmbeddedAssets embed.FS, name string) (*tiled.Map, error) {
	embeddedMap, err := tiled.LoadFile(name,
		tiled.WithFileSystem(EmbeddedAssets))
	if err != nil {
		return nil, err
	}
	return embeddedMap, nil
}

func LoadFont(size int) font.Face {
	//originally inspired by https://www.fatoldyeti.com/posts/roguelike16/
	trueTypeFont, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		fmt.Println("error loading font:", err)
	}
	fontFace, err := opentype.NewFace(trueTypeFont, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		fmt.Println("error loading font of correct size:", err)
	}
	return fontFace
}
