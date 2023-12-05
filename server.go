package main

import (
	"embed"
	"fmt"
	"guion-2d-project3/entity/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

//go:embed assets/*
var EmbeddedAssets embed.FS

func main() {
	gameObj := game.NewGame(EmbeddedAssets)

	go func(player *audio.Player) {
		player.SetVolume(0.4)
		for {
			if !player.IsPlaying() {
				err := player.Rewind()
				if err != nil {
					fmt.Println("failed to rewind background music")
				}
				player.Play()
			}
		}
	}(gameObj.Sounds.BGMFirstTown)
	err := ebiten.RunGame(&gameObj)
	if err != nil {
		fmt.Println("failed to run game:", err)
	}
}
