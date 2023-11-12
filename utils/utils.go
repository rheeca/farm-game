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
