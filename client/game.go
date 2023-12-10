package main

import (
	"embed"
	"encoding/json"
	"guion-2d-project3/entity/game"
	"guion-2d-project3/entity/loader"
	"guion-2d-project3/entity/model"
	"guion-2d-project3/utils"
	"log"

	"github.com/codecat/go-enet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lafriks/go-tiled"
)

type ClientGame struct {
	peer         enet.Peer
	host         enet.Host
	State        int
	Data         *game.GameData
	Maps         []*tiled.Map
	CurrentMap   int
	CurrentFrame int
	PlayerID     string
	Images       loader.ImageCollection
	Sounds       loader.SoundCollection
	UIState      model.UIState
}

func NewClientGame(embeddedAssets embed.FS, peer enet.Peer, host enet.Host) ClientGame {
	gameMaps := game.LoadMaps(embeddedAssets, "assets")
	currentMap := utils.FarmMap
	windowWidth := gameMaps[currentMap].Width * gameMaps[currentMap].TileWidth
	windowHeight := gameMaps[currentMap].Height * gameMaps[currentMap].TileHeight
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(utils.ProjectTitle)

	images := loader.NewImageCollection(embeddedAssets, "assets")
	sounds := loader.NewSoundCollection(embeddedAssets, "assets")
	game.SetConstants(gameMaps[currentMap], images)

	return ClientGame{
		peer:       peer,
		host:       host,
		State:      utils.GameStateWaitingForServer,
		Maps:       gameMaps,
		CurrentMap: currentMap,
		PlayerID:   peer.GetAddress().String(),
		Images:     images,
		Sounds:     sounds,
	}
}

func (g *ClientGame) Update() error {
	if g.State != utils.GameStateWaitingForServer {
		sendUpdatesToServer(g)
	}
	listenForEvents(g)
	return nil
}

func (g *ClientGame) Draw(screen *ebiten.Image) {
	drawOptions := ebiten.DrawImageOptions{}
	game.DrawMap(g.Maps[g.CurrentMap], g.Images.Tilesets, screen, drawOptions)

	if g.State == utils.GameStateWaitingForServer {
		return
	}

	if g.Data != nil {
		if g.Data.Environment != nil {
			if g.CurrentMap == utils.ForestMap {
				game.DrawTrees(g.Data.Environment.Trees, g.Images, screen, drawOptions)
			} else if g.CurrentMap == utils.FarmMap {
				game.DrawFarmPlots(g.Data.Environment.Plots, g.Images, screen, drawOptions)
			}
			game.DrawObjects(g.Data.Environment.Objects[g.CurrentMap], g.Images, screen, drawOptions)
		}

		if g.CurrentMap == utils.AnimalsMap {
			game.DrawChickens(g.Data.Chickens, screen, drawOptions)
			game.DrawCows(g.Data.Cows, screen, drawOptions)
		}

		game.DrawPlayers(g.CurrentMap, g.Data.Players, g.Images, screen, drawOptions)
		game.DrawBackpack(g.Data.Players[g.PlayerID], g.Images, screen, drawOptions)
	}
}

func (g *ClientGame) Layout(oWidth, oHeight int) (sWidth, sHeight int) {
	return oWidth, oHeight
}

func listenForEvents(g *ClientGame) {
	ev := g.host.Service(1000)
	switch ev.GetType() {
	case enet.EventConnect:
		log.Println("connected to the server!")
		g.State = utils.GameStatePlay

	case enet.EventDisconnect:
		log.Println("lost connection to the server!")

	case enet.EventReceive:
		packet := ev.GetPacket()

		var data model.DataPacket
		json.Unmarshal(packet.GetData(), &data)
		if data.Type == utils.PacketGameData {
			var gamePacket game.GameData
			body, _ := json.Marshal(data.Body)
			json.Unmarshal(body, &gamePacket)
			g.Data = &gamePacket
			g.CurrentMap = g.Data.Players[g.PlayerID].CurrentMap
		}

		packet.Destroy()
	}
}

func sendUpdatesToServer(g *ClientGame) {
	clientInput := getInput(g)
	dataObj := model.DataPacket{
		Type: utils.PacketClientInput,
		Body: clientInput,
	}
	data, _ := json.Marshal(dataObj)
	g.peer.SendString(string(data), 0, enet.PacketFlagReliable)
}

func getInput(g *ClientGame) (input model.ClientInputPacket) {
	input.PlayerID = g.PlayerID
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		input.Input = utils.InputKeyW
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		input.Input = utils.InputKeyA
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		input.Input = utils.InputKeyS
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		input.Input = utils.InputKeyD
	} else if ebiten.IsKeyPressed(ebiten.Key1) {
		input.Input = utils.InputKey1
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		input.Input = utils.InputKey2
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		input.Input = utils.InputKey3
	} else if ebiten.IsKeyPressed(ebiten.Key4) {
		input.Input = utils.InputKey4
	} else if ebiten.IsKeyPressed(ebiten.Key5) {
		input.Input = utils.InputKey5
	} else if ebiten.IsKeyPressed(ebiten.Key6) {
		input.Input = utils.InputKey6
	} else if ebiten.IsKeyPressed(ebiten.Key7) {
		input.Input = utils.InputKey7
	} else if ebiten.IsKeyPressed(ebiten.Key8) {
		input.Input = utils.InputKey8
	} else if ebiten.IsKeyPressed(ebiten.Key9) {
		input.Input = utils.InputKey9
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		input.Input = utils.InputMouseLeft
		input.MouseX = mouseX
		input.MouseY = mouseY
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		mouseX, mouseY := ebiten.CursorPosition()
		input.Input = utils.InputMouseRight
		input.MouseX = mouseX
		input.MouseY = mouseY
	} else {
		input.Input = utils.InputNone
	}

	return input
}
