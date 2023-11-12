package main

import (
	"embed"
	"fmt"
	"github.com/co0p/tankism/lib/collision"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
	"image"
	"log"
	"os"
	"path"
	"time"
)

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

//go:embed assets/*
var EmbeddedAssets embed.FS

type Game struct {
	Environment  Environment
	Player       Player
	Animals      []Animal
	CurrentFrame int
}

type Environment struct {
	Map       *tiled.Map
	Tileset   *ebiten.Image
	MapWidth  int
	MapHeight int
}

type Player struct {
	Spritesheet *ebiten.Image
	Frame       int
	Direction   int
	xLoc        int
	yLoc        int
	Width       int
	Height      int
}

type Animal struct {
	Spritesheet *ebiten.Image
	Frame       int
	Direction   int
	xLoc        int
	yLoc        int
	Destination int
	Path        []Location
	Width       int
	Height      int
}

type Location struct {
	X int
	Y int
}

func hasCollision(g *Game, playerX, playerY, objectX, objectY int) bool {
	playerBounds := collision.BoundingBox{
		// player bounding box made slightly smaller than the sprite
		X:      float64(playerX + g.Player.Width/4),
		Y:      float64(playerY + g.Player.Height/2),
		Width:  float64(g.Player.Width / 2),
		Height: float64(g.Player.Height / 2),
	}
	objectBounds := collision.BoundingBox{
		X:      float64(objectX),
		Y:      float64(objectY),
		Width:  float64(g.Environment.Map.TileWidth),
		Height: float64(g.Environment.Map.TileHeight),
	}
	if collision.AABBCollision(playerBounds, objectBounds) {
		return true
	}
	return false
}

func hasMapCollisions(g *Game, playerX, playerY int) bool {
	for tileY := 0; tileY < g.Environment.Map.Height; tileY += 1 {
		for tileX := 0; tileX < g.Environment.Map.Width; tileX += 1 {
			tileToDraw := g.Environment.Map.Layers[CollisionObjLayer].Tiles[tileY*g.Environment.Map.Width+tileX]
			if tileToDraw.ID == 0 {
				continue
			}
			tileXpos := g.Environment.Map.TileWidth * tileX
			tileYpos := g.Environment.Map.TileHeight * tileY
			if hasCollision(g, playerX, playerY, tileXpos, tileYpos) {
				return true
			}
		}
	}
	return false
}

func hasAnimalCollisions(g *Game, playerX, playerY int) bool {
	for _, animal := range g.Animals {
		if hasCollision(g, playerX, playerY, animal.xLoc, animal.yLoc) {
			return true
		}
	}
	return false
}

func updatePlayerFrame(g *Game) {
	if g.CurrentFrame%FrameDelay == 0 {
		g.Player.Frame += 1
		if g.Player.Frame >= AnimFrameCount {
			g.Player.Frame = 0
		}
	}
}

func updateAnimalFrame(g *Game, animal Animal) (frame int) {
	frame = animal.Frame
	if g.CurrentFrame%AnimalFrameDelay == 0 {
		frame += 1
		if frame >= AnimFrameCount {
			frame = 0
		}
	}
	return frame
}

