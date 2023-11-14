package utils

const (
	ProjectTitle      = "Project 3"
	MapFile           = "map.tmx"
	EnvImg            = "environment.png"
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
	MovementSpeed       = 3
	AnimalFrameDelay    = 12
	AnimalMovementSpeed = 1
)

// Directions
const (
	DOWN = iota
	LEFT
	RIGHT
	UP
)

type Location struct {
	X int
	Y int
}

var (
	TileWidth    int
	TileHeight   int
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
