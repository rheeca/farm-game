package utils

const (
	ProjectTitle      = "Project 3"
	MapFile           = "farm_map.tmx"
	PlayerImg         = "player.png"
	ChickenImg        = "chicken.png"
	DogImg            = "dog.png"
	FirstTownAudio    = "first-town.wav"
	GroundLayer       = 0
	CollisionObjLayer = 1
	SoundSampleRate   = 16000
)

const (
	StartingX           = 12
	StartingY           = 5
	StartingFrame       = 0
	AnimFrameCount      = 4
	FrameDelay          = 4
	MovementSpeed       = 2
	AnimalFrameDelay    = 12
	AnimalMovementSpeed = 1
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
	PlayerFrameCount   = 8
	PlayerFrameDelay   = 8
	PlayerSpriteWidth  = 96
	PlayerSpriteHeight = 96
)

// Maps
const (
	FarmMap = iota
	TownMap
	ForestMap
)

type Location struct {
	X int
	Y int
}

var (
	MapWidth      int
	MapHeight     int
	TileWidth     int
	TileHeight    int
	MapTileWidth  int
	MapTileHeight int

	Tilesets = []string{
		"grass_hill.png",
		"soil_ground.png",
		"hills.png",
		"water.png",
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
