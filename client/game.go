package main

import (
	"embed"
	"log"

	"github.com/codecat/go-enet"
	"github.com/hajimehoshi/ebiten/v2"
)

type ClientGame struct {
	peer enet.Peer
	host enet.Host
}

func NewClientGame(embeddedAssets embed.FS, peer enet.Peer, host enet.Host) ClientGame {
	return ClientGame{
		peer: peer,
		host: host,
	}
}

func (g *ClientGame) Update() error {
	// test send data
	g.peer.SendString("hello from client!", 0, enet.PacketFlagReliable)

	ev := g.host.Service(1000)
	switch ev.GetType() {
	case enet.EventConnect:
		log.Println("connected to the server!")

	case enet.EventDisconnect:
		log.Println("lost connection to the server!")

	case enet.EventReceive:
		packet := ev.GetPacket()

		// test receive data
		log.Println("received from server:", string(packet.GetData()))

		packet.Destroy()
	}
	return nil
}

func (g *ClientGame) Draw(screen *ebiten.Image) {
}

func (g *ClientGame) Layout(oWidth, oHeight int) (sWidth, sHeight int) {
	return oWidth, oHeight
}
