package main

import (
	"embed"
	"fmt"
	"image"
	"log"
	"os"
	"path"
	"time"

	"guion-2d-project3/entity/animal"
	"guion-2d-project3/entity/environment"
	"guion-2d-project3/entity/player"
	"guion-2d-project3/interfaces"
	"guion-2d-project3/utils"

	"github.com/co0p/tankism/lib/collision"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
)

//go:embed assets/*
var EmbeddedAssets embed.FS

type Game struct {
	Environment  *environment.Environment
	Player       *player.Player
	Animals      []*animal.Animal
	CurrentFrame int
}

func hasMapCollisions(g *Game, animObj interfaces.AnimatedSprite) bool {
	for tileY := 0; tileY < utils.MapTileHeight; tileY += 1 {
		for tileX := 0; tileX < utils.MapTileWidth; tileX += 1 {
			tile := g.Environment.Maps[0].Layers[utils.CollisionObjLayer].Tiles[tileY*utils.MapTileWidth+tileX]
			if tile.ID == 0 {
				continue
			}
			tileXpos := utils.TileWidth * tileX
			tileYpos := utils.TileHeight * tileY

			newX := animObj.GetXLoc() + animObj.GetDx()
			newY := animObj.GetYLoc() + animObj.GetDy()
			animBounds := collision.BoundingBox{
				// bounding box for animated object made slightly smaller than the sprite
				X:      float64(newX + animObj.GetWidth()/4),
				Y:      float64(newY + animObj.GetHeight()/2),
				Width:  float64(animObj.GetWidth() / 2),
				Height: float64(animObj.GetHeight() / 2),
			}
			tileBounds := collision.BoundingBox{
				X:      float64(tileXpos),
				Y:      float64(tileYpos),
				Width:  float64(utils.TileWidth),
				Height: float64(utils.TileHeight),
			}
			if collision.AABBCollision(animBounds, tileBounds) {
				return true
			}
		}
	}
	return false
}

func playerHasCollisions(g *Game) bool {
	// TODO: check for map collisions
	//if hasMapCollisions(g, g.Player) {
	//	return true
	//}

	// check for animated entities collisions
	for _, a := range g.Animals {
		if g.Player.HasCollisionWith(a) {
			return true
		}
	}
	return false
}

func animalHasCollisions(g *Game, animObj interfaces.AnimatedSprite) bool {
	// TODO: check for map collisions
	//if hasMapCollisions(g, animObj) {
	//	return true
	//}

	// check for collision with player
	if animObj.HasCollisionWith(g.Player) {
		return true
	}
	return false
}

func getPlayerInput(g *Game) {
	g.Player.UpdateFrame(g.CurrentFrame)
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.Player.XLoc > 0 {
		g.Player.Direction = utils.Left
		g.Player.State = utils.WalkState
		g.Player.Dx -= utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dx = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) &&
		g.Player.XLoc < utils.MapWidth-g.Player.SpriteWidth {
		g.Player.Direction = utils.Right
		g.Player.State = utils.WalkState
		g.Player.Dx += utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dx = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.Player.YLoc > 0 {
		g.Player.Direction = utils.Back
		g.Player.State = utils.WalkState
		g.Player.Dy -= utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dy = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) &&
		g.Player.YLoc < utils.MapHeight-g.Player.SpriteHeight {
		g.Player.Direction = utils.Front
		g.Player.State = utils.WalkState
		g.Player.Dy += utils.MovementSpeed
		if !playerHasCollisions(g) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dy = 0
		}
	} else {
		g.Player.State = utils.IdleState
	}
}

func updateAnimals(g *Game) {
	for i, a := range g.Animals {
		if a.XLoc == (a.Path[a.Destination].X*utils.TileWidth) &&
			a.YLoc == (a.Path[a.Destination].Y*utils.TileHeight) {

			// if animal has reached its destination, give it a new destination
			g.Animals[i].Destination = (g.Animals[i].Destination + 1) % len(a.Path)
		} else {
			// move animal towards destination
			if a.XLoc > a.Path[a.Destination].X*utils.TileWidth {
				g.Animals[i].Direction = utils.Left
				g.Animals[i].UpdateFrame(g.CurrentFrame)
				g.Animals[i].Dx -= utils.AnimalMovementSpeed
				if !animalHasCollisions(g, g.Animals[i]) {
					g.Animals[i].UpdateLocation()
				} else {
					g.Animals[i].Dx = 0
				}
			} else if a.XLoc < a.Path[a.Destination].X*utils.TileWidth {
				g.Animals[i].Direction = utils.Right
				g.Animals[i].UpdateFrame(g.CurrentFrame)
				g.Animals[i].Dx += utils.AnimalMovementSpeed
				if !animalHasCollisions(g, g.Animals[i]) {
					g.Animals[i].UpdateLocation()
				} else {
					g.Animals[i].Dx = 0
				}
			} else if a.YLoc > a.Path[a.Destination].Y*utils.TileHeight {
				g.Animals[i].Direction = utils.Back
				g.Animals[i].UpdateFrame(g.CurrentFrame)
				g.Animals[i].Dy -= utils.AnimalMovementSpeed
				if !animalHasCollisions(g, g.Animals[i]) {
					g.Animals[i].UpdateLocation()
				} else {
					g.Animals[i].Dy = 0
				}
			} else if a.YLoc < a.Path[a.Destination].Y*utils.TileHeight {
				g.Animals[i].Direction = utils.Front
				g.Animals[i].UpdateFrame(g.CurrentFrame)
				g.Animals[i].Dy += utils.AnimalMovementSpeed
				if !animalHasCollisions(g, g.Animals[i]) {
					g.Animals[i].UpdateLocation()
				} else {
					g.Animals[i].Dy = 0
				}
			}
		}
	}
}

