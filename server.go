package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"guion-2d-project3/entity/game"
	"guion-2d-project3/entity/model"
	"guion-2d-project3/utils"
	"log"

	"github.com/codecat/go-enet"
	"github.com/gofrs/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed client/assets/*
var EmbeddedAssets embed.FS

func main() {
	enet.Initialize()
	host, err := enet.NewHost(enet.NewListenAddress(utils.ServerPort), 32, 1, 0, 0)
	if err != nil {
		log.Fatal("failed to create host: ", err.Error())
		return
	}

	gameObj := game.NewGame(EmbeddedAssets)

	go runServer(host, gameObj)
	err = ebiten.RunGame(&gameObj)
	if err != nil {
		fmt.Println("failed to run game:", err)
	}
}

func runServer(host enet.Host, g game.Game) {
	defer func() {
		host.Destroy()
		enet.Deinitialize()
	}()
	for {
		ev := host.Service(1000)
		if ev.GetType() == enet.EventNone {
			continue
		}

		switch ev.GetType() {
		case enet.EventConnect:
			playerID := uuid.Must(uuid.NewV4()).String()
			log.Println(fmt.Sprintf("new peer connected: %s playerID: %s", ev.GetPeer().GetAddress(), playerID))

			// send player id of client
			jsonObj := model.DataPacket{
				Type: utils.PacketPlayerID,
				Body: model.PlayerIDPacket{
					PlayerID: playerID,
				},
			}
			jsonStr, _ := json.Marshal(jsonObj)
			ev.GetPeer().SendString(string(jsonStr), ev.GetChannelID(), enet.PacketFlagReliable)
		case enet.EventDisconnect:
			log.Println("peer disconnected: ", ev.GetPeer().GetAddress())
		case enet.EventReceive:
			processClientAction(&g, ev)
		}
	}
}

func processClientAction(g *game.Game, ev enet.Event) {
	packet := ev.GetPacket()
	defer packet.Destroy()

	// test receive data
	log.Println("received from client:", string(packet.GetData()))
}
