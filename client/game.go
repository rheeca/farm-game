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
		State:      utils.GameStateCustomChar,
		Maps:       gameMaps,
		CurrentMap: currentMap,
		Images:     images,
		Sounds:     sounds,
	}
}

func (g *ClientGame) Update() error {
	listenForEvents(g)
	return nil
}

func (g *ClientGame) Draw(screen *ebiten.Image) {
	drawOptions := ebiten.DrawImageOptions{}
	game.DrawMap(g.Maps[g.CurrentMap], g.Images.Tilesets, screen, drawOptions)
}

func (g *ClientGame) Layout(oWidth, oHeight int) (sWidth, sHeight int) {
	return oWidth, oHeight
}

func listenForEvents(g *ClientGame) {
	ev := g.host.Service(1000)
	switch ev.GetType() {
	case enet.EventConnect:
		log.Println("connected to the server!")

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
		}

		packet.Destroy()
	}
}
