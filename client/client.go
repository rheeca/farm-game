package main

import (
	"embed"
	"fmt"
	"guion-2d-project3/utils"
	"log"

	"github.com/codecat/go-enet"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var EmbeddedAssets embed.FS

func main() {
	enet.Initialize()
	host, err := enet.NewHost(nil, 1, 1, 0, 0)
	if err != nil {
		log.Fatal("failed to create host: ", err)
		return
	}

	peer, err := host.Connect(enet.NewAddress(utils.ServerAddress, utils.ServerPort), 1, 0)
	if err != nil {
		log.Fatal("failed to connect to server: ", err)
		return
	}

	gameObj := NewClientGame(EmbeddedAssets, peer, host)
	defer func() {
		host.Destroy()
		enet.Deinitialize()
	}()

	err = ebiten.RunGame(&gameObj)
	if err != nil {
		fmt.Println("failed to run game:", err)
	}
}