func getPlayerInput(g *Game) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.Player.xLoc > 0 {
		g.Player.Direction = LEFT
		updatePlayerFrame(g)
		newX := g.Player.xLoc - MovementSpeed
		if !hasMapCollisions(g, newX, g.Player.yLoc) && !hasAnimalCollisions(g, newX, g.Player.yLoc) {
			g.Player.xLoc = newX
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) &&
		g.Player.xLoc < g.Environment.MapWidth-g.Player.Width {
		g.Player.Direction = RIGHT
		updatePlayerFrame(g)
		newX := g.Player.xLoc + MovementSpeed
		if !hasMapCollisions(g, newX, g.Player.yLoc) && !hasAnimalCollisions(g, newX, g.Player.yLoc) {
			g.Player.xLoc = newX
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.Player.yLoc > 0 {
		g.Player.Direction = UP
		updatePlayerFrame(g)
		newY := g.Player.yLoc - MovementSpeed
		if !hasMapCollisions(g, g.Player.xLoc, newY) && !hasAnimalCollisions(g, g.Player.xLoc, newY) {
			g.Player.yLoc = newY
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) &&
		g.Player.yLoc < g.Environment.MapHeight-g.Player.Height {
		g.Player.Direction = DOWN
		updatePlayerFrame(g)
		newY := g.Player.yLoc + MovementSpeed
		if !hasMapCollisions(g, g.Player.xLoc, newY) && !hasAnimalCollisions(g, g.Player.xLoc, newY) {
			g.Player.yLoc = newY
		}
	} else {
		g.Player.Frame = StartingFrame
	}
}

func updateAnimals(g *Game) {
	for i, animal := range g.Animals {
		if animal.xLoc == (animal.Path[animal.Destination].X*TileWidth) &&
			animal.yLoc == (animal.Path[animal.Destination].Y*TileHeight) {

			// if animal has reached its destination, give it a new destination
			g.Animals[i].Destination = (g.Animals[i].Destination + 1) % len(animal.Path)
		} else {
			// move animal towards destination
			if animal.xLoc > animal.Path[animal.Destination].X*TileWidth {
				g.Animals[i].Direction = LEFT
				g.Animals[i].Frame = updateAnimalFrame(g, animal)
				newX := g.Animals[i].xLoc - AnimalMovementSpeed
				if !hasCollision(g, g.Player.xLoc, g.Player.yLoc, newX, g.Animals[i].yLoc) {
					g.Animals[i].xLoc = newX
				}
			} else if animal.xLoc < animal.Path[animal.Destination].X*TileWidth {
				g.Animals[i].Direction = RIGHT
				g.Animals[i].Frame = updateAnimalFrame(g, animal)
				newX := g.Animals[i].xLoc + AnimalMovementSpeed
				if !hasCollision(g, g.Player.xLoc, g.Player.yLoc, newX, g.Animals[i].yLoc) {
					g.Animals[i].xLoc = newX
				}
			} else if animal.yLoc > animal.Path[animal.Destination].Y*TileHeight {
				g.Animals[i].Direction = UP
				g.Animals[i].Frame = updateAnimalFrame(g, animal)
				newY := g.Animals[i].yLoc - AnimalMovementSpeed
				if !hasCollision(g, g.Player.xLoc, g.Player.yLoc, g.Animals[i].xLoc, newY) {
					g.Animals[i].yLoc = newY
				}
			} else if animal.yLoc < animal.Path[animal.Destination].Y*TileHeight {
				g.Animals[i].Direction = DOWN
				g.Animals[i].Frame = updateAnimalFrame(g, animal)
				newY := g.Animals[i].yLoc + AnimalMovementSpeed
				if !hasCollision(g, g.Player.xLoc, g.Player.yLoc, g.Animals[i].xLoc, newY) {
					g.Animals[i].yLoc = newY
				}
			}
		}
	}
}

func drawMapLayer(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions, layer int) {
	tilesetColumns := g.Environment.Map.Tilesets[0].Columns
	for tileY := 0; tileY < g.Environment.Map.Height; tileY += 1 {
		for tileX := 0; tileX < g.Environment.Map.Width; tileX += 1 {
			drawOptions.GeoM.Reset()
			TileXpos := float64(TileWidth * tileX)
			TileYpos := float64(TileHeight * tileY)
			drawOptions.GeoM.Translate(TileXpos, TileYpos)
			tileToDraw := g.Environment.Map.Layers[layer].Tiles[tileY*g.Environment.Map.Width+tileX]
			if tileToDraw.ID == 0 {
				continue
			}
			tileToDrawX := int(tileToDraw.ID) % tilesetColumns
			tileToDrawY := int(tileToDraw.ID) / tilesetColumns

			ebitenTileToDraw := g.Environment.Tileset.SubImage(image.Rect(tileToDrawX*TileWidth,
				tileToDrawY*TileHeight,
				tileToDrawX*TileWidth+TileWidth,
				tileToDrawY*TileHeight+TileHeight)).(*ebiten.Image)
			screen.DrawImage(ebitenTileToDraw, &drawOptions)
		}
	}
}

func (g *Game) Update() error {
	g.CurrentFrame += 1
	getPlayerInput(g)
	updateAnimals(g)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawOptions := ebiten.DrawImageOptions{}

	// draw map ground
	drawMapLayer(g, screen, drawOptions, GroundLayer)

	// draw map objects
	drawMapLayer(g, screen, drawOptions, CollisionObjLayer)

	// draw animals
	for _, animal := range g.Animals {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(animal.xLoc), float64(animal.yLoc))
		screen.DrawImage(animal.Spritesheet.SubImage(image.Rect(animal.Frame*animal.Width,
			animal.Direction*animal.Height,
			animal.Frame*animal.Width+animal.Width,
			animal.Direction*animal.Height+animal.Height)).(*ebiten.Image), &drawOptions)
	}

	// draw player
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(g.Player.xLoc), float64(g.Player.yLoc))
	screen.DrawImage(g.Player.Spritesheet.SubImage(image.Rect(g.Player.Frame*g.Player.Width,
		g.Player.Direction*g.Player.Height,
		g.Player.Frame*g.Player.Width+g.Player.Width,
		g.Player.Direction*g.Player.Height+g.Player.Height)).(*ebiten.Image), &drawOptions)
}

func (g *Game) Layout(oWidth, oHeight int) (sWidth, sHeight int) {
	return oWidth, oHeight
}

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
	gameMap, err := loadMapFromEmbedded(path.Join("assets", MapFile))
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	windowWidth := gameMap.Width * gameMap.TileWidth
	windowHeight := gameMap.Height * gameMap.TileHeight
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(ProjectTitle)

	TileWidth = gameMap.TileWidth
	TileHeight = gameMap.TileHeight

	// load environment
	embeddedFile, err := EmbeddedAssets.Open(path.Join("assets", EnvImg))
	if err != nil {
		log.Fatal("failed to load embedded image:", embeddedFile, err)
	}
	mapImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("error loading map image")
	}
	environment := Environment{
		Map:       gameMap,
		Tileset:   mapImage,
		MapWidth:  gameMap.Width * gameMap.TileWidth,
		MapHeight: gameMap.Height * gameMap.TileHeight,
	}

	// load audio
	audioContext := audio.NewContext(SoundSampleRate)
	bgmPlayer, err := loadWavFromEmbedded(FirstTownAudio, audioContext)
	if err != nil {
		fmt.Println("shutting down. error:", err)
		return
	}

	// load player
	embeddedFile, err = EmbeddedAssets.Open(path.Join("assets", PlayerImg))
	if err != nil {
		log.Fatal("failed to load embedded image:", embeddedFile, err)
	}
	playerImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("error loading player image")
	}
	player := Player{
		Spritesheet: playerImage,
		xLoc:        StartingX * gameMap.TileWidth,
		yLoc:        StartingY * gameMap.TileHeight,
		Width:       playerImage.Bounds().Dx() / AnimFrameCount,
		Height:      playerImage.Bounds().Dy() / AnimFrameCount,
	}

	// load chickens
	embeddedFile, err = EmbeddedAssets.Open(path.Join("assets", ChickenImg))
	if err != nil {
		log.Fatal("failed to load embedded image:", embeddedFile, err)
	}
	chickenImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("error loading chicken image")
	}
	chicken1 := Animal{
		Spritesheet: chickenImage,
		xLoc:        ChickenPath1[0].X * gameMap.TileWidth,
		yLoc:        ChickenPath1[0].Y * gameMap.TileHeight,
		Path:        ChickenPath1,
		Width:       chickenImage.Bounds().Dx() / AnimFrameCount,
		Height:      chickenImage.Bounds().Dy() / AnimFrameCount,
	}
	chicken2 := Animal{
		Spritesheet: chickenImage,
		xLoc:        ChickenPath2[0].X * gameMap.TileWidth,
		yLoc:        ChickenPath2[0].Y * gameMap.TileHeight,
		Path:        ChickenPath2,
		Width:       chickenImage.Bounds().Dx() / AnimFrameCount,
		Height:      chickenImage.Bounds().Dy() / AnimFrameCount,
	}

	// load dog
	embeddedFile, err = EmbeddedAssets.Open(path.Join("assets", DogImg))
	if err != nil {
		log.Fatal("failed to load embedded image:", embeddedFile, err)
	}
	dogImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("error loading dog image")
	}
	dog := Animal{
		Spritesheet: dogImage,
		xLoc:        DogPath[0].X * gameMap.TileWidth,
		yLoc:        DogPath[0].Y * gameMap.TileHeight,
		Path:        DogPath,
		Width:       dogImage.Bounds().Dx() / AnimFrameCount,
		Height:      dogImage.Bounds().Dy() / AnimFrameCount,
	}

	game := Game{
		Environment: environment,
		Player:      player,
		Animals:     []Animal{chicken1, chicken2, dog},
	}

	go func(player *audio.Player) {
		player.Play()
		time.Sleep(122 * time.Second)
		player.Rewind()
	}(bgmPlayer)
	err = ebiten.RunGame(&game)
	if err != nil {
		fmt.Println("failed to run game:", err)
	}
}
