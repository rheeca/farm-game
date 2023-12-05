package loader

import (
	"embed"
	"fmt"
	"guion-2d-project3/utils"
	"path"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type SoundCollection struct {
	BGMFirstTown   *audio.Player
	SFXChangeMap   *audio.Player
	SFXChicken     *audio.Player
	SFXChopTree    *audio.Player
	SFXCloseDoor   *audio.Player
	SFXCow         *audio.Player
	SFXCraft       *audio.Player
	SFXOpenDoor    *audio.Player
	SFXTillSoil    *audio.Player
	SFXWateringCan *audio.Player
}

func NewSoundCollection(EmbeddedAssets embed.FS) SoundCollection {
	audioContext := audio.NewContext(utils.SoundSampleRate)
	return SoundCollection{
		BGMFirstTown:   loadWavFromEmbedded(EmbeddedAssets, "first_town.wav", audioContext),
		SFXChangeMap:   loadWavFromEmbedded(EmbeddedAssets, "sfx_change_map.wav", audioContext),
		SFXChicken:     loadWavFromEmbedded(EmbeddedAssets, "sfx_chicken.wav", audioContext),
		SFXChopTree:    loadWavFromEmbedded(EmbeddedAssets, "sfx_chop_tree.wav", audioContext),
		SFXCloseDoor:   loadWavFromEmbedded(EmbeddedAssets, "sfx_close_door.wav", audioContext),
		SFXCow:         loadWavFromEmbedded(EmbeddedAssets, "sfx_cow.wav", audioContext),
		SFXCraft:       loadWavFromEmbedded(EmbeddedAssets, "sfx_craft.wav", audioContext),
		SFXOpenDoor:    loadWavFromEmbedded(EmbeddedAssets, "sfx_open_door.wav", audioContext),
		SFXTillSoil:    loadWavFromEmbedded(EmbeddedAssets, "sfx_till_soil.wav", audioContext),
		SFXWateringCan: loadWavFromEmbedded(EmbeddedAssets, "sfx_watering_can.wav", audioContext),
	}
}

func (s *SoundCollection) PlaySound(sound *audio.Player) {
	err := sound.Rewind()
	if err != nil {
		fmt.Println("failed to rewind sound")
	}
	sound.Play()
}

func loadWavFromEmbedded(EmbeddedAssets embed.FS, name string, context *audio.Context) (soundPlayer *audio.Player) {
	soundFile, err := EmbeddedAssets.Open(path.Join("client", "assets", "sounds", name))
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	sound, err := wav.DecodeWithoutResampling(soundFile)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	soundPlayer, err = context.NewPlayer(sound)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return soundPlayer
}
