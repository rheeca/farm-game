package utils

const (
	ProjectTitle    = "Project 3"
	MapFile         = "farm_map.tmx"
	PlayerImg       = "player.png"
	ChickenImg      = "chicken.png"
	FirstTownAudio  = "first-town.wav"
	GroundLayer     = 1
	ObjectsLayer    = 2
	SoundSampleRate = 16000
	UnitSize        = 32
	ToolsUIBoxSize  = 48
)

const (
	StartingX     = 12
	StartingY     = 5
	MovementSpeed = 2
	BackpackSize  = 9
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
	IdleState = iota
	WalkState
	RunState
	HoeState
	AxeState
	WateringState
	PlayerFrameCount   = 8
	PlayerFrameDelay   = 8
	PlayerSpriteWidth  = 96
	PlayerSpriteHeight = 96
)

// animal sprite sheet
const (
	ChickenIdleState      = 0
	AnimalRight           = 0
	AnimalLeft            = 1
	AnimalNumOfDirections = 2
	AnimalFrameCount      = 8
	AnimalFrameDelay      = 12
	AnimalMovementSpeed   = 1
)

// Maps
const (
	FarmMap = iota
	TownMap
	ForestMap
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

	ToolsUIX         int
	ToolsUIY         int
	FarmItemsColumns int
	ToolsFirstSlotX  int
	ToolsFirstSlotY  int
	ToolsFirstBoxX   int
	ToolsFirstBoxY   int

	Tilesets = []string{
		"grass_hill.png",
		"soil_ground.png",
		"hills.png",
		"water.png",
		"darker_soil_ground.png",
		"flowers_stones.png",
		"trees.png",
	}
	ChickenPath1 = []Location{
		{X: 2, Y: 2},
		{X: 6, Y: 2},
	}
	ChickenPath2 = []Location{
		{X: 4, Y: 3},
		{X: 4, Y: 6},
	}
	DogPath = []Location{
		{X: 13, Y: 11},
		{X: 21, Y: 11},
	}
)

// Error messages
const (
	ErrorLoadEmbeddedImage = "failed to load embedded image: %v"
	ErrorLoadEbitenImage   = "failed to load ebiten image: %v"
)