func drawMap(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions, currMap int) {
	for _, layer := range g.Environment.Maps[currMap].Layers {
		for tileY := 0; tileY < utils.MapTileHeight; tileY += 1 {
			for tileX := 0; tileX < utils.MapTileWidth; tileX += 1 {
				// find img of tile to draw
				tileToDraw := layer.Tiles[tileY*utils.MapTileWidth+tileX]
				if tileToDraw.IsNil() {
					continue
				}

				tileToDrawX := int(tileToDraw.ID) % tileToDraw.Tileset.Columns
				tileToDrawY := int(tileToDraw.ID) / tileToDraw.Tileset.Columns

				ebitenTileToDraw := g.Environment.Tilesets[tileToDraw.Tileset.Name].SubImage(image.Rect(tileToDrawX*utils.TileWidth,
					tileToDrawY*utils.TileHeight,
					tileToDrawX*utils.TileWidth+utils.TileWidth,
					tileToDrawY*utils.TileHeight+utils.TileHeight)).(*ebiten.Image)

				// draw tile
				drawOptions.GeoM.Reset()
				TileXpos := float64(utils.TileWidth * tileX)
				TileYpos := float64(utils.TileHeight * tileY)
				drawOptions.GeoM.Translate(TileXpos, TileYpos)
				screen.DrawImage(ebitenTileToDraw, &drawOptions)
			}
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

	drawMap(g, screen, drawOptions, g.Environment.CurrentMap)

	// draw animals
	for _, a := range g.Animals {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(a.XLoc), float64(a.YLoc))
		screen.DrawImage(a.Spritesheet.SubImage(image.Rect(a.Frame*a.Width,
			a.Direction*a.Height,
			a.Frame*a.Width+a.Width,
			a.Direction*a.Height+a.Height)).(*ebiten.Image), &drawOptions)
	}

	// draw player
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(g.Player.XLoc), float64(g.Player.YLoc))
	screen.DrawImage(g.Player.Spritesheet.SubImage(image.Rect(g.Player.Frame*g.Player.SpriteWidth,
		(g.Player.State*utils.NumOfDirections+g.Player.Direction)*g.Player.SpriteHeight,
		g.Player.Frame*g.Player.SpriteWidth+g.Player.SpriteWidth,
		(g.Player.State*utils.NumOfDirections+g.Player.Direction)*g.Player.SpriteHeight+g.Player.SpriteHeight)).(*ebiten.Image), &drawOptions)
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
	gameMap, err := loadMapFromEmbedded(path.Join("assets", utils.MapFile))
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	windowWidth := gameMap.Width * gameMap.TileWidth
	windowHeight := gameMap.Height * gameMap.TileHeight
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(utils.ProjectTitle)

	utils.MapWidth = gameMap.Width * gameMap.TileWidth
	utils.MapHeight = gameMap.Height * gameMap.TileHeight
	utils.TileWidth = gameMap.TileWidth
	utils.TileHeight = gameMap.TileHeight
	utils.MapTileWidth = gameMap.Width
	utils.MapTileHeight = gameMap.Height

	// load environment
	env := environment.NewEnvironment(EmbeddedAssets, []*tiled.Map{gameMap})

	// load audio
	audioContext := audio.NewContext(utils.SoundSampleRate)
	bgmPlayer, err := loadWavFromEmbedded(utils.FirstTownAudio, audioContext)
	if err != nil {
		fmt.Println("shutting down. error:", err)
		return
	}

	// load player
	embeddedFile, err := EmbeddedAssets.Open(path.Join("assets", utils.PlayerImg))
	if err != nil {
		log.Fatal("failed to load embedded image:", embeddedFile, err)
	}
	playerImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("error loading player image")
	}
	playerChar := player.NewPlayer(playerImage)

	// load chickens
	embeddedFile, err = EmbeddedAssets.Open(path.Join("assets", utils.ChickenImg))
	if err != nil {
		log.Fatal("failed to load embedded image:", embeddedFile, err)
	}
	chickenImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("error loading chicken image")
	}
	chicken1 := animal.NewAnimal(chickenImage, utils.ChickenPath1)
	chicken2 := animal.NewAnimal(chickenImage, utils.ChickenPath2)

	// load dog
	embeddedFile, err = EmbeddedAssets.Open(path.Join("assets", utils.DogImg))
	if err != nil {
		log.Fatal("failed to load embedded image:", embeddedFile, err)
	}
	dogImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("error loading dog image")
	}
	dog := animal.NewAnimal(dogImage, utils.DogPath)

	game := Game{
		Environment: env,
		Player:      playerChar,
		Animals:     []*animal.Animal{chicken1, chicken2, dog},
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
